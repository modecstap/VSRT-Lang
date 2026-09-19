package session

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	domain "VSRT-Lang/internal/session"
)

func TestHandler_SaveRecord(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		path          string
		body          string
		auth          bool
		service       *fakeService
		wantStatus    int
		wantErrorCode string
	}{
		{
			name: "saves record into session",
			path: "/sessions/3/records",
			body: `{"phrase":"hello","context":"world"}`,
			auth: true,
			service: &fakeService{
				addRecord: func(cmd domain.AddRecordCommand) (domain.Record, error) {
					if cmd.UserId != testUserID || cmd.SessionId != 3 || cmd.Phrase != "hello" || cmd.Context != "world" {
						return domain.Record{}, errors.New("unexpected add record arguments")
					}
					return domain.Record{Phrase: cmd.Phrase}, nil
				},
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:          "rejects unknown route",
			path:          "/sessions/3",
			body:          `{"phrase":"hello","context":"world"}`,
			auth:          true,
			service:       &fakeService{},
			wantStatus:    http.StatusNotFound,
			wantErrorCode: "not_found",
		},
		{
			name:          "rejects invalid session id",
			path:          "/sessions/abc/records",
			body:          `{"phrase":"hello","context":"world"}`,
			auth:          true,
			service:       &fakeService{},
			wantStatus:    http.StatusBadRequest,
			wantErrorCode: "invalid_request",
		},
		{
			name:          "rejects invalid json",
			path:          "/sessions/3/records",
			body:          `{not-json`,
			auth:          true,
			service:       &fakeService{},
			wantStatus:    http.StatusBadRequest,
			wantErrorCode: "invalid_request",
		},
		{
			name:          "rejects missing user context",
			path:          "/sessions/3/records",
			body:          `{"phrase":"hello","context":"world"}`,
			auth:          false,
			service:       &fakeService{},
			wantStatus:    http.StatusUnauthorized,
			wantErrorCode: "unauthorized",
		},
		{
			name: "maps missing session",
			path: "/sessions/3/records",
			body: `{"phrase":"hello","context":"world"}`,
			auth: true,
			service: &fakeService{
				addRecord: func(domain.AddRecordCommand) (domain.Record, error) {
					return domain.Record{}, domain.ErrSessionNotFound
				},
			},
			wantStatus:    http.StatusNotFound,
			wantErrorCode: "session_not_found",
		},
		{
			name: "maps other service errors",
			path: "/sessions/3/records",
			body: `{"phrase":"hello","context":"world"}`,
			auth: true,
			service: &fakeService{
				addRecord: func(domain.AddRecordCommand) (domain.Record, error) {
					return domain.Record{}, errors.New("translator down")
				},
			},
			wantStatus:    http.StatusInternalServerError,
			wantErrorCode: "record_save_failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := NewHandler(tt.service)
			req := httptest.NewRequest(http.MethodPost, tt.path, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			var rw *httptest.ResponseRecorder
			if tt.auth {
				rw = serveAuthed(t, testJWT(t), testUserID, h.SaveRecord, req)
			} else {
				rw = httptest.NewRecorder()
				h.SaveRecord(rw, req)
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

			var got map[string]any
			if err := json.NewDecoder(rw.Body).Decode(&got); err != nil {
				t.Fatalf("decode response: %v", err)
			}
			if got["message"] != "record saved" {
				t.Fatalf("response = %#v", got)
			}
		})
	}
}
