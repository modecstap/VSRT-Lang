package router

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"VSRT-Lang/internal/auth"
	"VSRT-Lang/internal/database/memory"
	"VSRT-Lang/internal/http/handlers"
	authhandler "VSRT-Lang/internal/http/handlers/auth"
	passwordresethandler "VSRT-Lang/internal/http/handlers/passwordreset"
	sessionhandler "VSRT-Lang/internal/http/handlers/session"
	userhandler "VSRT-Lang/internal/http/handlers/user"
	"VSRT-Lang/internal/http/middleware"
	"VSRT-Lang/internal/passwordreset"
	"VSRT-Lang/internal/session"
	"VSRT-Lang/internal/translators/stub"
	"VSRT-Lang/internal/user"
)

type tokenPairResponse struct {
	AccessToken  string `json:"AccessToken"`
	RefreshToken string `json:"RefreshToken"`
}

func newTestMux() *http.ServeMux {
	users := memory.NewUserRepository()
	tokens := memory.NewRefreshTokenRepository()
	sessions := memory.NewSessionRepository()
	links := memory.NewPasswordResetRepository()
	jwtSvc := auth.NewJWTService("test-secret")
	authSvc := auth.NewService(users, tokens, jwtSvc)
	sessionSvc := session.NewService(sessions, stub.Translator{})
	tx := memory.NewPasswordResetTransactor(users, links, tokens)
	resetSvc := passwordreset.NewService(
		users,
		links,
		nopMailer{},
		tx,
		"http://localhost",
		nil,
	)

	return NewServeMux(Handlers{
		Auth:           authhandler.NewAuth(authSvc),
		PasswordReset:  passwordresethandler.NewHandler(resetSvc),
		Session:        sessionhandler.NewHandler(sessionSvc),
		User:           userhandler.NewHandler(user.NewService(users), sessions),
		AuthMiddleware: middleware.Auth(jwtSvc),
	})
}

type nopMailer struct{}

func (nopMailer) Send(to, link string) error { return nil }

func doJSON(t *testing.T, mux http.Handler, method, path, token string, body any) *httptest.ResponseRecorder {
	t.Helper()

	var reader io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal body: %v", err)
		}
		reader = bytes.NewReader(raw)
	}

	req := httptest.NewRequest(method, path, reader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rw := httptest.NewRecorder()
	mux.ServeHTTP(rw, req)
	return rw
}

