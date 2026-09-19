package session

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	domain "VSRT-Lang/internal/session"
	"VSRT-Lang/internal/user"
)

func TestHandler_DeleteSession(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		path          string
		auth          bool
		service       *fakeService
		wantStatus    int
		wantErrorCode string
		wantPlainBody string
	}{
		{
			name: "deletes owned session",
			path: "/sessions/7",
			auth: true,
			service: &fakeService{
				deleteSession: func(userID user.UserId, sessionID int64) error {
					if userID != testUserID || sessionID != 7 {
						return errors.New("unexpected delete arguments")
					}
					return nil
				},
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name:          "rejects unknown route",
			path:          "/other/7",
			auth:          true,
			service:       &fakeService{},
			wantStatus:    http.StatusNotFound,
			wantPlainBody: "not found",
		},
		{
			name:          "rejects invalid session id",
			path:          "/sessions/abc",
			auth:          true,
			service:       &fakeService{},
			wantStatus:    http.StatusBadRequest,
			wantPlainBody: "invalid session ID",
		},
		{
			name:          "rejects missing user context",
			path:          "/sessions/7",
			auth:          false,
			service:       &fakeService{},
			wantStatus:    http.StatusUnauthorized,
			wantPlainBody: "unauthorized",
		},
		{
			name: "maps unauthorized as not found",
			path: "/sessions/7",
			auth: true,
			service: &fakeService{
				deleteSession: func(user.UserId, int64) error {
					return domain.ErrUnauthorized
				},
			},
			wantStatus:    http.StatusNotFound,
			wantErrorCode: "session_not_found",
		},
		{
			name: "maps other service errors",
			path: "/sessions/7",
			auth: true,
			service: &fakeService{
				deleteSession: func(user.UserId, int64) error {
					return errors.New("db down")
				},
			},
			wantStatus:    http.StatusInternalServerError,
			wantErrorCode: "session_delete_failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := NewHandler(tt.service)
			req := httptest.NewRequest(http.MethodDelete, tt.path, nil)

			var rw *httptest.ResponseRecorder
			if tt.auth {
				rw = serveAuthed(t, testJWT(t), testUserID, h.DeleteSession, req)
			} else {
				rw = httptest.NewRecorder()
				h.DeleteSession(rw, req)
			}

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

			if tt.wantPlainBody != "" && !strings.Contains(rw.Body.String(), tt.wantPlainBody) {
				t.Fatalf("body = %q, want substring %q", rw.Body.String(), tt.wantPlainBody)
			}
		})
	}
}
