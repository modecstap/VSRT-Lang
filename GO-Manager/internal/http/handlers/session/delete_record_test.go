package session

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	domain "VSRT-Lang/internal/session"
	"VSRT-Lang/internal/user"
)

func TestHandler_DeleteRecord(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		path          string
		auth          bool
		service       *fakeService
		wantStatus    int
		wantErrorCode string
	}{
		{
			name: "deletes owned record",
			path: "/sessions/7/records/hello",
			auth: true,
			service: &fakeService{
				deleteRecord: func(userID user.UserId, sessionID int64, phrase string) error {
					if userID != testUserID || sessionID != 7 || phrase != "hello" {
						return errors.New("unexpected delete record arguments")
					}
					return nil
				},
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "keeps decoded phrase",
			path: "/sessions/7/records/hello%20world",
			auth: true,
			service: &fakeService{
				deleteRecord: func(userID user.UserId, sessionID int64, phrase string) error {
					if userID != testUserID || sessionID != 7 || phrase != "hello world" {
						return errors.New("unexpected delete record arguments")
					}
					return nil
				},
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name:          "rejects missing phrase",
			path:          "/sessions/7/records",
			auth:          true,
			service:       &fakeService{},
			wantStatus:    http.StatusBadRequest,
			wantErrorCode: "invalid_request",
		},
		{
			name:          "rejects trailing slash",
			path:          "/sessions/7/records/",
			auth:          true,
			service:       &fakeService{},
			wantStatus:    http.StatusBadRequest,
			wantErrorCode: "invalid_request",
		},
		{
			name:          "rejects whitespace phrase",
			path:          "/sessions/7/records/%20",
			auth:          true,
			service:       &fakeService{},
			wantStatus:    http.StatusBadRequest,
			wantErrorCode: "invalid_request",
		},
		{
			name:          "rejects extra segment",
			path:          "/sessions/7/records/a/b",
			auth:          true,
			service:       &fakeService{},
			wantStatus:    http.StatusNotFound,
			wantErrorCode: "not_found",
		},
		{
			name:          "rejects invalid session id",
			path:          "/sessions/abc/records/hello",
			auth:          true,
			service:       &fakeService{},
			wantStatus:    http.StatusBadRequest,
			wantErrorCode: "invalid_request",
		},
		{
			name:          "rejects missing user context",
			path:          "/sessions/7/records/hello",
			auth:          false,
			service:       &fakeService{},
			wantStatus:    http.StatusUnauthorized,
			wantErrorCode: "unauthorized",
		},
		{
			name: "maps missing session",
			path: "/sessions/7/records/hello",
			auth: true,
			service: &fakeService{
				deleteRecord: func(user.UserId, int64, string) error {
					return domain.ErrSessionNotFound
				},
			},
			wantStatus:    http.StatusNotFound,
			wantErrorCode: "session_not_found",
		},
		{
			name: "maps unauthorized as not found",
			path: "/sessions/7/records/hello",
			auth: true,
			service: &fakeService{
				deleteRecord: func(user.UserId, int64, string) error {
					return domain.ErrUnauthorized
				},
			},
			wantStatus:    http.StatusNotFound,
			wantErrorCode: "session_not_found",
		},
		{
			name: "maps other service errors",
			path: "/sessions/7/records/hello",
			auth: true,
			service: &fakeService{
				deleteRecord: func(user.UserId, int64, string) error {
					return errors.New("db down")
				},
			},
			wantStatus:    http.StatusInternalServerError,
			wantErrorCode: "record_delete_failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := NewHandler(tt.service)
			req := httptest.NewRequest(http.MethodDelete, tt.path, nil)

			var rw *httptest.ResponseRecorder
			if tt.auth {
				rw = serveAuthed(t, testJWT(t), testUserID, h.DeleteRecord, req)
			} else {
				rw = httptest.NewRecorder()
				h.DeleteRecord(rw, req)
			}

			if rw.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body=%s", rw.Code, tt.wantStatus, rw.Body.String())
			}
			if tt.wantStatus == http.StatusNoContent && rw.Body.Len() != 0 {
				t.Fatalf("body = %q, want empty", rw.Body.String())
			}
			if tt.wantErrorCode != "" {
				got := decodeJSONError(t, rw)
				if got.Error.Code != tt.wantErrorCode {
					t.Fatalf("error code = %q, want %q", got.Error.Code, tt.wantErrorCode)
				}
			}
		})
	}
}
