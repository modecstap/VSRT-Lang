package card

import (
	domain "VSRT-Lang/internal/card"
	"VSRT-Lang/internal/session"
	"VSRT-Lang/internal/user"
)

type Service interface {
	GetCards(user.UserId, int64, int) ([]session.Record, error)
	UpdateKnowledge(domain.UpdateKnowledgeCommand) error
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}
