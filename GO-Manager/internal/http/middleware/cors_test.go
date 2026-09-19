package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORS(t *testing.T) {
	t.Parallel()

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	tests := []struct {
		name            string
		cfg             CORSConfig
		method          string
		origin          string
		wantStatus      int
		wantAllowOrigin string
		wantCreds       string
		wantNext        bool
	}{
		{
			name: "allows wildcard origin",
			cfg: CORSConfig{
				AllowedOrigins: []string{"*"},
				AllowedMethods: []string{http.MethodGet, http.MethodPost},
				AllowedHeaders: []string{"Content-Type", "Authorization"},
			},
			method:          http.MethodGet,
			origin:          "https://app.example.com",
			wantStatus:      http.StatusOK,
			wantAllowOrigin: "*",
			wantNext:        true,
		},
		{
			name: "echoes matching origin",
			cfg: CORSConfig{
				AllowedOrigins:   []string{"https://app.example.com"},
				AllowedMethods:   []string{http.MethodGet},
				AllowCredentials: true,
			},
			method:          http.MethodGet,
			origin:          "https://app.example.com",
			wantStatus:      http.StatusOK,
			wantAllowOrigin: "https://app.example.com",
			wantCreds:       "true",
			wantNext:        true,
		},
		{
			name: "ignores unknown origin",
			cfg: CORSConfig{
				AllowedOrigins: []string{"https://app.example.com"},
			},
			method:     http.MethodGet,
			origin:     "https://evil.example.com",
			wantStatus: http.StatusOK,
			wantNext:   true,
		},
		{
			name: "handles preflight without calling next",
			cfg: CORSConfig{
				AllowedOrigins: []string{"*"},
				AllowedMethods: []string{http.MethodGet, http.MethodOptions},
			},
			method:          http.MethodOptions,
			origin:          "https://app.example.com",
			wantStatus:      http.StatusNoContent,
			wantAllowOrigin: "*",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			called := false
			handler := CORS(tt.cfg)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
				next.ServeHTTP(w, r)
			}))

			req := httptest.NewRequest(tt.method, "/login", nil)
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}
			rw := httptest.NewRecorder()
			handler.ServeHTTP(rw, req)

			if rw.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rw.Code, tt.wantStatus)
			}
			if got := rw.Header().Get("Access-Control-Allow-Origin"); got != tt.wantAllowOrigin {
				t.Fatalf("Allow-Origin = %q, want %q", got, tt.wantAllowOrigin)
			}
			if got := rw.Header().Get("Access-Control-Allow-Credentials"); got != tt.wantCreds {
				t.Fatalf("Allow-Credentials = %q, want %q", got, tt.wantCreds)
			}
			if called != tt.wantNext {
				t.Fatalf("next called = %v, want %v", called, tt.wantNext)
			}
		})
	}
}
