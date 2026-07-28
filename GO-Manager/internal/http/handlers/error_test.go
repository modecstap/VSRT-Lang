package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteError(t *testing.T) {
	recorder := httptest.NewRecorder()

	WriteError(
		recorder,
		http.StatusBadRequest,
		"invalid_request",
		"invalid request",
		ErrorDetail{Field: "email", Code: "required", Message: "email is required"},
	)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}

	if recorder.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("expected application/json content type, got %q", recorder.Header().Get("Content-Type"))
	}

	var response ErrorResponse
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	if response.Error.Code != "invalid_request" {
		t.Fatalf("expected error code %q, got %q", "invalid_request", response.Error.Code)
	}

	if response.Error.Message != "invalid request" {
		t.Fatalf("expected error message %q, got %q", "invalid request", response.Error.Message)
	}

	if len(response.Error.Details) != 1 {
		t.Fatalf("expected 1 error detail, got %d", len(response.Error.Details))
	}
}
