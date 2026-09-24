package router

import (
	"VSRT-Lang/internal/auth"
	"VSRT-Lang/internal/database/postgres/password_reset_repository"
	"VSRT-Lang/internal/database/postgres/refresh_token_repository"
	"VSRT-Lang/internal/database/postgres/session_repository"
	"VSRT-Lang/internal/database/postgres/user_repository"
	"VSRT-Lang/internal/mail"
	"VSRT-Lang/internal/net_translator"
	"VSRT-Lang/internal/passwordreset"
	"VSRT-Lang/internal/session"
	"VSRT-Lang/internal/translators/stub"
	"VSRT-Lang/internal/user"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
)

type ServerDependens struct {
	UserRepo             *user_repository.Repository
	UserService          *user.Service
	SessionRepo          *session_repository.Repository
	SessionService       *session.Service
	AuthService          *auth.Service
	JwtService           *auth.JWTService
	PasswordResetService *passwordreset.Service
	Host                 string
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

	sessionRepo := session_repository.New(db)
	translator, err := setupTranslator(mode)
	if err != nil {
		return nil, err
	}
	sessionService := session.NewService(sessionRepo, translator)

	userRepo := user_repository.New(db)
	userService := user.NewService(userRepo)
	tokenRepo := refresh_token_repository.New(db)
	secret, ok := os.LookupEnv("SECRET_KEY")
	if !ok {
		return nil, errors.New("environment variable SECRET_KEY not found")
	}
	jwtService := auth.NewJWTService(secret)
	authService := auth.NewService(
		userRepo,
		tokenRepo,
		jwtService,
	)

	publicSiteURL := os.Getenv("PUBLIC_SITE_URL")
	mailCfg := mail.LoadConfig()
	mailer := mail.NewSender(mailCfg)
	linkRepo := password_reset_repository.New(db)
	tx := password_reset_repository.NewTransactor(db)
	passwordResetService := passwordreset.NewService(
		userRepo,
		linkRepo,
		mailer,
		tx,
		publicSiteURL,
		nil,
	)

	return &ServerDependens{
		UserRepo:             userRepo,
		UserService:          userService,
		SessionRepo:          sessionRepo,
		SessionService:       sessionService,
		AuthService:          authService,
		JwtService:           jwtService,
		PasswordResetService: passwordResetService,
		Host:                 host,
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
	return *net_translator.NewTranslator(translatorHost), nil
}
