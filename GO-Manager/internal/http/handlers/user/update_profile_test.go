package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"VSRT-Lang/internal/http/handlers"
	domain "VSRT-Lang/internal/user"
)

type fakeProfileService struct {
	id       domain.UserId
	username string
	email    string
	err      error
	called   int
}

func (f *fakeProfileService) UpdateProfile(id domain.UserId, username, email string) error {
	f.called++
	f.id = id
	f.username = username
	f.email = email
	return f.err
}

func (f *fakeProfileService) SaveAvatar(domain.UserId, []byte) error {
	return errors.New("not implemented")
}

func (f *fakeProfileService) Get(domain.UserId) (*domain.User, error) {
	return nil, errors.New("not implemented")
}

func TestHandler_UpdateProfile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name, body, code, msg, username, email string
		auth                                   bool
		err                                    error
		status                                 int
	}{
		{name: "username required", auth: true, body: `{"username":"","email":"a@b.c"}`, err: domain.ErrUsernameRequired, status: http.StatusBadRequest, code: "username_required", msg: "username required", username: "", email: "a@b.c"},
		{name: "email required", auth: true, body: `{"username":"ada","email":""}`, err: domain.ErrEmailRequired, status: http.StatusBadRequest, code: "email_required", msg: "email required", username: "ada"},
		{name: "invalid email", auth: true, body: `{"username":"ada","email":"userexample.com"}`, err: domain.ErrInvalidEmail, status: http.StatusBadRequest, code: "invalid_email", msg: "invalid email", username: "ada", email: "userexample.com"},
		{name: "username taken", auth: true, body: `{"username":"ada","email":"ada@example.com"}`, err: domain.ErrUsernameTaken, status: http.StatusConflict, code: "username_taken", msg: "username taken", username: "ada", email: "ada@example.com"},
		{name: "email taken", auth: true, body: `{"username":"ada","email":"ada@example.com"}`, err: domain.ErrEmailTaken, status: http.StatusConflict, code: "email_taken", msg: "email taken", username: "ada", email: "ada@example.com"},
		{name: "missing user", auth: true, body: `{"username":"ada","email":"ada@example.com"}`, err: domain.ErrNotFound, status: http.StatusNotFound, code: "user_not_found", msg: "user not found", username: "ada", email: "ada@example.com"},
		{name: "save failed", auth: true, body: `{"username":"ada","email":"ada@example.com"}`, err: errors.New("db down"), status: http.StatusInternalServerError, code: "profile_save_failed", msg: "db down", username: "ada", email: "ada@example.com"},
		{name: "success", auth: true, body: `{"id":"other","username":"ada","email":"ada@example.com"}`, status: http.StatusNoContent, username: "ada", email: "ada@example.com"},
		{name: "bad json", auth: true, body: `{`, status: http.StatusBadRequest, code: "invalid_request", msg: "invalid request"},
		{name: "no context", body: `{"username":"ada","email":"ada@example.com"}`, status: http.StatusUnauthorized, code: "unauthorized", msg: "missing user context"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := &fakeProfileService{err: tt.err}
			req := httptest.NewRequest(http.MethodPost, "/users/me", strings.NewReader(tt.body))
			rw := httptest.NewRecorder()
			if tt.auth {
				rw = serveAuthed(t, NewHandler(svc, nil).UpdateProfile, req, "user-1")
			} else {
				NewHandler(svc, nil).UpdateProfile(rw, req)
			}
			if rw.Code != tt.status {
				t.Fatalf("status = %d, want %d; body=%s", rw.Code, tt.status, rw.Body.String())
			}
			wantCall := tt.code != "invalid_request" && tt.code != "unauthorized"
			if wantCall != (svc.called == 1 && svc.id == "user-1" && svc.username == tt.username && svc.email == tt.email) {
				t.Fatalf("UpdateProfile called=%d id=%q username=%q email=%q", svc.called, svc.id, svc.username, svc.email)
			}
			if tt.status == http.StatusNoContent {
				if rw.Body.Len() != 0 {
					t.Fatalf("body = %s, want empty", rw.Body.String())
				}
				return
			}
			var got handlers.ErrorResponse
			if json.NewDecoder(rw.Body).Decode(&got) != nil || got.Error.Code != tt.code || got.Error.Message != tt.msg {
				t.Fatalf("error = %+v body=%s", got.Error, rw.Body.String())
			}
		})
	}
}
