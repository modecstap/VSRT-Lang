package memory

import (
	"VSRT-Lang/internal/auth"
	"errors"
	"sync"
	"time"
)

type UserRepository struct {
	mu              sync.RWMutex
	usersByID       map[string]*auth.User
	usersByEmail    map[string]*auth.User
	usersByUsername map[string]*auth.User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		usersByID:       make(map[string]*auth.User),
		usersByEmail:    make(map[string]*auth.User),
		usersByUsername: make(map[string]*auth.User),
	}
}

func (r *UserRepository) Create(user *auth.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.usersByEmail[user.Email]; exists {
		return errors.New("user already exists")
	}

	if user.ID == "" {
		user.ID = "user-" + time.Now().Format("20060102150405")
	}

	r.usersByID[user.ID] = user
	r.usersByEmail[user.Email] = user
	r.usersByUsername[user.Username] = user
	return nil
}

func (r *UserRepository) FindByEmail(email string) (*auth.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.usersByEmail[email]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r *UserRepository) FindByUsername(username string) (*auth.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.usersByUsername[username]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r *UserRepository) FindByID(id string) (*auth.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.usersByID[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}
