package passwordreset

import (
	"VSRT-Lang/internal/passwordreset"
)

type Handler struct {
	service *passwordreset.Service
}

func NewHandler(service *passwordreset.Service) *Handler {
	return &Handler{service: service}
}