func TestNewServeMux_AuthAndSessionLifecycle(t *testing.T) {
	mux := newTestMux()

	register := doJSON(t, mux, http.MethodPost, "/register", "", map[string]string{
		"username": "demo",
		"email":    "demo@example.com",
		"password": "Password1",
	})
	if register.Code != http.StatusOK {
		t.Fatalf("register status = %d, body=%s", register.Code, register.Body.String())
	}

	login := doJSON(t, mux, http.MethodPost, "/login", "", map[string]string{
		"login":    "demo",
		"password": "Password1",
	})
	if login.Code != http.StatusOK {
		t.Fatalf("login status = %d, body=%s", login.Code, login.Body.String())
	}

	var tokens tokenPairResponse
	if err := json.NewDecoder(login.Body).Decode(&tokens); err != nil {
		t.Fatalf("decode tokens: %v", err)
	}
	if tokens.AccessToken == "" {
		t.Fatal("expected access token")
	}

	for _, path := range []string{"/users/me", "/users/me?id=not-demo"} {
		me := doJSON(t, mux, http.MethodGet, path, tokens.AccessToken, nil)
		var got struct {
			Username string  `json:"username"`
			Email    string  `json:"email"`
			Avatar   *string `json:"avatar"`
		}
		if me.Code != http.StatusOK || json.Unmarshal(me.Body.Bytes(), &got) != nil ||
			got.Username != "demo" || got.Email != "demo@example.com" || got.Avatar != nil {
			t.Fatalf("%s status=%d body=%s", path, me.Code, me.Body.String())
		}
	}
	if notMe := doJSON(t, mux, http.MethodGet, "/users/not-me", tokens.AccessToken, nil); notMe.Code != http.StatusNotFound {
		t.Fatalf("not-me status = %d, want 404; body=%s", notMe.Code, notMe.Body.String())
	}

	created := doJSON(t, mux, http.MethodPost, "/sessions", tokens.AccessToken, map[string]string{
		"name": "daily",
	})
	if created.Code != http.StatusCreated {
		t.Fatalf("create session status = %d, body=%s", created.Code, created.Body.String())
	}

	var sessionResp struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}
	if err := json.NewDecoder(created.Body).Decode(&sessionResp); err != nil {
		t.Fatalf("decode session: %v", err)
	}
	if sessionResp.ID == 0 || sessionResp.Name != "daily" {
		t.Fatalf("session = %+v", sessionResp)
	}

	saved := doJSON(t, mux, http.MethodPost, "/sessions/"+strconv.FormatInt(sessionResp.ID, 10)+"/records", tokens.AccessToken, map[string]string{
		"phrase":  "hello",
		"context": "world",
	})
	if saved.Code != http.StatusCreated {
		t.Fatalf("save record status = %d, body=%s", saved.Code, saved.Body.String())
	}

	records := doJSON(t, mux, http.MethodGet, "/sessions/"+strconv.FormatInt(sessionResp.ID, 10)+"/records", tokens.AccessToken, nil)
	if records.Code != http.StatusOK {
		t.Fatalf("get records status = %d, body=%s", records.Code, records.Body.String())
	}

	var recordsResp struct {
		Records []session.Record `json:"records"`
	}
	if err := json.NewDecoder(records.Body).Decode(&recordsResp); err != nil {
		t.Fatalf("decode records: %v", err)
	}
	if len(recordsResp.Records) != 1 || recordsResp.Records[0].Phrase != "hello" {
		t.Fatalf("records = %+v", recordsResp.Records)
	}

	userSessions := doJSON(t, mux, http.MethodGet, "/users/sessions", tokens.AccessToken, nil)
	if userSessions.Code != http.StatusOK {
		t.Fatalf("user sessions status = %d, body=%s", userSessions.Code, userSessions.Body.String())
	}

	var listed []session.Session
	if err := json.NewDecoder(userSessions.Body).Decode(&listed); err != nil {
		t.Fatalf("decode user sessions: %v", err)
	}
	if len(listed) != 1 || listed[0].ID != sessionResp.ID {
		t.Fatalf("listed sessions = %+v", listed)
	}

	deleted := doJSON(t, mux, http.MethodDelete, "/sessions/"+strconv.FormatInt(sessionResp.ID, 10), tokens.AccessToken, nil)
	if deleted.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d, body=%s", deleted.Code, deleted.Body.String())
	}

	missing := doJSON(t, mux, http.MethodGet, "/sessions/"+strconv.FormatInt(sessionResp.ID, 10)+"/records", tokens.AccessToken, nil)
	if missing.Code != http.StatusNotFound {
		t.Fatalf("deleted session get status = %d, want 404; body=%s", missing.Code, missing.Body.String())
	}
}

