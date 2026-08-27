package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *loggingResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *loggingResponseWriter) Write(body []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	return w.ResponseWriter.Write(body)
}

func RequestLogger(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			started := time.Now()
			body := readRequestBody(r)
			response := &loggingResponseWriter{ResponseWriter: w}

			next.ServeHTTP(response, r)

			attrs := []any{
				"method", r.Method,
				"path", r.URL.Path,
				"status", response.status,
				"duration_ms", time.Since(started).Milliseconds(),
			}
			if len(body) > 0 {
				attrs = append(attrs, "request_body", body)
			}
			logger.Info("http request", attrs...)
		})
	}
}

func readRequestBody(r *http.Request) map[string]any {
	if r.Body == nil || r.ContentLength == 0 {
		return nil
	}

	data, err := io.ReadAll(r.Body)
	r.Body.Close()
	r.Body = io.NopCloser(bytes.NewReader(data))
	if err != nil {
		return nil
	}

	var payload map[string]any
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil
	}
	for _, key := range []string{"password", "token", "access_token", "refresh_token", "authorization"} {
		if _, ok := payload[key]; ok {
			payload[key] = "[REDACTED]"
		}
	}

	return payload
}
