package user

import (
	"VSRT-Lang/internal/session"
	domain "VSRT-Lang/internal/user"
)

type SessionRepository interface {
	FindByUser(userId domain.UserId) ([]session.Session, error)
}

type Handler struct {
	userRepo    domain.Repository
	sessionRepo SessionRepository
}

func NewHandler(
	userRepo domain.Repository, 
	sessionRepo SessionRepository,
	) *Handler {
	return &Handler{
		userRepo: userRepo,
		sessionRepo: sessionRepo ,
	}
}
