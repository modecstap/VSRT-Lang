package session

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	domain "VSRT-Lang/internal/session"
	"VSRT-Lang/internal/user"
)

func TestHandler_GetRecords(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		path            string
		auth            bool
		service         *fakeService
		wantStatus      int
		wantErrorCode   string
		wantRecordCount int
	}{
		{
			name: "returns session records",
			path: "/sessions/9/records",
			auth: true,
			service: &fakeService{
				getSession: func(userID user.UserId, sessionID int64) (domain.Session, error) {
					if userID != testUserID || sessionID != 9 {
						return domain.Session{}, errors.New("unexpected get session arguments")
					}
					return domain.Session{
						ID:   9,
						User: testUserID,
						Records: map[string]*domain.Record{
							"hello": {Phrase: "hello", Count: 1},
						},
					}, nil
				},
			},
			wantStatus:      http.StatusOK,
			wantRecordCount: 1,
		},
		{
			name:          "rejects unknown route",
			path:          "/sessions/9",
			auth:          true,
			service:       &fakeService{},
			wantStatus:    http.StatusNotFound,
			wantErrorCode: "not_found",
		},
		{
			name:          "rejects invalid session id",
			path:          "/sessions/abc/records",
			auth:          true,
			service:       &fakeService{},
			wantStatus:    http.StatusBadRequest,
			wantErrorCode: "invalid_request",
		},
		{
			name:          "rejects missing user context",
			path:          "/sessions/9/records",
			auth:          false,
			service:       &fakeService{},
			wantStatus:    http.StatusUnauthorized,
			wantErrorCode: "unauthorized",
		},
		{
			name: "maps missing session",
			path: "/sessions/9/records",
			auth: true,
			service: &fakeService{
				getSession: func(user.UserId, int64) (domain.Session, error) {
					return domain.Session{}, errors.New("session not found")
				},
			},
			wantStatus:    http.StatusNotFound,
			wantErrorCode: "session_not_found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := NewHandler(tt.service)
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)

			var rw *httptest.ResponseRecorder
			if tt.auth {
				rw = serveAuthed(t, testJWT(t), testUserID, h.GetRecords, req)
			} else {
				rw = httptest.NewRecorder()
				h.GetRecords(rw, req)
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

			var got recordsResponse
			if err := json.NewDecoder(rw.Body).Decode(&got); err != nil {
				t.Fatalf("decode response: %v", err)
			}
			if len(got.Records) != tt.wantRecordCount {
				t.Fatalf("records = %d, want %d", len(got.Records), tt.wantRecordCount)
			}
		})
	}
}
