package test

import (
	"VSRT-Lang/internal/database/memory"
	"VSRT-Lang/internal/session"
	"VSRT-Lang/internal/translators/stub"
	"VSRT-Lang/internal/user"
	"errors"
	"testing"
)

const ANOTHER_USER = user.UserId("Another User")

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

func TestGetSession(t *testing.T) {
	service, _ := setupService()
	userID, sessionID, err := createSession(service)
	if err != nil {
		t.Fatalf("expected session to be created, but got error: %v", err)
	}

	got, err := service.GetSession(userID, sessionID)
	if err != nil {
		t.Fatalf("expected to get session, but got error: %v", err)
	}

	if got.ID != sessionID {
		t.Fatalf("expected session ID to be %d, but got %d", sessionID, got.ID)
	}
	if got.Name != "test-session" {
		t.Fatalf("expected session name to be %q, but got %q", "test-session", got.Name)
	}
	if got.User != userID {
		t.Fatalf("expected session user to be %q, but got %q", userID, got.User)
	}
}

func TestGetSessionByAnotherUser(t *testing.T) {
	service, _ := setupService()
	_, sessionID, err := createSession(service)
	if err != nil {
		t.Fatalf("expected session to be created, but got error: %v", err)
	}

	_, err = service.GetSession(ANOTHER_USER, sessionID)
	if !errors.Is(err, session.ErrSessionNotFound) {
		t.Fatalf("expected error %q, but got %v", session.ErrSessionNotFound, err)
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

func TestAddRecordByAnotherUser(t *testing.T) {
	service, _ := setupService()
	userId, sessionId, err := createSession(service)
	if err != nil {
		t.Fatalf("expected to get sessions, but got error: %v", err)
	}

	command := session.AddRecordCommand{
		UserId:    ANOTHER_USER,
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

	err = service.DeleteSession(userId, sessionId)

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

	err = service.DeleteSession(ANOTHER_USER, sessionId)
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

func TestDeleteRecord(t *testing.T) {
	service, _ := setupService()
	userId, sessionId, err := createSession(service)
	service.AddRecord(session.AddRecordCommand{
		UserId:    userId,
		SessionId: sessionId,
		Phrase:    "test",
		Context:   "test context",
	})

	err = service.DeleteRecord(userId, sessionId, "test")
	if err != nil {
		t.Fatalf("expected to delete Record, but got error: %v", err)
	}

	sessions, err := service.GetSessions(userId)
	if err != nil {
		t.Fatalf("expected to get sessions, but got error: %v", err)
	}
	_, ok := sessions[0].Records["test"]
	if ok {
		t.Fatal("record must be deleted but sill exist")
	}
}

func TestDeleteRecordByAnotherUser(t *testing.T) {
	service, _ := setupService()
	userId, sessionId, err := createSession(service)
	service.AddRecord(session.AddRecordCommand{
		UserId:    userId,
		SessionId: sessionId,
		Phrase:    "test",
		Context:   "test context",
	})

	err = service.DeleteRecord(ANOTHER_USER, sessionId, "test")

	sessions, err := service.GetSessions(userId)
	if err != nil {
		t.Fatalf("expected to get sessions, but got error: %v", err)
	}
	_, ok := sessions[0].Records["test"]
	if !ok {
		t.Fatal("record must exist")
	}
}

func TestAddRecordWhithEmptyContext(t *testing.T) {
	service, _ := setupService()
	userId, sessionId, err := createSession(service)
	newRecord := session.AddRecordCommand{
		UserId:    userId,
		SessionId: sessionId,
		Phrase:    "test",
		Context:   "",
	}

	record, err := service.AddRecord(newRecord)
	if err != nil {
		t.Fatalf("expected to add record, but got error: %v", err)
	}
	for _, c := range record.Contexts {
		if c.Phrase == "" {
			t.Fatalf("empty context must be not added")
		}
	}
}

func TestAddRecordWhithExistingContext(t *testing.T) {
	service, _ := setupService()
	userId, sessionId, err := createSession(service)
	newRecord := session.AddRecordCommand{
		UserId:    userId,
		SessionId: sessionId,
		Phrase:    "test",
		Context:   "test context",
	}

	record, err := service.AddRecord(newRecord)
	if err != nil {
		t.Fatalf("expected to add record, but got error: %v", err)
	}
	firstContextCount := len(record.Contexts)

	record, err = service.AddRecord(newRecord)
	if err != nil {
		t.Fatalf("expected to add record, but got error: %v", err)
	}
	secondContextCount := len(record.Contexts)

	if firstContextCount != secondContextCount {
		t.Fatalf("context count must be %d but got %d", firstContextCount, secondContextCount)
	}
}
