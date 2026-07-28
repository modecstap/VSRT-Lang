package main

import (
	"VSRT-Lang/internal/auth"
	"VSRT-Lang/internal/database/postgres/migrations"
	"VSRT-Lang/internal/database/postgres/refresh_token_repository"
	"VSRT-Lang/internal/database/postgres/user_repository"
	myHttp "VSRT-Lang/internal/http"
	auth_handler "VSRT-Lang/internal/http/handlers/auth"
	"VSRT-Lang/internal/http/middleware"
	"database/sql"
	"net/http"
	_ "github.com/lib/pq"
)

func main() {
	host := "localhost:8080"

	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=vsrt_lang sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	if err := migrations.Run(db); err != nil {
		panic(err)
	}

	service := auth.NewService(
		user_repository.New(db),
		refresh_token_repository.New(db),
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
