package user_test

import (
	"bytes"
	"errors"
	"testing"

	"VSRT-Lang/internal/database/memory"
	"VSRT-Lang/internal/user"
)

func TestService_UpdateProfile(t *testing.T) {
	t.Parallel()

	seed := func(t *testing.T) (*user.Service, *memory.UserRepository, *user.User) {
		t.Helper()
		repo := memory.NewUserRepository()
		owner := user.NewUser("alice", "alice@example.com", "hash-alice")
		if err := repo.Create(owner); err != nil {
			t.Fatalf("Create owner: %v", err)
		}
		avatar := user.Avatar{Bytes: []byte{1, 2, 3}, MediaType: "image/png"}
		if err := repo.SaveAvatar(user.UserId(owner.ID), avatar); err != nil {
			t.Fatalf("SaveAvatar: %v", err)
		}
		other := user.NewUser("bob", "bob@example.com", "hash-bob")
		if err := repo.Create(other); err != nil {
			t.Fatalf("Create other: %v", err)
		}
		return user.NewService(repo), repo, owner
	}

	assertOwner := func(t *testing.T, repo *memory.UserRepository, id, username, email string) {
		t.Helper()
		got, err := repo.FindByID(id)
		if err != nil {
			t.Fatalf("FindByID: %v", err)
		}
		if got.Username != username || got.Email != email || got.Password != "hash-alice" || !bytes.Equal(got.Avatar.Bytes, []byte{1, 2, 3}) {
			t.Fatalf("stored user = %+v", got)
		}
	}

	t.Run("success keeps password and avatar", func(t *testing.T) {
		svc, repo, owner := seed(t)
		if err := svc.UpdateProfile(user.UserId(owner.ID), "alice2", "alice2@example.com"); err != nil {
			t.Fatalf("UpdateProfile: %v", err)
		}
		assertOwner(t, repo, owner.ID, "alice2", "alice2@example.com")
	})

	t.Run("same username and email", func(t *testing.T) {
		svc, repo, owner := seed(t)
		if err := svc.UpdateProfile(user.UserId(owner.ID), "alice", "alice@example.com"); err != nil {
			t.Fatalf("UpdateProfile: %v", err)
		}
		assertOwner(t, repo, owner.ID, "alice", "alice@example.com")
	})

	t.Run("username taken", func(t *testing.T) {
		svc, repo, owner := seed(t)
		err := svc.UpdateProfile(user.UserId(owner.ID), "bob", "alice2@example.com")
		if !errors.Is(err, user.ErrUsernameTaken) {
			t.Fatalf("UpdateProfile error = %v", err)
		}
		assertOwner(t, repo, owner.ID, "alice", "alice@example.com")
	})

	t.Run("email taken", func(t *testing.T) {
		svc, repo, owner := seed(t)
		err := svc.UpdateProfile(user.UserId(owner.ID), "alice2", "bob@example.com")
		if !errors.Is(err, user.ErrEmailTaken) {
			t.Fatalf("UpdateProfile error = %v", err)
		}
		assertOwner(t, repo, owner.ID, "alice", "alice@example.com")
	})

	t.Run("empty username", func(t *testing.T) {
		svc, repo, owner := seed(t)
		err := svc.UpdateProfile(user.UserId(owner.ID), "", "alice2@example.com")
		if !errors.Is(err, user.ErrUsernameRequired) {
			t.Fatalf("UpdateProfile error = %v", err)
		}
		assertOwner(t, repo, owner.ID, "alice", "alice@example.com")
	})
}
