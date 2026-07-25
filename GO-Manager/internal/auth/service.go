package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type Service struct {
	userRepo    UserRepository
	refreshRepo RefreshTokenRepository
	jwtService  *JWTService
}

func NewService(userRepo UserRepository, refreshRepo RefreshTokenRepository, jwtService *JWTService) *Service {
	return &Service{
		userRepo:    userRepo,
		refreshRepo: refreshRepo,
		jwtService:  jwtService,
	}
}

func (s *Service) Register(username, email, password string) (*User, error) {
	if username == "" {
		username = email
	}
	if !ValidateEmail(email) {
		return nil, errInvalidEmail
	}
	if !ValidatePassword(password) {
		return nil, errWeakPassword
	}

	hashed, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &User{Username: username, Email: email, Password: hashed, CreatedAt: time.Now()}
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) Login(login, password string) (*TokenPair, error) {
	user, err := s.userRepo.FindByUsername(login)
	if err != nil {
		return nil, err
	}

	if user == nil {
		user, err = s.userRepo.FindByEmail(login)
		if err != nil {
			return nil, err
		}
	}

	if !ComparePassword(user.Password, password) {
		return nil, fmt.Errorf("invalid password")
	}

	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, user.Username, user.Email)
	if err != nil {
		return nil, err
	}

	refreshTokenValue := s.jwtService.GenerateRefreshToken()
	refreshHash := hashToken(refreshTokenValue)
	refreshToken := &RefreshToken{
		UserID:    user.ID,
		TokenHash: refreshHash,
		ExpiresAt: time.Now().Add(refreshTokenTTL),
		CreatedAt: time.Now(),
	}
	if err := s.refreshRepo.Save(refreshToken); err != nil {
		return nil, err
	}

	return &TokenPair{AccessToken: accessToken, RefreshToken: refreshTokenValue}, nil
}

func (s *Service) Refresh(accessToken, refreshTokenValue string) (*TokenPair, error) {
	claims, err := s.jwtService.ParseAccessToken(accessToken)
	if err != nil {
		return nil, err
	}

	refreshHash := hashToken(refreshTokenValue)
	storedToken, err := s.refreshRepo.FindByHash(refreshHash)
	if err != nil {
		return nil, err
	}
	if storedToken.Revoked || time.Now().After(storedToken.ExpiresAt) {
		return nil, fmt.Errorf("refresh token expired")
	}

	if err := s.refreshRepo.RevokeByHash(refreshHash); err != nil {
		return nil, err
	}

	newAccessToken, err := s.jwtService.GenerateAccessToken(claims.Subject, claims.Username, claims.Email)
	if err != nil {
		return nil, err
	}

	newRefreshTokenValue := s.jwtService.GenerateRefreshToken()
	newRefreshHash := hashToken(newRefreshTokenValue)
	newToken := &RefreshToken{UserID: claims.Subject, TokenHash: newRefreshHash, ExpiresAt: time.Now().Add(refreshTokenTTL), CreatedAt: time.Now()}
	if err := s.refreshRepo.Save(newToken); err != nil {
		return nil, err
	}

	return &TokenPair{AccessToken: newAccessToken, RefreshToken: newRefreshTokenValue}, nil
}

func (s *Service) Logout(refreshTokenValue string) error {
	refreshHash := hashToken(refreshTokenValue)
	return s.refreshRepo.RevokeByHash(refreshHash)
}

func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
