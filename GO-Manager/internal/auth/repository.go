package auth

import (
	"time"
)

type User struct {
	ID        string
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
}

type RefreshToken struct {
	ID        string
	UserID    string
	TokenHash string
	ExpiresAt time.Time
	CreatedAt time.Time
	Revoked   bool
}

type UserRepository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
	FindByUsername(email string) (*User, error)
	FindByID(id string) (*User, error)
}

type RefreshTokenRepository interface {
	Save(token *RefreshToken) error
	FindByHash(hash string) (*RefreshToken, error)
	RevokeByHash(hash string) error
	DeleteByHash(hash string) error
}
