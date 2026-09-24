package passwordreset

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"VSRT-Lang/internal/database/memory"
	handlers "VSRT-Lang/internal/http/handlers"
	"VSRT-Lang/internal/passwordreset"
	"VSRT-Lang/internal/user"
)

type stubMailer struct {
	err   error
	calls int
}

func (m *stubMailer) Send(to, link string) error {
	m.calls++
	return m.err
}

func setupHandler(t *testing.T, mail *stubMailer) (*Handler, *memory.UserRepository) {
	t.Helper()
	users := memory.NewUserRepository()
	links := memory.NewPasswordResetRepository()
	refresh := memory.NewRefreshTokenRepository()
	tx := memory.NewPasswordResetTransactor(users, links, refresh)
	clk := time.Date(2026, 9, 24, 12, 0, 0, 0, time.UTC)
	svc := passwordreset.NewService(users, links, mail, tx, "http://localhost", func() time.Time { return clk })
	return NewHandler(svc), users
}

func doForgot(t *testing.T, h *Handler, body any) *httptest.ResponseRecorder {
	t.Helper()
	raw, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/forgot-password", bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")
	rw := httptest.NewRecorder()
	h.ForgotPassword(rw, req)
	return rw
}

func doReset(t *testing.T, h *Handler, body any) *httptest.ResponseRecorder {
	t.Helper()
	raw, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/reset-password", bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")
	rw := httptest.NewRecorder()
	h.ResetPassword(rw, req)
	return rw
}

func TestForgotPassword(t *testing.T) {
	mail := &stubMailer{}
	h, users := setupHandler(t, mail)
	hashed, _ := user.HashPassword("Password1")
	_ = users.Create(user.NewUser("demo", "demo@example.com", hashed))

	cases := []struct {
		name   string
		body   any
		status int
		code   string
		msg    string
		calls  int
	}{
		{name: "invalid email", body: map[string]string{"email": "bad"}, status: 400, code: "invalid_email", msg: "invalid email", calls: 0},
		{name: "unknown email", body: map[string]string{"email": "no@example.com"}, status: 204, calls: 0},
		{name: "success", body: map[string]string{"email": "demo@example.com"}, status: 204, calls: 1},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mail.calls = 0
			rw := doForgot(t, h, c.body)
			if rw.Code != c.status {
				t.Fatalf("status = %d, want %d; body=%s", rw.Code, c.status, rw.Body.String())
			}
			if c.status == 204 {
				if mail.calls != c.calls {
					t.Fatalf("mail calls = %d, want %d", mail.calls, c.calls)
				}
				return
			}
			var resp handlers.ErrorResponse
			if err := json.NewDecoder(rw.Body).Decode(&resp); err != nil {
				t.Fatalf("decode: %v", err)
			}
			if resp.Error.Code != c.code || resp.Error.Message != c.msg {
				t.Fatalf("error = %+v", resp.Error)
			}
			if mail.calls != c.calls {
				t.Fatalf("mail calls = %d, want %d", mail.calls, c.calls)
			}
		})
	}

	t.Run("rate limited", func(t *testing.T) {
		mail.calls = 0
		_ = doForgot(t, h, map[string]string{"email": "demo@example.com"})
		rw := doForgot(t, h, map[string]string{"email": "demo@example.com"})
		if rw.Code != 400 {
			t.Fatalf("status = %d", rw.Code)
		}
		var resp handlers.ErrorResponse
		_ = json.NewDecoder(rw.Body).Decode(&resp)
		if resp.Error.Code != "rate_limited" || resp.Error.Message != passwordreset.ErrRateLimited.Error() {
			t.Fatalf("error = %+v", resp.Error)
		}
	})

	t.Run("send failed", func(t *testing.T) {
		users2 := memory.NewUserRepository()
		links := memory.NewPasswordResetRepository()
		refresh := memory.NewRefreshTokenRepository()
		tx := memory.NewPasswordResetTransactor(users2, links, refresh)
		failMail := &stubMailer{err: errors.New("down")}
		clk := time.Date(2026, 9, 24, 12, 0, 0, 0, time.UTC)
		svc := passwordreset.NewService(users2, links, failMail, tx, "http://localhost", func() time.Time { return clk })
		h2 := NewHandler(svc)
		hashed, _ := user.HashPassword("Password1")
		_ = users2.Create(user.NewUser("demo", "demo@example.com", hashed))

		rw := doForgot(t, h2, map[string]string{"email": "demo@example.com"})
		if rw.Code != 400 {
			t.Fatalf("status = %d", rw.Code)
		}
		var resp handlers.ErrorResponse
		_ = json.NewDecoder(rw.Body).Decode(&resp)
		if resp.Error.Code != "email_failed" || resp.Error.Message != passwordreset.ErrSendFailed.Error() {
			t.Fatalf("error = %+v", resp.Error)
		}
	})
}

func TestResetPassword(t *testing.T) {
	mail := &stubMailer{}
	h, users := setupHandler(t, mail)
	hashed, _ := user.HashPassword("Password1")
	_ = users.Create(user.NewUser("demo", "demo@example.com", hashed))
	_ = doForgot(t, h, map[string]string{"email": "demo@example.com"})

	cases := []struct {
		name   string
		body   any
		status int
		code   string
		msg    string
	}{
		{name: "empty token", body: map[string]string{"token": "", "password": "NewPass99"}, status: 400, code: "invalid_link", msg: passwordreset.ErrInvalidLink.Error()},
		{name: "weak password", body: map[string]string{"token": "anything", "password": "weak"}, status: 400, code: "weak_password", msg: user.ErrWeakPassword.Error()},
		{name: "dead link", body: map[string]string{"token": "dead", "password": "NewPass99"}, status: 400, code: "invalid_link", msg: passwordreset.ErrInvalidLink.Error()},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			rw := doReset(t, h, c.body)
			if rw.Code != c.status {
				t.Fatalf("status = %d, want %d; body=%s", rw.Code, c.status, rw.Body.String())
			}
			var resp handlers.ErrorResponse
			if err := json.NewDecoder(rw.Body).Decode(&resp); err != nil {
				t.Fatalf("decode: %v", err)
			}
			if resp.Error.Code != c.code || resp.Error.Message != c.msg {
				t.Fatalf("error = %+v", resp.Error)
			}
		})
	}
}
