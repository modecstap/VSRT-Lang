package middleware

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestLoggerRestoresBodyAndRedactsSecrets(t *testing.T) {
	var logs bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&logs, nil))
	original := `{"email":"user@example.com","password":"secret","repeat_password":"secret2","token":"raw-token"}`
	handler := RequestLogger(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		if string(body) != original {
			t.Fatalf("unexpected body: %s", body)
		}
		w.WriteHeader(http.StatusCreated)
	}))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(original))
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, recorder.Code)
	}
	output := logs.String()
	if !strings.Contains(output, "method=POST") || !strings.Contains(output, "path=/register") || !strings.Contains(output, "status=201") {
		t.Fatalf("request metadata missing from log: %s", output)
	}
	redactedCount := strings.Count(output, "[REDACTED]")
	if redactedCount < 3 {
		t.Fatalf("expected at least 3 [REDACTED], got %d in %s", redactedCount, output)
	}
	for _, secret := range []string{"secret", "secret2", "raw-token"} {
		if strings.Contains(output, secret) {
			t.Fatalf("secret %q was not redacted: %s", secret, output)
		}
	}
}
