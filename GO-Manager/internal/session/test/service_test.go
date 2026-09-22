package test

import (
	"errors"
	"testing"

	"VSRT-Lang/internal/database/memory"
	"VSRT-Lang/internal/session"
	"VSRT-Lang/internal/translators/stub"
	"VSRT-Lang/internal/user"
)

const (
	userAnna  = user.UserId("anna")
	userOther = user.UserId("other")
)

var errServiceUnavailable = errors.New("service unavailable")

func setupService() *session.Service {
	return session.NewService(memory.NewSessionRepository(), stub.Translator{})
}

func mustCreateSession(t *testing.T, svc *session.Service, uid user.UserId, name string) int64 {
	t.Helper()
	id, err := svc.NewSession(uid, name)
	if err != nil {
		t.Fatalf("NewSession(%q): %v", name, err)
	}
	return id
}

func mustAddRecord(t *testing.T, svc *session.Service, uid user.UserId, sessionID int64, phrase, context string) session.Record {
	t.Helper()
	record, err := svc.AddRecord(session.AddRecordCommand{
		UserId:    uid,
		SessionId: sessionID,
		Phrase:    phrase,
		Context:   context,
	})
	if err != nil {
		t.Fatalf("AddRecord(%q): %v", phrase, err)
	}
	return record
}

func mustGetSession(t *testing.T, svc *session.Service, uid user.UserId, sessionID int64) session.Session {
	t.Helper()
	got, err := svc.GetSession(uid, sessionID)
	if err != nil {
		t.Fatalf("GetSession: %v", err)
	}
	return got
}

func recordsOf(t *testing.T, svc *session.Service, uid user.UserId, sessionID int64) []session.Record {
	t.Helper()
	got := mustGetSession(t, svc, uid, sessionID)
	return got.GetRecords()
}

func containsSession(sessions []session.Session, name string) bool {
	for _, s := range sessions {
		if s.Name == name {
			return true
		}
	}
	return false
}

type failSaveRepo struct {
	*memory.SessionRepository
	err error
}

func (r failSaveRepo) Save(*session.Session) (int64, error) {
	return 0, r.err
}

func TestCreateSessionWithName(t *testing.T) {
	svc := setupService()
	id := mustCreateSession(t, svc, userAnna, "Урок 1")

	got := mustGetSession(t, svc, userAnna, id)
	if got.Name != "Урок 1" {
		t.Fatalf("session name = %q, want %q", got.Name, "Урок 1")
	}
}

func TestCreateSessionWithoutName(t *testing.T) {
	svc := setupService()

	_, err := svc.NewSession(userAnna, "")
	if !errors.Is(err, session.ErrSessionNameRequired) {
		t.Fatalf("error = %v, want %v", err, session.ErrSessionNameRequired)
	}

	sessions, err := svc.GetSessions(userAnna)
	if err != nil {
		t.Fatalf("GetSessions: %v", err)
	}
	if len(sessions) != 0 {
		t.Fatalf("got %d sessions, want 0", len(sessions))
	}
}

func TestListOwnSessions(t *testing.T) {
	svc := setupService()
	mustCreateSession(t, svc, userAnna, "Урок 1")
	mustCreateSession(t, svc, userAnna, "Урок 2")
	mustCreateSession(t, svc, userOther, "Чужой урок")

	sessions, err := svc.GetSessions(userAnna)
	if err != nil {
		t.Fatalf("GetSessions: %v", err)
	}
	if !containsSession(sessions, "Урок 1") || !containsSession(sessions, "Урок 2") {
		t.Fatalf("list missing own sessions: %+v", sessions)
	}
	if containsSession(sessions, "Чужой урок") {
		t.Fatal("list contains foreign session")
	}
}

func TestListSessionsEmpty(t *testing.T) {
	svc := setupService()

	sessions, err := svc.GetSessions(userAnna)
	if err != nil {
		t.Fatalf("GetSessions: %v", err)
	}
	if len(sessions) != 0 {
		t.Fatalf("got %d sessions, want 0", len(sessions))
	}
}

func TestGetOwnSession(t *testing.T) {
	svc := setupService()
	id := mustCreateSession(t, svc, userAnna, "Урок 1")

	got := mustGetSession(t, svc, userAnna, id)
	if got.Name != "Урок 1" {
		t.Fatalf("session name = %q, want %q", got.Name, "Урок 1")
	}
}

func TestGetForeignSessionNotFound(t *testing.T) {
	svc := setupService()
	id := mustCreateSession(t, svc, userOther, "Чужой урок")

	_, err := svc.GetSession(userAnna, id)
	if !errors.Is(err, session.ErrSessionNotFound) {
		t.Fatalf("error = %v, want %v", err, session.ErrSessionNotFound)
	}
}

func TestDeleteOwnSession(t *testing.T) {
	svc := setupService()
	id := mustCreateSession(t, svc, userAnna, "Урок 1")
	mustAddRecord(t, svc, userAnna, id, "hello", "")

	if err := svc.DeleteSession(userAnna, id); err != nil {
		t.Fatalf("DeleteSession: %v", err)
	}

	sessions, err := svc.GetSessions(userAnna)
	if err != nil {
		t.Fatalf("GetSessions: %v", err)
	}
	if containsSession(sessions, "Урок 1") {
		t.Fatal("deleted session still in list")
	}

	_, err = svc.GetSession(userAnna, id)
	if !errors.Is(err, session.ErrSessionNotFound) {
		t.Fatalf("session records still available, error = %v", err)
	}
}

func TestDeleteForeignSessionNotFound(t *testing.T) {
	svc := setupService()
	id := mustCreateSession(t, svc, userOther, "Чужой урок")

	err := svc.DeleteSession(userAnna, id)
	if !errors.Is(err, session.ErrSessionNotFound) {
		t.Fatalf("error = %v, want %v", err, session.ErrSessionNotFound)
	}

	if _, err := svc.GetSession(userOther, id); err != nil {
		t.Fatalf("foreign session was deleted: %v", err)
	}
}

func TestCreateSessionWhenRepositoryUnavailable(t *testing.T) {
	repo := memory.NewSessionRepository()
	svc := session.NewService(failSaveRepo{SessionRepository: repo, err: errServiceUnavailable}, stub.Translator{})

	_, err := svc.NewSession(userAnna, "Урок 1")
	if err == nil {
		t.Fatal("expected service error")
	}

	sessions, err := repo.FindByUser(userAnna)
	if err != nil {
		t.Fatalf("FindByUser: %v", err)
	}
	if len(sessions) != 0 {
		t.Fatalf("got %d sessions, want 0", len(sessions))
	}
}
