package auth

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"VSRT-Lang/internal/auth"
	"VSRT-Lang/internal/database/memory"
	"VSRT-Lang/internal/http/handlers"
)

func newAuthHandler() *Handler {
	return NewAuth(auth.NewService(
		memory.NewUserRepository(),
		memory.NewRefreshTokenRepository(),
		auth.NewJWTService("test-secret"),
	))
}

func decodeJSONError(t *testing.T, rw *httptest.ResponseRecorder) handlers.ErrorResponse {
	t.Helper()
	var resp handlers.ErrorResponse
	if err := json.NewDecoder(rw.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error response: %v; body=%s", err, rw.Body.String())
	}
	return resp
}
