package memory

import (
	"bytes"
	"errors"
	"testing"

	"VSRT-Lang/internal/user"
)

func TestUserRepository_Save_ProfileKeys(t *testing.T) {
	t.Parallel()

	repo := NewUserRepository()
	alice := user.NewUser("alice", "alice@example.com", "secret")
	bob := user.NewUser("bob", "bob@example.com", "secret")
	if err := repo.Create(alice); err != nil {
		t.Fatalf("Create alice: %v", err)
	}
	if err := repo.Create(bob); err != nil {
		t.Fatalf("Create bob: %v", err)
	}

	same, err := repo.FindByID(alice.ID)
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}
	if err := repo.Save(same); err != nil {
		t.Fatalf("Save same keys: %v", err)
	}

	takenName, err := repo.FindByID(alice.ID)
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}
	takenName.Username = bob.Username
	if err := repo.Save(takenName); !errors.Is(err, user.ErrUsernameTaken) {
		t.Fatalf("Save taken username: %v", err)
	}
	if !profileKeysIntact(t, repo, alice, bob) {
		t.Fatal("username conflict changed stored keys")
	}

	takenMail, err := repo.FindByID(alice.ID)
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}
	takenMail.Email = bob.Email
	if err := repo.Save(takenMail); !errors.Is(err, user.ErrEmailTaken) {
		t.Fatalf("Save taken email: %v", err)
	}
	if !profileKeysIntact(t, repo, alice, bob) {
		t.Fatal("email conflict changed stored keys")
	}

	missing := user.NewUser("nope", "nope@example.com", "secret")
	missing.ID = "missing"
	if err := repo.Save(missing); !errors.Is(err, user.ErrNotFound) {
		t.Fatalf("Save missing: %v", err)
	}
}

func profileKeysIntact(t *testing.T, repo *UserRepository, alice, bob *user.User) bool {
	t.Helper()
	gotAlice, err := repo.FindByUsername(alice.Username)
	if err != nil || gotAlice.ID != alice.ID || gotAlice.Email != alice.Email {
		return false
	}
	gotBob, err := repo.FindByUsername(bob.Username)
	if err != nil || gotBob.ID != bob.ID || gotBob.Email != bob.Email {
		return false
	}
	byAliceMail, err := repo.FindByEmail(alice.Email)
	if err != nil || byAliceMail.ID != alice.ID {
		return false
	}
	byBobMail, err := repo.FindByEmail(bob.Email)
	return err == nil && byBobMail.ID == bob.ID
}

func TestUserRepository_SaveAndGetAvatar(t *testing.T) {
	t.Parallel()

	repo := NewUserRepository()
	u := user.NewUser("demo", "demo@example.com", "secret")
	if err := repo.Create(u); err != nil {
		t.Fatalf("Create: %v", err)
	}

	first := user.Avatar{Bytes: []byte{1, 2, 3}, MediaType: "image/png"}
	if err := repo.SaveAvatar(user.UserId(u.ID), first); err != nil {
		t.Fatalf("SaveAvatar: %v", err)
	}
	first.Bytes[0] = 9

	got, err := repo.FindByID(u.ID)
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}
	if got.Avatar.MediaType != "image/png" || !bytes.Equal(got.Avatar.Bytes, []byte{1, 2, 3}) {
		t.Fatalf("got %+v, want copied first avatar", got.Avatar)
	}

	second := user.Avatar{Bytes: []byte{4, 5}, MediaType: "image/png"}
	if err := repo.SaveAvatar(user.UserId(u.ID), second); err != nil {
		t.Fatalf("replace SaveAvatar: %v", err)
	}
	got, err = repo.FindByID(u.ID)
	if err != nil {
		t.Fatalf("FindByID after replace: %v", err)
	}
	if !bytes.Equal(got.Avatar.Bytes, []byte{4, 5}) {
		t.Fatalf("got %v, want replaced bytes", got.Avatar.Bytes)
	}

	got.Avatar.Bytes[0] = 7
	again, err := repo.FindByID(u.ID)
	if err != nil {
		t.Fatalf("FindByID after mutate: %v", err)
	}
	if !bytes.Equal(again.Avatar.Bytes, []byte{4, 5}) {
		t.Fatalf("got %v, want stored bytes", again.Avatar.Bytes)
	}
}

func TestUserRepository_AvatarErrors(t *testing.T) {
	t.Parallel()

	repo := NewUserRepository()
	u := user.NewUser("demo", "demo@example.com", "secret")
	if err := repo.Create(u); err != nil {
		t.Fatalf("Create: %v", err)
	}

	got, err := repo.FindByID(u.ID)
	if err != nil {
		t.Fatalf("FindByID without save: %v", err)
	}
	if got.Avatar.Bytes != nil {
		t.Fatalf("Avatar.Bytes = %v, want nil", got.Avatar.Bytes)
	}

	missing := user.UserId("missing")
	if err := repo.SaveAvatar(missing, user.Avatar{Bytes: []byte{1}, MediaType: "image/png"}); err == nil || err.Error() != "user not found" {
		t.Fatalf("SaveAvatar missing user: %v", err)
	}
	_, err = repo.FindByID(string(missing))
	if err == nil || err.Error() != "user not found" {
		t.Fatalf("FindByID missing user: %v", err)
	}
}
