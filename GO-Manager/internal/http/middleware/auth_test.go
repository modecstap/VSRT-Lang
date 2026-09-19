package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"VSRT-Lang/internal/auth"
)

func TestAuth(t *testing.T) {
	t.Parallel()

	jwtSvc := auth.NewJWTService("test-secret")
	validToken, err := jwtSvc.GenerateAccessToken("user-1", "demo", "demo@example.com")
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	tests := []struct {
		name       string
		header     string
		wantStatus int
		wantUser   string
		wantEmail  string
	}{
		{
			name:       "allows valid bearer token",
			header:     "Bearer " + validToken,
			wantStatus: http.StatusOK,
			wantUser:   "user-1",
			wantEmail:  "demo@example.com",
		},
		{
			name:       "rejects missing header",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "rejects non-bearer scheme",
			header:     "Token " + validToken,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "rejects malformed header",
			header:     "Bearer",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "rejects invalid token",
			header:     "Bearer not-a-token",
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				userID, email, ok := UserFromContext(r.Context())
				if !ok {
					t.Fatal("expected user context")
				}
				if string(userID) != tt.wantUser {
					t.Fatalf("userID = %q, want %q", userID, tt.wantUser)
				}
				if email != tt.wantEmail {
					t.Fatalf("email = %q, want %q", email, tt.wantEmail)
				}
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/sessions", nil)
			if tt.header != "" {
				req.Header.Set("Authorization", tt.header)
			}
			rw := httptest.NewRecorder()
			Auth(jwtSvc)(next).ServeHTTP(rw, req)

			if rw.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body=%s", rw.Code, tt.wantStatus, rw.Body.String())
			}
		})
	}
}

func TestUserFromContext(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	userID, email, ok := UserFromContext(req.Context())
	if ok || userID != "" || email != "" {
		t.Fatalf("expected empty context, got %q %q %v", userID, email, ok)
	}
}

func TestAuthRejectsTokenFromDifferentSecret(t *testing.T) {
	t.Parallel()

	other := auth.NewJWTService("other-secret")
	token, err := other.GenerateAccessToken("user-1", "demo", "demo@example.com")
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rw := httptest.NewRecorder()
	Auth(auth.NewJWTService("test-secret"))(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("handler should not run")
	})).ServeHTTP(rw, req)

	if rw.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", rw.Code)
	}
	if !strings.Contains(rw.Body.String(), "invalid token") {
		t.Fatalf("body = %q", rw.Body.String())
	}
}