func TestNewServeMux_DeleteRecord(t *testing.T) {
	mux := newTestMux()

	demo := loginUser(t, mux, "demo", "demo@example.com", "Password1")
	sessionID := createSession(t, mux, demo)
	saveRecord(t, mux, demo, sessionID, "hello")
	saveRecord(t, mux, demo, sessionID, "world")

	other := loginUser(t, mux, "other", "other@example.com", "Password1")
	foreign := doJSON(t, mux, http.MethodDelete, "/sessions/"+sessionID+"/records/hello", other, nil)
	if foreign.Code != http.StatusNotFound {
		t.Fatalf("foreign delete status = %d, want 404; body=%s", foreign.Code, foreign.Body.String())
	}
	var foreignErr handlers.ErrorResponse
	if err := json.NewDecoder(foreign.Body).Decode(&foreignErr); err != nil {
		t.Fatalf("decode foreign delete: %v", err)
	}
	if foreignErr.Error.Code != "session_not_found" {
		t.Fatalf("foreign delete code = %q, want session_not_found", foreignErr.Error.Code)
	}
	if phrases := recordPhrases(t, mux, demo, sessionID); !containsPhrase(phrases, "hello") {
		t.Fatalf("phrases after foreign delete = %v, want hello", phrases)
	}

	deleted := doJSON(t, mux, http.MethodDelete, "/sessions/"+sessionID+"/records/hello", demo, nil)
	if deleted.Code != http.StatusNoContent {
		t.Fatalf("delete record status = %d, want 204; body=%s", deleted.Code, deleted.Body.String())
	}
	phrases := recordPhrases(t, mux, demo, sessionID)
	if containsPhrase(phrases, "hello") || !containsPhrase(phrases, "world") {
		t.Fatalf("phrases after delete = %v, want world only", phrases)
	}

	missingPhrase := doJSON(t, mux, http.MethodDelete, "/sessions/"+sessionID+"/records", demo, nil)
	if missingPhrase.Code != http.StatusBadRequest {
		t.Fatalf("missing phrase status = %d, want 400; body=%s", missingPhrase.Code, missingPhrase.Body.String())
	}
	if phrases = recordPhrases(t, mux, demo, sessionID); !containsPhrase(phrases, "world") {
		t.Fatalf("phrases after empty delete = %v, want world", phrases)
	}

	extra := doJSON(t, mux, http.MethodDelete, "/sessions/"+sessionID+"/extra", demo, nil)
	if extra.Code != http.StatusMethodNotAllowed {
		t.Fatalf("extra path status = %d, want 405; body=%s", extra.Code, extra.Body.String())
	}
	if phrases = recordPhrases(t, mux, demo, sessionID); !containsPhrase(phrases, "world") {
		t.Fatalf("phrases after extra path = %v, want world", phrases)
	}
}

func loginUser(t *testing.T, mux http.Handler, username, email, password string) string {
	t.Helper()

	register := doJSON(t, mux, http.MethodPost, "/register", "", map[string]string{
		"username": username,
		"email":    email,
		"password": password,
	})
	if register.Code != http.StatusOK {
		t.Fatalf("register %s status = %d, body=%s", username, register.Code, register.Body.String())
	}

	login := doJSON(t, mux, http.MethodPost, "/login", "", map[string]string{
		"login":    username,
		"password": password,
	})
	if login.Code != http.StatusOK {
		t.Fatalf("login %s status = %d, body=%s", username, login.Code, login.Body.String())
	}

	var tokens tokenPairResponse
	if err := json.NewDecoder(login.Body).Decode(&tokens); err != nil {
		t.Fatalf("decode tokens: %v", err)
	}
	if tokens.AccessToken == "" {
		t.Fatal("expected access token")
	}
	return tokens.AccessToken
}

func createSession(t *testing.T, mux http.Handler, token string) string {
	t.Helper()

	created := doJSON(t, mux, http.MethodPost, "/sessions", token, map[string]string{"name": "daily"})
	if created.Code != http.StatusCreated {
		t.Fatalf("create session status = %d, body=%s", created.Code, created.Body.String())
	}

	var sessionResp struct {
		ID int64 `json:"id"`
	}
	if err := json.NewDecoder(created.Body).Decode(&sessionResp); err != nil {
		t.Fatalf("decode session: %v", err)
	}
	if sessionResp.ID == 0 {
		t.Fatal("expected session id")
	}
	return strconv.FormatInt(sessionResp.ID, 10)
}

func saveRecord(t *testing.T, mux http.Handler, token, sessionID, phrase string) {
	t.Helper()

	saved := doJSON(t, mux, http.MethodPost, "/sessions/"+sessionID+"/records", token, map[string]string{
		"phrase":  phrase,
		"context": "sample",
	})
	if saved.Code != http.StatusCreated {
		t.Fatalf("save %s status = %d, body=%s", phrase, saved.Code, saved.Body.String())
	}
}

