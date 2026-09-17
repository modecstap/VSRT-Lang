package router

import (
	_ "VSRT-Lang/docs"
	"VSRT-Lang/internal/http/handlers/auth"
	sessionhandler "VSRT-Lang/internal/http/handlers/session"
	"VSRT-Lang/internal/http/handlers/user"
	"VSRT-Lang/internal/http/middleware"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

type Handlers struct {
	Auth           *auth.Handler
	Session        *sessionhandler.Handler
	User           *user.Handler
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
		mux.Handle("DELETE /sessions/", protected(http.HandlerFunc(h.Session.DeleteSession)))
		mux.Handle("POST /sessions/", protected(http.HandlerFunc(h.Session.SaveRecord)))
		mux.Handle("GET /sessions/", protected(http.HandlerFunc(h.Session.GetRecords)))
	}

	if h.User != nil {
		if protected == nil {
			protected = func(next http.Handler) http.Handler { return next }
		}

		mux.Handle("GET /users/", protected(http.HandlerFunc(h.User.GetUserSessions)))
	}

	mux.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	return mux
}
