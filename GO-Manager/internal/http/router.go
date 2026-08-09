package router

import (
	"VSRT-Lang/internal/http/handlers/auth"
	sessionhandler "VSRT-Lang/internal/http/handlers/session"
	"VSRT-Lang/internal/http/handlers/user"
	"VSRT-Lang/internal/http/middleware"
	"net/http"
)

type Handlers struct {
	Auth *auth.Handler
	Session *sessionhandler.Handler
	User *user.Handler
	AuthMiddleware middleware.Middleware
}

func NewServeMux(h Handlers) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /register", h.Auth.Register)
	mux.HandleFunc("POST /login", h.Auth.Login)
	protected := h.AuthMiddleware

	if h.Session != nil {
		if protected == nil {
			protected = func(next http.Handler) http.Handler { return next }
		}

		mux.Handle("POST /sessions", protected(http.HandlerFunc(h.Session.CreateSession)))
		mux.Handle("POST /sessions/", protected(http.HandlerFunc(h.Session.SaveRecord)))
		mux.Handle("GET /sessions/", protected(http.HandlerFunc(h.Session.GetRecords)))
	}

	if h.User != nil {
		if protected == nil {
			protected = func(next http.Handler) http.Handler { return next }
		}

		mux.Handle("GET /users/", protected(http.HandlerFunc(h.User.GetUserSessions)))
	}

	RegisterSwagger(mux)

	return mux
}
