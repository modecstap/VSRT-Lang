package memory

import (
	"sync"

	"VSRT-Lang/internal/passwordreset"
	"VSRT-Lang/internal/user"
)

type PasswordResetRepository struct {
	mu     sync.RWMutex
	byUser map[user.UserId]*passwordreset.ResetLink
	byHash map[string]*passwordreset.ResetLink
}

func NewPasswordResetRepository() *PasswordResetRepository {
	return &PasswordResetRepository{
		byUser: make(map[user.UserId]*passwordreset.ResetLink),
		byHash: make(map[string]*passwordreset.ResetLink),
	}
}

func (r *PasswordResetRepository) Save(link *passwordreset.ResetLink) (user.UserId, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if old, ok := r.byUser[link.User]; ok && old.TokenHash != "" {
		delete(r.byHash, old.TokenHash)
	}

	stored := cloneResetLink(link)
	r.byUser[stored.User] = stored
	if stored.TokenHash != "" {
		r.byHash[stored.TokenHash] = stored
	}
	return stored.User, nil
}

func (r *PasswordResetRepository) Take(id user.UserId) (*passwordreset.ResetLink, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	stored, ok := r.byUser[id]
	if !ok {
		return nil, nil
	}
	return cloneResetLink(stored), nil
}

func (r *PasswordResetRepository) TakeForUpdate(id user.UserId) (*passwordreset.ResetLink, error) {
	return r.Take(id)
}

func (r *PasswordResetRepository) TakeByHash(hash string) (*passwordreset.ResetLink, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if hash == "" {
		return nil, nil
	}
	stored, ok := r.byHash[hash]
	if !ok {
		return nil, nil
	}
	return cloneResetLink(stored), nil
}

func cloneResetLink(src *passwordreset.ResetLink) *passwordreset.ResetLink {
	dst := *src
	if src.ExpiresAt != nil {
		t := *src.ExpiresAt
		dst.ExpiresAt = &t
	}
	return &dst
}

type PasswordResetTransactor struct {
	Users   *UserRepository
	Links   *PasswordResetRepository
	Refresh *RefreshTokenRepository
}

func NewPasswordResetTransactor(
	users *UserRepository,
	links *PasswordResetRepository,
	refresh *RefreshTokenRepository,
) *PasswordResetTransactor {
	return &PasswordResetTransactor{Users: users, Links: links, Refresh: refresh}
}

func (t *PasswordResetTransactor) Within(fn func(passwordreset.TxRepos) error) error {
	return fn(passwordreset.TxRepos{
		Users:   t.Users,
		Links:   t.Links,
		Refresh: t.Refresh,
	})
}
