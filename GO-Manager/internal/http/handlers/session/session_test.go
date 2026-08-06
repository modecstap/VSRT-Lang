package session

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"VSRT-Lang/internal"
	"VSRT-Lang/internal/auth"
	"VSRT-Lang/internal/http/middleware"
	domain "VSRT-Lang/internal/session"
)

type stubRepository struct {
	sessions map[int]*domain.Session
}

func (r *stubRepository) Save(s *domain.Session) error {
	if r.sessions == nil {
		r.sessions = make(map[int]*domain.Session)
	}
	r.sessions[int(s.ID)] = s
	return nil
}

func (r *stubRepository) Take(sessionID int) (domain.Session, error) {
	s, ok := r.sessions[sessionID]
	if !ok {
		return domain.Session{}, errors.New("session not found")
	}
	return *s, nil
}

func TestCreateSession(t *testing.T) {
	repo := &stubRepository{}
	h := NewHandler(repo, internal.MockTranslator{})

	jwt := auth.NewJWTService("StrongSecretString")
	token, err := jwt.GenerateAccessToken("1", "demo", "demo@example.com")
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/sessions", strings.NewReader(`{"name":"demo"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rw := httptest.NewRecorder()

	middleware.Auth(jwt)(http.HandlerFunc(h.CreateSession)).ServeHTTP(rw, req)

	if rw.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rw.Code)
	}

	if len(repo.sessions) != 1 {
		t.Fatalf("expected one saved session, got %d", len(repo.sessions))
	}
}

func TestSaveAndGetRecords(t *testing.T) {
	repo := &stubRepository{}
	h := NewHandler(repo, internal.MockTranslator{})

	jwt := auth.NewJWTService("StrongSecretString")
	token, err := jwt.GenerateAccessToken("1", "demo", "demo@example.com")
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	createReq := httptest.NewRequest(http.MethodPost, "/sessions", strings.NewReader(`{"name":"demo"}`))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", "Bearer "+token)
	createRW := httptest.NewRecorder()
	middleware.Auth(jwt)(http.HandlerFunc(h.CreateSession)).ServeHTTP(createRW, createReq)

	var created struct {
		ID int64 `json:"id"`
	}
	if err := json.NewDecoder(createRW.Body).Decode(&created); err != nil {
		t.Fatalf("decode created session: %v", err)
	}

	saveReq := httptest.NewRequest(http.MethodPost, "/sessions/1/records", strings.NewReader(`{"phrase":"hello","context":"world"}`))
	saveReq.Header.Set("Content-Type", "application/json")
	saveReq.Header.Set("Authorization", "Bearer "+token)
	saveRW := httptest.NewRecorder()
	middleware.Auth(jwt)(http.HandlerFunc(h.SaveRecord)).ServeHTTP(saveRW, saveReq)

	if saveRW.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", saveRW.Code)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/sessions/1/records", nil)
	getReq.Header.Set("Authorization", "Bearer "+token)
	getRW := httptest.NewRecorder()
	middleware.Auth(jwt)(http.HandlerFunc(h.GetRecords)).ServeHTTP(getRW, getReq)

	if getRW.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", getRW.Code)
	}

	var got struct {
		Records []domain.Record `json:"records"`
	}
	if err := json.NewDecoder(getRW.Body).Decode(&got); err != nil {
		t.Fatalf("decode records: %v", err)
	}

	if len(got.Records) != 1 {
		t.Fatalf("expected one record, got %d", len(got.Records))
	}
}
