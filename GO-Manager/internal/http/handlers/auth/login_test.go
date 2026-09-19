package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"VSRT-Lang/internal/auth"
)

func TestHandler_Login(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		body          string
		seedUser      bool
		wantStatus    int
		wantErrorCode string
	}{
		{
			name:       "returns tokens for valid credentials",
			body:       `{"login":"demo","password":"Password1"}`,
			seedUser:   true,
			wantStatus: http.StatusOK,
		},
		{
			name:          "rejects invalid json",
			body:          `{not-json`,
			wantStatus:    http.StatusBadRequest,
			wantErrorCode: "invalid_request",
		},
		{
			name:          "rejects unknown user",
			body:          `{"login":"missing","password":"Password1"}`,
			wantStatus:    http.StatusUnauthorized,
			wantErrorCode: "unauthorized",
		},
		{
			name:          "rejects wrong password",
			body:          `{"login":"demo","password":"WrongPass1"}`,
			seedUser:      true,
			wantStatus:    http.StatusUnauthorized,
			wantErrorCode: "unauthorized",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := newAuthHandler()
			if tt.seedUser {
				if _, err := h.service.Register("demo", "demo@example.com", "Password1"); err != nil {
					t.Fatalf("seed user: %v", err)
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rw := httptest.NewRecorder()
			h.Login(rw, req)

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

			var tokens auth.TokenPair
			if err := json.NewDecoder(rw.Body).Decode(&tokens); err != nil {
				t.Fatalf("decode tokens: %v", err)
			}
			if tokens.AccessToken == "" || tokens.RefreshToken == "" {
				t.Fatalf("expected tokens, got %+v", tokens)
			}
		})
	}
}
