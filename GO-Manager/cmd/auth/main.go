package main

import (
	"VSRT-Lang/internal/auth"
	"VSRT-Lang/internal/database/memory"
	myHttp "VSRT-Lang/internal/http"
	auth_handler "VSRT-Lang/internal/http/handlers/auth"
	"VSRT-Lang/internal/http/middleware"
	"net/http"
)

func main() {
	host := "localhost:8080"

	service := auth.NewService(
		memory.NewUserRepository(),
		memory.NewRefreshTokenRepository(),
		auth.NewJWTService("StrongSecretString"),
	)
	authHandler := auth_handler.NewAuth(service)

	handlers := myHttp.Handlers{
		Auth: authHandler,
	}

	cors := middleware.CORS(middleware.CORSConfig{
		AllowedOrigins: []string{
			"*",
		},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"Content-Type",
			"Authorization",
		},
		AllowCredentials: true,
	})

	mux := myHttp.NewServeMux(handlers)

	http.ListenAndServe(host, cors(mux))
}
