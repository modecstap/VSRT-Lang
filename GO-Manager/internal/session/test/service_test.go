package test

import (
	"VSRT-Lang/internal/database/memory"
	"VSRT-Lang/internal/session"
	"VSRT-Lang/internal/translators/stub"
	"VSRT-Lang/internal/user"
	"testing"
)

func TestCreateSession(t *testing.T) {
	repo := memory.NewSessionRepository()
	trans := stub.Translator{}
	service := session.NewService(repo, trans)

	userID := user.UserId("test-user")
	sessionName := "test-session"

	sessionId, err := service.NewSession(userID, sessionName)
	if err != nil {
		t.Fatalf("expected session to be created, but got error: %v", err)
	}
	loadedSession, err := repo.Take(sessionId)
	if err != nil {
		t.Fatalf("expected session to be created, but got error: %v", err)
	}

	if loadedSession.Name != sessionName {
		t.Fatalf("expected session name to be %q, but got %q", sessionName, loadedSession.Name)
	}
	if loadedSession.User != userID {
		t.Fatalf("expected session user to be %q, but got %q", userID, loadedSession.User)
	}
}

func TestGetSessions(t *testing.T) {
	repo := memory.NewSessionRepository()
	trans := stub.Translator{}
	service := session.NewService(repo, trans)

	service.NewSession(user.UserId("test-user"), "test-session")

	sessions, err := service.GetSessions(user.UserId("test-user"))
	if err != nil {
		t.Fatalf("expected to get sessions, but got error: %v", err)
	}

	if len(sessions) != 1 {
		t.Fatalf("expected to get 1 session, but got %d", len(sessions))
	}
	if sessions[0].Name != "test-session" {
		t.Fatalf("expected session name to be %q, but got %q", "test-session", sessions[0].Name)
	}
}

func TestAddRecord(t *testing.T) {
	repo := memory.NewSessionRepository()
	trans := stub.Translator{}
	service := session.NewService(repo, trans)
	userId := user.UserId("test-user")
	sessionName := "test-session"
	sessionId, err := service.NewSession(userId, sessionName)
	if err != nil {
		t.Fatalf("expected to get sessions, but got error: %v", err)
	}

	command := session.AddRecordCommand{
		UserId:    userId,
		SessionId: sessionId,
		Phrase:    "test",
		Context:   "test context",
	}
	record, err := service.AddRecord(command)

	sessions, err := service.GetSessions(userId)
	if err != nil {
		t.Fatalf("expected to get sessions, but got error: %v", err)
	}

	s := sessions[0]
	
	sRecord, ok := s.Records["test"]
	if !ok {
		t.Fatalf("record must be in saved session")
	}
	if sRecord.Phrase != record.Phrase {
		t.Fatalf(
			"expected equality record in session and responce record, but got %q and %q",
			sRecord.Phrase, record.Phrase,
		)
	}
}
