package user

import (
	"VSRT-Lang/internal/session"
	domain "VSRT-Lang/internal/user"
)

type SessionRepository interface {
	FindByUser(userId domain.UserId) ([]session.Session, error)
}

type AvatarService interface {
	SaveAvatar(id domain.UserId, raw []byte) error
}

type Handler struct {
	avatarService AvatarService
	sessionRepo   SessionRepository
}

func NewHandler(
	avatarService AvatarService,
	sessionRepo SessionRepository,
) *Handler {
	return &Handler{
		avatarService: avatarService,
		sessionRepo:   sessionRepo,
	}
}
