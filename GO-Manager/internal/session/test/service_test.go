package test

import (
	"VSRT-Lang/internal/database/memory"
	"VSRT-Lang/internal/session"
	"VSRT-Lang/internal/translators/stub"
	"VSRT-Lang/internal/user"
	"errors"
	"testing"
)

func setupService() (*session.Service, *memory.SessionRepository) {
	repo := memory.NewSessionRepository()
	trans := stub.Translator{}
	service := session.NewService(repo, trans)
	return service, repo
}

func createSession(service *session.Service) (user.UserId, int64, error) {
	userId := user.UserId("test-user")
	sessionName := "test-session"
	sessionId, err := service.NewSession(userId, sessionName)
	return userId, sessionId, err
}

func TestCreateSession(t *testing.T) {
	service, repo := setupService()

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
	service, _ := setupService()
	_, _, err := createSession(service)
	if err != nil {
		t.Fatalf("expected to get sessions, but got error: %v", err)
	}

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
	service, _ := setupService()
	userId, sessionId, err := createSession(service)
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

func TestAddRecordByanotherUser(t *testing.T) {
	service, _ := setupService()
	userId, sessionId, err := createSession(service)
	if err != nil {
		t.Fatalf("expected to get sessions, but got error: %v", err)
	}

	command := session.AddRecordCommand{
		UserId:    "Another User",
		SessionId: sessionId,
		Phrase:    "test",
		Context:   "test context",
	}

	_, err = service.AddRecord(command)

	if !errors.Is(err, session.ErrSessionNotFound) {
		t.Fatalf(
			"must be error «session not found» but got %v", err,
		)
	}

	sessions, err := service.GetSessions(userId)
	if err != nil {
		t.Fatalf("expected to get sessions, but got error: %v", err)
	}
	s := sessions[0]
	if len(s.Records) != 0 {
		t.Fatalf("record count be 0 but got %d", len(s.Records))
	}
}

func TestDeleteSession(t *testing.T) {
	service, _ := setupService()
	userId, sessionId, err := createSession(service)
	if err != nil {
		t.Fatalf("expected to get sessions, but got error: %v", err)
	}

	err = service.Delete(userId, sessionId)

	sessions, err := service.GetSessions(userId)
	if len(sessions) != 0 {
		t.Fatalf(
			"count session must be 0 but got %d", len(sessions),
		)
	}
}

func TestDeleteAnotherUserSession(t *testing.T) {
	service, _ := setupService()
	userId, sessionId, err := createSession(service)
	if err != nil {
		t.Fatalf("expected to get sessions, but got error: %v", err)
	}

	err = service.Delete("WrongUser", sessionId)
	if !errors.Is(err, session.ErrUnauthorized) {
		t.Fatalf(
			"must be error «unauthorized» but got %v", err,
		)
	}

	sessions, err := service.GetSessions(userId)
	if len(sessions) != 1 {
		t.Fatalf(
			"count session must be 1 but got %d", len(sessions),
		)
	}
}
