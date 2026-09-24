package memory

import (
	"VSRT-Lang/internal/user"
	"errors"
	"sync"
	"time"
)

type UserRepository struct {
	mu              sync.RWMutex
	usersByID       map[string]*user.User
	usersByEmail    map[string]*user.User
	usersByUsername map[string]*user.User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		usersByID:       make(map[string]*user.User),
		usersByEmail:    make(map[string]*user.User),
		usersByUsername: make(map[string]*user.User),
	}
}

func (r *UserRepository) Create(u *user.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.usersByEmail[u.Email]; exists {
		return errors.New("user already exists")
	}

	if u.ID == "" {
		u.ID = "user-" + time.Now().Format("20060102150405")
	}

	stored := cloneUser(u)
	r.usersByID[stored.ID] = stored
	r.usersByEmail[stored.Email] = stored
	r.usersByUsername[stored.Username] = stored
	return nil
}

func (r *UserRepository) FindByEmail(email string) (*user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	stored, ok := r.usersByEmail[email]
	if !ok {
		return nil, user.ErrNotFound
	}
	return cloneUser(stored), nil
}

func (r *UserRepository) FindByUsername(username string) (*user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	stored, ok := r.usersByUsername[username]
	if !ok {
		return nil, user.ErrNotFound
	}
	return cloneUser(stored), nil
}

func (r *UserRepository) FindByID(id string) (*user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	stored, ok := r.usersByID[id]
	if !ok {
		return nil, user.ErrNotFound
	}
	return cloneUser(stored), nil
}

func (r *UserRepository) Save(u *user.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	old, ok := r.usersByID[u.ID]
	if !ok {
		return user.ErrNotFound
	}
	if other, exists := r.usersByUsername[u.Username]; exists && other.ID != u.ID {
		return user.ErrUsernameTaken
	} else if other, exists := r.usersByEmail[u.Email]; exists && other.ID != u.ID {
		return user.ErrEmailTaken
	}
	if old.Email != u.Email {
		delete(r.usersByEmail, old.Email)
	}
	if old.Username != u.Username {
		delete(r.usersByUsername, old.Username)
	}
	stored := cloneUser(u)
	r.usersByID[stored.ID] = stored
	r.usersByEmail[stored.Email] = stored
	r.usersByUsername[stored.Username] = stored
	return nil
}

func (r *UserRepository) SaveAvatar(id user.UserId, avatar user.Avatar) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	stored, ok := r.usersByID[string(id)]
	if !ok {
		return user.ErrNotFound
	}
	copied := make([]byte, len(avatar.Bytes))
	copy(copied, avatar.Bytes)
	stored.Avatar = user.Avatar{Bytes: copied, MediaType: avatar.MediaType}
	return nil
}

func cloneUser(src *user.User) *user.User {
	dst := *src
	if src.Avatar.Bytes != nil {
		dst.Avatar.Bytes = make([]byte, len(src.Avatar.Bytes))
		copy(dst.Avatar.Bytes, src.Avatar.Bytes)
	}
	return &dst
}
