package session

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"VSRT-Lang/internal/auth"
	"VSRT-Lang/internal/http/handlers"
	"VSRT-Lang/internal/http/middleware"
	domain "VSRT-Lang/internal/session"
	"VSRT-Lang/internal/user"
)

const testSecret = "test-secret"
const testUserID = "1"

type fakeService struct {
	newSession    func(user.UserId, string) (int64, error)
	getSession    func(user.UserId, int64) (domain.Session, error)
	deleteSession func(user.UserId, int64) error
	deleteRecord  func(user.UserId, int64, string) error
	addRecord     func(domain.AddRecordCommand) (domain.Record, error)
}

func (f *fakeService) NewSession(userID user.UserId, name string) (int64, error) {
	if f.newSession == nil {
		return 0, errors.New("unexpected NewSession")
	}
	return f.newSession(userID, name)
}

func (f *fakeService) GetSession(userID user.UserId, sessionID int64) (domain.Session, error) {
	if f.getSession == nil {
		return domain.Session{}, errors.New("unexpected GetSession")
	}
	return f.getSession(userID, sessionID)
}

func (f *fakeService) DeleteSession(userID user.UserId, sessionID int64) error {
	if f.deleteSession == nil {
		return errors.New("unexpected DeleteSession")
	}
	return f.deleteSession(userID, sessionID)
}

func (f *fakeService) DeleteRecord(userID user.UserId, sessionID int64, phrase string) error {
	if f.deleteRecord == nil {
		return errors.New("unexpected DeleteRecord")
	}
	return f.deleteRecord(userID, sessionID, phrase)
}

func (f *fakeService) AddRecord(cmd domain.AddRecordCommand) (domain.Record, error) {
	if f.addRecord == nil {
		return domain.Record{}, errors.New("unexpected AddRecord")
	}
	return f.addRecord(cmd)
}

func testJWT(t *testing.T) *auth.JWTService {
	t.Helper()
	return auth.NewJWTService(testSecret)
}

func bearerToken(t *testing.T, jwtSvc *auth.JWTService, userID string) string {
	t.Helper()
	token, err := jwtSvc.GenerateAccessToken(userID, "demo", "demo@example.com")
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}
	return token
}

func serveAuthed(t *testing.T, jwtSvc *auth.JWTService, userID string, handler http.HandlerFunc, req *http.Request) *httptest.ResponseRecorder {
	t.Helper()
	req.Header.Set("Authorization", "Bearer "+bearerToken(t, jwtSvc, userID))
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
