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
	sessionhandler "VSRT-Lang/internal/http/handlers/session"
	userhandler "VSRT-Lang/internal/http/handlers/user"
	"VSRT-Lang/internal/http/middleware"
	"VSRT-Lang/internal/session"
	"VSRT-Lang/internal/translators/stub"
)

type tokenPairResponse struct {
	AccessToken  string `json:"AccessToken"`
	RefreshToken string `json:"RefreshToken"`
}

func newTestMux() *http.ServeMux {
	users := memory.NewUserRepository()
	tokens := memory.NewRefreshTokenRepository()
	sessions := memory.NewSessionRepository()
	jwtSvc := auth.NewJWTService("test-secret")
	authSvc := auth.NewService(users, tokens, jwtSvc)
	sessionSvc := session.NewService(sessions, stub.Translator{})

	return NewServeMux(Handlers{
		Auth:           authhandler.NewAuth(authSvc),
		Session:        sessionhandler.NewHandler(sessionSvc),
		User:           userhandler.NewHandler(users, sessions),
		AuthMiddleware: middleware.Auth(jwtSvc),
	})
}

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
		{name: "save record", method: http.MethodPost, path: "/sessions/1/records", body: map[string]string{"phrase": "hello", "context": "world"}},
		{name: "get records", method: http.MethodGet, path: "/sessions/1/records"},
		{name: "user sessions", method: http.MethodGet, path: "/users/sessions"},
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
