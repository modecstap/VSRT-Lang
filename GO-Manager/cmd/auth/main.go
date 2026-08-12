package main

import (
	"VSRT-Lang/internal"
	"VSRT-Lang/internal/auth"
	"VSRT-Lang/internal/database/postgres"
	"VSRT-Lang/internal/database/postgres/migrations"
	"VSRT-Lang/internal/database/postgres/refresh_token_repository"
	"VSRT-Lang/internal/database/postgres/session_repository"
	"VSRT-Lang/internal/database/postgres/user_repository"
	myHttp "VSRT-Lang/internal/http"
	auth_handler "VSRT-Lang/internal/http/handlers/auth"
	session_handler "VSRT-Lang/internal/http/handlers/session"
	user_handler "VSRT-Lang/internal/http/handlers/user"
	"VSRT-Lang/internal/http/middleware"
	"database/sql"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	host := "localhost:8080"

	db := setupDb()

	cors, mux := setupServer(db)

	http.ListenAndServe(host, cors(mux))
}

func setupDb() *sql.DB {
	dbConfig, err := postgres.LoadDBConfig()
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("postgres", dbConfig.ConnString())
	if err != nil {
		panic(err)
	}

	err = migrations.Run(db)
	if err != nil {
		panic(err)
	}
	return db
}

func setupServer(db *sql.DB) (func(http.Handler) http.Handler, *http.ServeMux) {
	userRepo := user_repository.New(db)
	tokenRepo := refresh_token_repository.New(db)
	jwtService := auth.NewJWTService("StrongSecretString")
	service := auth.NewService(
		userRepo,
		tokenRepo,
		jwtService,
	)
	authHandler := auth_handler.NewAuth(service)

	sessionRepo := session_repository.New(db)
	sessionTranslator := internal.MockTranslator{}
	sessionHandler := session_handler.NewHandler(sessionRepo, sessionTranslator)

	userHandler := user_handler.NewHandler(userRepo, sessionRepo)

	handlers := myHttp.Handlers{
		Auth:           authHandler,
		Session:        sessionHandler,
		User:           userHandler,
		AuthMiddleware: middleware.Auth(jwtService),
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
	return cors, mux
}
