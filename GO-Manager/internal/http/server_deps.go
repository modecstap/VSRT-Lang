package router

import (
	"VSRT-Lang/internal/auth"
	"VSRT-Lang/internal/database/postgres/refresh_token_repository"
	"VSRT-Lang/internal/database/postgres/session_repository"
	"VSRT-Lang/internal/database/postgres/user_repository"
	"VSRT-Lang/internal/net_translator"
	"VSRT-Lang/internal/session"
	"VSRT-Lang/internal/translators/stub"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type ServerDependens struct {
	UserRepo    *user_repository.Repository
	SessionRepo *session_repository.Repository
	AuthService *auth.Service
	JwtService  *auth.JWTService
	Translator  session.Translator
	Host        string
}

func DependensFromEnv(db *sql.DB) (*ServerDependens, error) {
	mode, ok := os.LookupEnv("MODE")
	if !ok {
		slog.Info("environment variable MODE not found")
		mode = "prod"
	}
	slog.Info(fmt.Sprintf("Started in %s mode", mode))

	host, ok := os.LookupEnv("MANAGER_HOST")
	if !ok {
		return nil, errors.New("environment variable MANAGER_HOST not found")
	}

	secret, ok := os.LookupEnv("SECRET_KEY")
	if !ok {
		return nil, errors.New("environment variable SECRET_KEY not found")
	}

	userRepo := user_repository.New(db)
	tokenRepo := refresh_token_repository.New(db)
	sessionRepo := session_repository.New(db)

	translator, err := setupTranslator(mode)
	if err != nil {
		return nil, err
	}

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
	}, nil
}

func setupTranslator(mode string) (session.Translator, error) {
	if mode == "debug" {
		return stub.Translator{}, nil
	}

	translatorHost, ok := os.LookupEnv("TRANSLATOR_HOST")
	if !ok {
		return net_translator.Translator{}, errors.New("environment variable TRANSLATOR_HOST not found")
	}
	translator := net_translator.Translator{
		Client: http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       5 * time.Second,
		},
		Backend: fmt.Sprintf("http://%s/api/translator", translatorHost),
	}
	return translator, nil
}
