package session

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"VSRT-Lang/internal/user"
)

func TestHandler_CreateSession(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		body          string
		auth          bool
		service       *fakeService
		wantStatus    int
		wantErrorCode string
		wantID        int64
		wantName      string
	}{
		{
			name: "creates session for authenticated user",
			body: `{"name":"demo"}`,
			auth: true,
			service: &fakeService{
				newSession: func(userID user.UserId, name string) (int64, error) {
					if userID != testUserID || name != "demo" {
						return 0, errors.New("unexpected session arguments")
					}
					return 42, nil
				},
			},
			wantStatus: http.StatusCreated,
			wantID:     42,
			wantName:   "demo",
		},
		{
			name:          "rejects invalid json",
			body:          `{not-json`,
			auth:          true,
			service:       &fakeService{},
			wantStatus:    http.StatusBadRequest,
			wantErrorCode: "invalid_request",
		},
		{
			name:          "rejects missing user context",
			body:          `{"name":"demo"}`,
			auth:          false,
			service:       &fakeService{},
			wantStatus:    http.StatusUnauthorized,
			wantErrorCode: "unauthorized",
		},
		{
			name: "maps service failure",
			body: `{"name":"demo"}`,
			auth: true,
			service: &fakeService{
				newSession: func(user.UserId, string) (int64, error) {
					return 0, errors.New("db down")
				},
			},
			wantStatus:    http.StatusInternalServerError,
			wantErrorCode: "session_save_failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := NewHandler(tt.service)
			req := httptest.NewRequest(http.MethodPost, "/sessions", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			var rw *httptest.ResponseRecorder
			if tt.auth {
				rw = serveAuthed(t, testJWT(t), testUserID, h.CreateSession, req)
			} else {
				rw = httptest.NewRecorder()
				h.CreateSession(rw, req)
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

			if rw.Header().Get("Content-Type") != "application/json" {
				t.Fatalf("Content-Type = %q, want application/json", rw.Header().Get("Content-Type"))
			}

			var got sessionResponse
			if err := json.NewDecoder(rw.Body).Decode(&got); err != nil {
				t.Fatalf("decode response: %v", err)
			}
			if got.ID != tt.wantID || got.Name != tt.wantName {
				t.Fatalf("response = %+v, want id=%d name=%q", got, tt.wantID, tt.wantName)
			}
		})
	}
}
