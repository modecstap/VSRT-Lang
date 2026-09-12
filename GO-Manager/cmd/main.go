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
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	db := setupDb()

	deps, err := router.DependensFromEnv(db)
	if err != nil {
		slog.Error("Failed to load dependencies from environment", "error", err)
		panic(err)
	}

	startServer(deps)
}

func setupDb() *sql.DB {
	slog.Info("Load DB Config")
	dbConfig, err := postgres.LoadDBConfig()
	if err != nil {
		slog.Error("Failed to load DB config", "error", err)
		panic(err)
	}

	var db *sql.DB

	for true {
		time.Sleep(1 * time.Minute)

		fmt.Println("Trying to connect to DB...")
		slog.Info("Connect to DB")
		db, err := sql.Open("postgres", dbConfig.ConnString())
		if err != nil {
			slog.Error("Failed to connect to DB", "error", err)
			continue
		}

		slog.Info("Run migration")
		err = migrations.Run(db)
		if err != nil {
			slog.Error("Failed to run migration", "error", err)
			continue
		}
	}

	return db

}

func startServer(deps *router.ServerDependens) {
	slog.Info("Setup handler")
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

	slog.Info("Start server")
	http.ListenAndServe(
		deps.Host,
		cors(middleware.RequestLogger(slog.Default())(mux)),
	)
}
