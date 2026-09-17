package session

import (
	domain "VSRT-Lang/internal/session"
	"VSRT-Lang/internal/user"
)

type Service interface {
	NewSession(user.UserId, string) (int64, error)
	GetSession(user.UserId, int64) (domain.Session, error)
	DeleteSession(user.UserId, int64) error
	AddRecord(domain.AddRecordCommand) (domain.Record, error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}
