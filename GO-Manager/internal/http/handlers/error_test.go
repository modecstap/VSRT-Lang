package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		status    int
		code      string
		message   string
		details   []ErrorDetail
		wantCount int
	}{
		{
			name:      "writes details",
			status:    http.StatusBadRequest,
			code:      "invalid_request",
			message:   "invalid request",
			details:   []ErrorDetail{{Field: "email", Code: "required", Message: "email is required"}},
			wantCount: 1,
		},
		{
			name:      "omits empty details",
			status:    http.StatusUnauthorized,
			code:      "unauthorized",
			message:   "missing token",
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			recorder := httptest.NewRecorder()
			WriteError(recorder, tt.status, tt.code, tt.message, tt.details...)

			if recorder.Code != tt.status {
				t.Fatalf("status = %d, want %d", recorder.Code, tt.status)
			}
			if recorder.Header().Get("Content-Type") != "application/json" {
				t.Fatalf("Content-Type = %q", recorder.Header().Get("Content-Type"))
			}

			var response ErrorResponse
			if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
				t.Fatalf("decode: %v", err)
			}
			if response.Error.Code != tt.code {
				t.Fatalf("code = %q, want %q", response.Error.Code, tt.code)
			}
			if response.Error.Message != tt.message {
				t.Fatalf("message = %q, want %q", response.Error.Message, tt.message)
			}
			if len(response.Error.Details) != tt.wantCount {
				t.Fatalf("details = %d, want %d", len(response.Error.Details), tt.wantCount)
			}
		})
	}
}
