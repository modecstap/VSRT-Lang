package auth

import (
	"time"
)

type RefreshToken struct {
	ID        string
	UserID    string
	TokenHash string
	ExpiresAt time.Time
	CreatedAt time.Time
	Revoked   bool
}

type RefreshTokenRepository interface {
	Save(token *RefreshToken) error
	FindByHash(hash string) (*RefreshToken, error)
	RevokeByHash(hash string) error
	DeleteByHash(hash string) error
}