func recordPhrases(t *testing.T, mux http.Handler, token, sessionID string) []string {
	t.Helper()

	records := doJSON(t, mux, http.MethodGet, "/sessions/"+sessionID+"/records", token, nil)
	if records.Code != http.StatusOK {
		t.Fatalf("get records status = %d, want 200; body=%s", records.Code, records.Body.String())
	}

	var recordsResp struct {
		Records []session.Record `json:"records"`
	}
	if err := json.NewDecoder(records.Body).Decode(&recordsResp); err != nil {
		t.Fatalf("decode records: %v", err)
	}

	phrases := make([]string, 0, len(recordsResp.Records))
	for _, record := range recordsResp.Records {
		phrases = append(phrases, record.Phrase)
	}
	return phrases
}

func containsPhrase(phrases []string, phrase string) bool {
	for _, item := range phrases {
		if item == phrase {
			return true
		}
	}
	return false
}

func TestNewServeMux_ProtectedRoutesRequireAuth(t *testing.T) {
	t.Parallel()

	mux := newTestMux()

	tests := []struct {
		name   string
		method string
		path   string
		body   any
	}{
		{name: "create session", method: http.MethodPost, path: "/sessions", body: map[string]string{"name": "demo"}},
		{name: "delete session", method: http.MethodDelete, path: "/sessions/1"},
		{name: "delete record", method: http.MethodDelete, path: "/sessions/1/records/hello"},
		{name: "save record", method: http.MethodPost, path: "/sessions/1/records", body: map[string]string{"phrase": "hello", "context": "world"}},
		{name: "get records", method: http.MethodGet, path: "/sessions/1/records"},
		{name: "user sessions", method: http.MethodGet, path: "/users/sessions"},
		{name: "current user", method: http.MethodGet, path: "/users/me"},
		{name: "save avatar", method: http.MethodPost, path: "/users/avatar"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rw := doJSON(t, mux, tt.method, tt.path, "", tt.body)
			if rw.Code != http.StatusUnauthorized {
				t.Fatalf("status = %d, want 401; body=%s", rw.Code, rw.Body.String())
			}
		})
	}

	for _, tt := range []struct{ header, want string }{
		{"", "unauthorized\n"},
		{"Token abc", "invalid authorization header\n"},
		{"Bearer not-a-token", "invalid token\n"},
	} {
		req := httptest.NewRequest(http.MethodGet, "/users/me", nil)
		if tt.header != "" {
			req.Header.Set("Authorization", tt.header)
		}
		rw := httptest.NewRecorder()
		mux.ServeHTTP(rw, req)
		if rw.Code != http.StatusUnauthorized || rw.Body.String() != tt.want {
			t.Fatalf("header %q status=%d body=%q, want 401 %q", tt.header, rw.Code, rw.Body.String(), tt.want)
		}
	}
}

func TestNewServeMux_UnknownRoute(t *testing.T) {
	t.Parallel()

	mux := newTestMux()
	rw := doJSON(t, mux, http.MethodGet, "/does-not-exist", "", nil)
	if rw.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rw.Code)
	}
}

func TestNewServeMux_NilOptionalHandlers(t *testing.T) {
	t.Parallel()

	mux := NewServeMux(Handlers{
		Auth: authhandler.NewAuth(auth.NewService(
			memory.NewUserRepository(),
			memory.NewRefreshTokenRepository(),
			auth.NewJWTService("test-secret"),
		)),
	})

	rw := doJSON(t, mux, http.MethodPost, "/sessions", "", map[string]string{"name": "demo"})
	if rw.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404 when session handler is nil", rw.Code)
	}

	rw = doJSON(t, mux, http.MethodGet, "/users/sessions", "", nil)
	if rw.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404 when user handler is nil", rw.Code)
	}
}

func TestNewServeMux_LoginRejectsUnknownUser(t *testing.T) {
	t.Parallel()

	mux := newTestMux()
	rw := doJSON(t, mux, http.MethodPost, "/login", "", map[string]string{
		"login":    "missing",
		"password": "Password1",
	})
	if rw.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401; body=%s", rw.Code, rw.Body.String())
	}

	var resp handlers.ErrorResponse
	if err := json.NewDecoder(rw.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp.Error.Code != "unauthorized" {
		t.Fatalf("error code = %q, want unauthorized", resp.Error.Code)
	}
}
