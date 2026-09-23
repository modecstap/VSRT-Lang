package user

import (
	"VSRT-Lang/internal/session"
	domain "VSRT-Lang/internal/user"
)

type SessionRepository interface {
	FindByUser(userId domain.UserId) ([]session.Session, error)
}

type Service interface {
	SaveAvatar(id domain.UserId, raw []byte) error
	Get(id domain.UserId) (*domain.User, error)
}

type Handler struct {
	service     Service
	sessionRepo SessionRepository
}

func NewHandler(
	service Service,
	sessionRepo SessionRepository,
) *Handler {
	return &Handler{
		service:     service,
		sessionRepo: sessionRepo,
	}
}
