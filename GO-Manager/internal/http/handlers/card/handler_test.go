package card

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"VSRT-Lang/internal/auth"
	domain "VSRT-Lang/internal/card"
	"VSRT-Lang/internal/http/handlers"
	"VSRT-Lang/internal/http/middleware"
	"VSRT-Lang/internal/session"
	"VSRT-Lang/internal/user"
)

const testSecret = "test-secret"
const testUserID = "1"

type fakeService struct {
	getCards        func(user.UserId, int64, int) ([]session.Record, error)
	updateKnowledge func(domain.UpdateKnowledgeCommand) error
}

func (f *fakeService) GetCards(userID user.UserId, sessionID int64, limit int) ([]session.Record, error) {
	if f.getCards == nil {
		return nil, errors.New("unexpected GetCards")
	}
	return f.getCards(userID, sessionID, limit)
}

func (f *fakeService) UpdateKnowledge(cmd domain.UpdateKnowledgeCommand) error {
	if f.updateKnowledge == nil {
		return errors.New("unexpected UpdateKnowledge")
	}
	return f.updateKnowledge(cmd)
}

func testJWT(t *testing.T) *auth.JWTService {
	t.Helper()
	return auth.NewJWTService(testSecret)
}

func serveAuthed(t *testing.T, jwtSvc *auth.JWTService, userID string, handler http.HandlerFunc, req *http.Request) *httptest.ResponseRecorder {
	t.Helper()
	token, err := jwtSvc.GenerateAccessToken(userID, "demo", "demo@example.com")
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	rw := httptest.NewRecorder()
	middleware.Auth(jwtSvc)(http.HandlerFunc(handler)).ServeHTTP(rw, req)
	return rw
}

func decodeJSONError(t *testing.T, rw *httptest.ResponseRecorder) handlers.ErrorResponse {
	t.Helper()
	var resp handlers.ErrorResponse
	if err := json.NewDecoder(rw.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error response: %v; body=%s", err, rw.Body.String())
	}
	return resp
}
