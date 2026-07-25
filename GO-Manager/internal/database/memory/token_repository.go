package memory

import (
	"VSRT-Lang/internal/auth"
	"errors"
	"sync"
)

type RefreshTokenRepository struct {
	mu           sync.RWMutex
	tokensByHash map[string]*auth.RefreshToken
}

func NewRefreshTokenRepository() *RefreshTokenRepository {
	return &RefreshTokenRepository{tokensByHash: make(map[string]*auth.RefreshToken)}
}

func (r *RefreshTokenRepository) Save(token *auth.RefreshToken) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tokensByHash[token.TokenHash] = token
	return nil
}

func (r *RefreshTokenRepository) FindByHash(hash string) (*auth.RefreshToken, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	token, ok := r.tokensByHash[hash]
	if !ok {
		return nil, errors.New("refresh token not found")
	}
	return token, nil
}

func (r *RefreshTokenRepository) RevokeByHash(hash string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	token, ok := r.tokensByHash[hash]
	if !ok {
		return errors.New("refresh token not found")
	}
	token.Revoked = true
	return nil
}

func (r *RefreshTokenRepository) DeleteByHash(hash string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.tokensByHash, hash)
	return nil
}
