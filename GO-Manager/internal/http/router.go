package router

import (
	"VSRT-Lang/internal/http/handlers/auth"
	"net/http"
)

type Handlers struct {
	Auth *auth.Handler
}

func NewServeMux(h Handlers) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /register", h.Auth.Register)
	mux.HandleFunc("POST /login", h.Auth.Login)
	// register embedded swagger handlers
	RegisterSwagger(mux)

	return mux
}
