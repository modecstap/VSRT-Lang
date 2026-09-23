package user

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"VSRT-Lang/internal/http/handlers"
	domain "VSRT-Lang/internal/user"
)

type fakeUserService struct {
	user   *domain.User
	err    error
	gotID  domain.UserId
	called int
}

func (f *fakeUserService) SaveAvatar(domain.UserId, []byte) error {
	return errors.New("not implemented")
}

func (f *fakeUserService) Get(id domain.UserId) (*domain.User, error) {
	f.called++
	f.gotID = id
	return f.user, f.err
}

func TestHandler_GetMe(t *testing.T) {
	t.Parallel()
	png := []byte{0x89, 0x50, 0x4e, 0x47}
	dataURL := "data:image/png;base64," + base64.StdEncoding.EncodeToString(png)
	stored := &domain.User{Username: "stored-name", Email: "stored@example.com", Avatar: domain.Avatar{Bytes: png, MediaType: "image/jpeg"}}
	plain := &domain.User{Username: "stored-name", Email: "stored@example.com"}
	tests := []struct {
		name, path, code, msg string
		auth                  bool
		user                  *domain.User
		err                   error
		status                int
		avatar                *string
	}{
		{name: "png avatar", path: "/users/me", auth: true, user: stored, status: http.StatusOK, avatar: &dataURL},
		{name: "null avatar", path: "/users/me", auth: true, user: plain, status: http.StatusOK},
		{name: "no context", path: "/users/me", status: http.StatusUnauthorized, code: "unauthorized", msg: "missing user context"},
		{name: "missing user", path: "/users/me", auth: true, err: errors.New("user not found"), status: http.StatusNotFound, code: "user_not_found", msg: "user not found"},
		{name: "lookup failed", path: "/users/me", auth: true, err: errors.New("db down"), status: http.StatusInternalServerError, code: "user_lookup_failed", msg: "db down"},
		{name: "ignores query id", path: "/users/me?id=other-id", auth: true, user: stored, status: http.StatusOK, avatar: &dataURL},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := &fakeUserService{user: tt.user, err: tt.err}
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rw := httptest.NewRecorder()
			if tt.auth {
				rw = serveAuthed(t, NewHandler(svc, nil).GetMe, req, "user-1")
			} else {
				NewHandler(svc, nil).GetMe(rw, req)
			}
			if rw.Code != tt.status || tt.code == "" && rw.Header().Get("Content-Type") != "application/json" {
				t.Fatalf("status = %d, want %d; body=%s", rw.Code, tt.status, rw.Body.String())
			}
			if tt.auth != (svc.called == 1 && svc.gotID == "user-1") {
				t.Fatalf("Get called=%d id=%q", svc.called, svc.gotID)
			}
			if tt.code != "" {
				var got handlers.ErrorResponse
				if json.NewDecoder(rw.Body).Decode(&got) != nil || got.Error.Code != tt.code || got.Error.Message != tt.msg {
					t.Fatalf("error = %+v body=%s", got.Error, rw.Body.String())
				}
				return
			}
			var body map[string]json.RawMessage
			if json.Unmarshal(rw.Body.Bytes(), &body) != nil || len(body) != 3 || body["username"] == nil || body["email"] == nil || body["avatar"] == nil {
				t.Fatalf("body keys: %v", body)
			}
			var username, email string
			var avatar *string
			badUser := json.Unmarshal(body["username"], &username) != nil || json.Unmarshal(body["email"], &email) != nil || username != "stored-name" || email != "stored@example.com"
			badAvatar := json.Unmarshal(body["avatar"], &avatar) != nil || (tt.avatar == nil && string(body["avatar"]) != "null") || (tt.avatar != nil && (avatar == nil || *avatar != *tt.avatar))
			if badUser || badAvatar {
				t.Fatalf("body = %s", rw.Body.String())
			}
		})
	}
}
