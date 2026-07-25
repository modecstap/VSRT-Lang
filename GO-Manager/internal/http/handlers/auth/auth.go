package auth

import (
	"VSRT-Lang/internal/auth"
)

type Handler struct {
	service *auth.Service
}

func NewAuth(service *auth.Service) *Handler {
	return &Handler{service: service}
}
