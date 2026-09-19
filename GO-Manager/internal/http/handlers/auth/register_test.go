package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_Register(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		body          string
		setup         func(*testing.T, *Handler)
		wantStatus    int
		wantErrorCode string
		wantEmail     string
	}{
		{
			name:       "creates user",
			body:       `{"username":"demo","email":"demo@example.com","password":"Password1"}`,
			wantStatus: http.StatusOK,
			wantEmail:  "demo@example.com",
		},
		{
			name:          "rejects invalid json",
			body:          `{not-json`,
			wantStatus:    http.StatusBadRequest,
			wantErrorCode: "invalid_request",
		},
		{
			name:          "rejects invalid email",
			body:          `{"username":"demo","email":"not-an-email","password":"Password1"}`,
			wantStatus:    http.StatusBadRequest,
			wantErrorCode: "registration_failed",
		},
		{
			name:          "rejects weak password",
			body:          `{"username":"demo","email":"demo@example.com","password":"short"}`,
			wantStatus:    http.StatusBadRequest,
			wantErrorCode: "registration_failed",
		},
		{
			name: "rejects duplicate email",
			body: `{"username":"demo2","email":"demo@example.com","password":"Password1"}`,
			setup: func(t *testing.T, h *Handler) {
				t.Helper()
				if _, err := h.service.Register("demo", "demo@example.com", "Password1"); err != nil {
					t.Fatalf("seed user: %v", err)
				}
			},
			wantStatus:    http.StatusBadRequest,
			wantErrorCode: "registration_failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := newAuthHandler()
			if tt.setup != nil {
				tt.setup(t, h)
			}

			req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rw := httptest.NewRecorder()
			h.Register(rw, req)

			if rw.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body=%s", rw.Code, tt.wantStatus, rw.Body.String())
			}

			if tt.wantErrorCode != "" {
				got := decodeJSONError(t, rw)
				if got.Error.Code != tt.wantErrorCode {
					t.Fatalf("error code = %q, want %q", got.Error.Code, tt.wantErrorCode)
				}
				return
			}

			var got map[string]any
			if err := json.NewDecoder(rw.Body).Decode(&got); err != nil {
				t.Fatalf("decode response: %v", err)
			}
			if got["email"] != tt.wantEmail {
				t.Fatalf("email = %#v, want %q", got["email"], tt.wantEmail)
			}
			if got["id"] == "" || got["id"] == nil {
				t.Fatal("expected user id in response")
			}
		})
	}
}

func TestHandler_RegisterUsesEmailAsUsernameWhenMissing(t *testing.T) {
	t.Parallel()

	h := newAuthHandler()
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{"email":"solo@example.com","password":"Password1"}`))
	rw := httptest.NewRecorder()
	h.Register(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rw.Code, rw.Body.String())
	}

	login := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"login":"solo@example.com","password":"Password1"}`))
	loginRW := httptest.NewRecorder()
	h.Login(loginRW, login)

	if loginRW.Code != http.StatusOK {
		t.Fatalf("login status = %d, want 200; body=%s", loginRW.Code, loginRW.Body.String())
	}
}
