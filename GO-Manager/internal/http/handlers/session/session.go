package session

import (
	domain "VSRT-Lang/internal/session"
)

type Repository interface {
	Save(session *domain.Session) error
	Take(sessionID int) (domain.Session, error)
}

type Handler struct {
	repo       Repository
	translator domain.Translator
}

func NewHandler(repo Repository, translator domain.Translator) *Handler {
	return &Handler{repo: repo, translator: translator}
}
