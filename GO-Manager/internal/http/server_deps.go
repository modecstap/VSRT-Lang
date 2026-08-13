package router

import (
	"VSRT-Lang/internal"
	"VSRT-Lang/internal/auth"
	"VSRT-Lang/internal/database/postgres/refresh_token_repository"
	"VSRT-Lang/internal/database/postgres/session_repository"
	"VSRT-Lang/internal/database/postgres/user_repository"
	"VSRT-Lang/internal/session"
	"database/sql"
	"os"
)

type ServerDependens struct {
	UserRepo    *user_repository.Repository
	SessionRepo *session_repository.Repository
	AuthService *auth.Service
	JwtService  *auth.JWTService
	Translator  session.Translator
	Host        string
}

func DependensFromEnv(db *sql.DB) *ServerDependens {
	host, ok := os.LookupEnv("MANAGER_HOST")
	if !ok {
		panic("environment variable MANAGER_HOST not found")
	}

	secret, ok := os.LookupEnv("SECRET_KEY")
	if !ok {
		panic("environment variable SECRET_KEY not found")
	}

	userRepo := user_repository.New(db)
	tokenRepo := refresh_token_repository.New(db)
	sessionRepo := session_repository.New(db)

	translator := internal.MockTranslator{}
	jwtService := auth.NewJWTService(secret)
	authService := auth.NewService(
		userRepo,
		tokenRepo,
		jwtService,
	)

	return &ServerDependens{
		UserRepo:    userRepo,
		SessionRepo: sessionRepo,
		AuthService: authService,
		JwtService:  jwtService,
		Translator:  translator,
		Host:        host,
	}
}
