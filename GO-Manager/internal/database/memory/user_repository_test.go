package memory

import (
	"bytes"
	"testing"

	"VSRT-Lang/internal/user"
)

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
