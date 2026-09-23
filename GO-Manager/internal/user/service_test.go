package user

import (
	"bytes"
	"errors"
	"image"
	"image/gif"
	"image/png"
	"testing"
)

type avatarRepo struct {
	byID map[string]*User
}

func newAvatarRepo(users ...*User) *avatarRepo {
	r := &avatarRepo{byID: map[string]*User{}}
	for _, u := range users {
		r.byID[u.ID] = u
	}
	return r
}

func (r *avatarRepo) Create(u *User) error {
	r.byID[u.ID] = u
	return nil
}

func (r *avatarRepo) FindByEmail(string) (*User, error) {
	return nil, errors.New("user not found")
}

func (r *avatarRepo) FindByUsername(string) (*User, error) {
	return nil, errors.New("user not found")
}

func (r *avatarRepo) FindByID(id string) (*User, error) {
	u, ok := r.byID[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return u, nil
}

func (r *avatarRepo) SaveAvatar(id UserId, avatar Avatar) error {
	u, ok := r.byID[string(id)]
	if !ok {
		return errors.New("user not found")
	}
	copied := append([]byte(nil), avatar.Bytes...)
	u.Avatar = Avatar{Bytes: copied, MediaType: avatar.MediaType}
	return nil
}

func tinyPNG(t *testing.T, w, h int) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := png.Encode(&buf, image.NewRGBA(image.Rect(0, 0, w, h))); err != nil {
		t.Fatalf("png.Encode: %v", err)
	}
	return buf.Bytes()
}

func tinyGIF(t *testing.T) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := gif.Encode(&buf, image.NewRGBA(image.Rect(0, 0, 2, 2)), nil); err != nil {
		t.Fatalf("gif.Encode: %v", err)
	}
	return buf.Bytes()
}

func TestService_Get(t *testing.T) {
	owner := NewUser("owner", "owner@example.com", "Password1")
	if err := owner.SetAvatar(tinyPNG(t, 4, 4)); err != nil {
		t.Fatalf("SetAvatar: %v", err)
	}
	got, err := NewService(newAvatarRepo(owner)).Get(UserId(owner.ID))
	if err != nil || got.Username != owner.Username || got.Email != owner.Email || !bytes.Equal(got.Avatar.Bytes, owner.Avatar.Bytes) {
		t.Fatalf("Get = %+v, %v", got, err)
	}
	if _, err = NewService(newAvatarRepo()).Get("missing"); err == nil || err.Error() != "user not found" {
		t.Fatalf("error = %v, want user not found", err)
	}
}

func TestService_SaveAvatar_StoresPNGForThatUserOnly(t *testing.T) {
	owner := NewUser("owner", "owner@example.com", "Password1")
	other := NewUser("other", "other@example.com", "Password1")
	repo := newAvatarRepo(owner, other)

	raw := tinyPNG(t, 4, 4)
	if err := NewService(repo).SaveAvatar(UserId(owner.ID), raw); err != nil {
		t.Fatalf("SaveAvatar: %v", err)
	}

	want := &User{}
	if err := want.SetAvatar(raw); err != nil {
		t.Fatalf("SetAvatar: %v", err)
	}
	got, err := repo.FindByID(owner.ID)
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}
	if got.Avatar.MediaType != "image/png" || !bytes.Equal(got.Avatar.Bytes, want.Avatar.Bytes) {
		t.Fatal("stored avatar mismatch")
	}
	untouched, err := repo.FindByID(other.ID)
	if err != nil {
		t.Fatalf("FindByID other: %v", err)
	}
	if untouched.Avatar.Bytes != nil {
		t.Fatal("other user avatar changed")
	}
}

func TestService_SaveAvatar_ReplacesBytes(t *testing.T) {
	owner := NewUser("owner", "owner@example.com", "Password1")
	repo := newAvatarRepo(owner)
	svc := NewService(repo)
	if err := svc.SaveAvatar(UserId(owner.ID), tinyPNG(t, 4, 4)); err != nil {
		t.Fatalf("first SaveAvatar: %v", err)
	}
	second := tinyPNG(t, 8, 8)
	if err := svc.SaveAvatar(UserId(owner.ID), second); err != nil {
		t.Fatalf("second SaveAvatar: %v", err)
	}
	want := &User{}
	if err := want.SetAvatar(second); err != nil {
		t.Fatalf("SetAvatar: %v", err)
	}
	got, err := repo.FindByID(owner.ID)
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}
	if !bytes.Equal(got.Avatar.Bytes, want.Avatar.Bytes) {
		t.Fatal("second save did not replace bytes")
	}
}

func TestService_SaveAvatar_UnknownID(t *testing.T) {
	err := NewService(newAvatarRepo()).SaveAvatar("missing", []byte{1})
	if err == nil || err.Error() != "user not found" {
		t.Fatalf("error = %v, want user not found", err)
	}
}

func TestService_SaveAvatar_GIFLeavesZeroAvatar(t *testing.T) {
	owner := NewUser("owner", "owner@example.com", "Password1")
	repo := newAvatarRepo(owner)
	err := NewService(repo).SaveAvatar(UserId(owner.ID), tinyGIF(t))
	if err != ErrInvalidAvatarType {
		t.Fatalf("error = %v, want ErrInvalidAvatarType", err)
	}
	got, findErr := repo.FindByID(owner.ID)
	if findErr != nil {
		t.Fatalf("FindByID: %v", findErr)
	}
	if got.Avatar.Bytes != nil {
		t.Fatal("Avatar.Bytes = non-nil, want nil")
	}
}
