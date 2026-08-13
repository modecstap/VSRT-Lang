package main

import (
	"VSRT-Lang/internal/database/postgres"
	"VSRT-Lang/internal/database/postgres/migrations"
	myHttp "VSRT-Lang/internal/http"
	router "VSRT-Lang/internal/http"
	auth_handler "VSRT-Lang/internal/http/handlers/auth"
	session_handler "VSRT-Lang/internal/http/handlers/session"
	user_handler "VSRT-Lang/internal/http/handlers/user"
	"VSRT-Lang/internal/http/middleware"
	"database/sql"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {

	db := setupDb()

	deps := router.DependensFromEnv(db)

	startServer(deps)
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

func startServer(deps *router.ServerDependens) {
	handlers := myHttp.Handlers{
		Auth:           auth_handler.NewAuth(deps.AuthService),
		Session:        session_handler.NewHandler(deps.SessionRepo, deps.Translator),
		User:           user_handler.NewHandler(deps.UserRepo, deps.SessionRepo),
		AuthMiddleware: middleware.Auth(deps.JwtService),
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

	http.ListenAndServe(deps.Host, cors(mux))
}
