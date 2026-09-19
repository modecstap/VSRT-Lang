package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"VSRT-Lang/internal/auth"
	"VSRT-Lang/internal/http/handlers"
	"VSRT-Lang/internal/http/middleware"
	"VSRT-Lang/internal/session"
	domain "VSRT-Lang/internal/user"
)

type fakeUserRepo struct{}

func (fakeUserRepo) Create(*domain.User) error                   { return nil }
func (fakeUserRepo) FindByEmail(string) (*domain.User, error)    { return nil, nil }
func (fakeUserRepo) FindByUsername(string) (*domain.User, error) { return nil, nil }
func (fakeUserRepo) FindByID(string) (*domain.User, error)       { return nil, nil }

type fakeSessionRepo struct {
	sessions []session.Session
	err      error
	lastUser domain.UserId
}

func (r *fakeSessionRepo) FindByUser(userID domain.UserId) ([]session.Session, error) {
	r.lastUser = userID
	if r.err != nil {
		return nil, r.err
	}
	return r.sessions, nil
}

func serveAuthed(t *testing.T, handler http.HandlerFunc, req *http.Request, userID string) *httptest.ResponseRecorder {
	t.Helper()
	jwtSvc := auth.NewJWTService("test-secret")
	token, err := jwtSvc.GenerateAccessToken(userID, "demo", "demo@example.com")
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	rw := httptest.NewRecorder()
	middleware.Auth(jwtSvc)(handler).ServeHTTP(rw, req)
	return rw
}

func TestHandler_GetUserSessions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		path          string
		auth          bool
		repo          *fakeSessionRepo
		wantStatus    int
		wantErrorCode string
		wantCount     int
	}{
		{
			name: "returns sessions for authenticated user",
			path: "/users/sessions",
			auth: true,
			repo: &fakeSessionRepo{
				sessions: []session.Session{
					{ID: 1, User: "1", Name: "demo"},
				},
			},
			wantStatus: http.StatusOK,
			wantCount:  1,
		},
		{
			name:          "rejects unknown route",
			path:          "/users/other",
			auth:          true,
			repo:          &fakeSessionRepo{},
			wantStatus:    http.StatusNotFound,
			wantErrorCode: "not_found",
		},
		{
			name:          "rejects missing user context",
			path:          "/users/sessions",
			auth:          false,
			repo:          &fakeSessionRepo{},
			wantStatus:    http.StatusInternalServerError,
			wantErrorCode: "invalid_token",
		},
		{
			name:          "maps repository errors",
			path:          "/users/sessions",
			auth:          true,
			repo:          &fakeSessionRepo{err: errors.New("not found")},
			wantStatus:    http.StatusNotFound,
			wantErrorCode: "sessions_not_found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := NewHandler(fakeUserRepo{}, tt.repo)
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)

			var rw *httptest.ResponseRecorder
			if tt.auth {
				rw = serveAuthed(t, h.GetUserSessions, req, "1")
			} else {
				rw = httptest.NewRecorder()
				h.GetUserSessions(rw, req)
			}

			if rw.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body=%s", rw.Code, tt.wantStatus, rw.Body.String())
			}

			if tt.wantErrorCode != "" {
				var got handlers.ErrorResponse
				if err := json.NewDecoder(rw.Body).Decode(&got); err != nil {
					t.Fatalf("decode error: %v", err)
				}
				if got.Error.Code != tt.wantErrorCode {
					t.Fatalf("error code = %q, want %q", got.Error.Code, tt.wantErrorCode)
				}
				return
			}

			if tt.auth && tt.repo.lastUser != "1" {
				t.Fatalf("FindByUser got %q, want 1", tt.repo.lastUser)
			}

			var got []session.Session
			if err := json.NewDecoder(rw.Body).Decode(&got); err != nil {
				t.Fatalf("decode sessions: %v", err)
			}
			if len(got) != tt.wantCount {
				t.Fatalf("sessions = %d, want %d", len(got), tt.wantCount)
			}
		})
	}
}
