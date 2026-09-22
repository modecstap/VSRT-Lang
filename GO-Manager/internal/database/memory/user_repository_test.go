package memory

import (
	"bytes"
	"errors"
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

	got, err := repo.GetAvatar(user.UserId(u.ID))
	if err != nil {
		t.Fatalf("GetAvatar: %v", err)
	}
	if got.MediaType != "image/png" || !bytes.Equal(got.Bytes, []byte{1, 2, 3}) {
		t.Fatalf("got %+v, want copied first avatar", got)
	}

	second := user.Avatar{Bytes: []byte{4, 5}, MediaType: "image/png"}
	if err := repo.SaveAvatar(user.UserId(u.ID), second); err != nil {
		t.Fatalf("replace SaveAvatar: %v", err)
	}
	got, err = repo.GetAvatar(user.UserId(u.ID))
	if err != nil {
		t.Fatalf("GetAvatar after replace: %v", err)
	}
	if !bytes.Equal(got.Bytes, []byte{4, 5}) {
		t.Fatalf("got %v, want replaced bytes", got.Bytes)
	}
}

func TestUserRepository_AvatarErrors(t *testing.T) {
	t.Parallel()

	repo := NewUserRepository()
	u := user.NewUser("demo", "demo@example.com", "secret")
	if err := repo.Create(u); err != nil {
		t.Fatalf("Create: %v", err)
	}

	_, err := repo.GetAvatar(user.UserId(u.ID))
	if !errors.Is(err, user.ErrNoAvatar) {
		t.Fatalf("GetAvatar without save: %v, want ErrNoAvatar", err)
	}

	missing := user.UserId("missing")
	if err := repo.SaveAvatar(missing, user.Avatar{Bytes: []byte{1}, MediaType: "image/png"}); err == nil || err.Error() != "user not found" {
		t.Fatalf("SaveAvatar missing user: %v", err)
	}
	_, err = repo.GetAvatar(missing)
	if err == nil || err.Error() != "user not found" {
		t.Fatalf("GetAvatar missing user: %v", err)
	}
}
