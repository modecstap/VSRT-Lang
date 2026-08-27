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
	handler := RequestLogger(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		if string(body) != `{"email":"user@example.com","password":"secret"}` {
			t.Fatalf("unexpected body: %s", body)
		}
		w.WriteHeader(http.StatusCreated)
	}))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{"email":"user@example.com","password":"secret"}`))
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, recorder.Code)
	}
	output := logs.String()
	if !strings.Contains(output, "method=POST") || !strings.Contains(output, "path=/register") || !strings.Contains(output, "status=201") {
		t.Fatalf("request metadata missing from log: %s", output)
	}
	if !strings.Contains(output, "password:[REDACTED]") || strings.Contains(output, "password:secret") {
		t.Fatalf("secret was not redacted: %s", output)
	}
}
