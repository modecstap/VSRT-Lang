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
	service := session.NewService(repo,trans)

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

