package session

import (
	domain "VSRT-Lang/internal/session"
	"VSRT-Lang/internal/user"
)

type Repository interface {
	Save(session *domain.Session) error
	Take(sessionID int) (domain.Session, error)
	FindByUser(userID user.UserId) ([]domain.Session, error)
	Delete(sessionID int) error
}

type Handler struct {
	repo       Repository
	translator domain.Translator
}

func NewHandler(repo Repository, translator domain.Translator) *Handler {
	return &Handler{repo: repo, translator: translator}
}
