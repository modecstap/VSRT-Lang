package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"VSRT-Lang/internal/user"
)

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type Service struct {
	userRepo    user.Repository
	refreshRepo RefreshTokenRepository
	jwtService  *JWTService
}

func NewService(userRepo user.Repository, refreshRepo RefreshTokenRepository, jwtService *JWTService) *Service {
	return &Service{
		userRepo:    userRepo,
		refreshRepo: refreshRepo,
		jwtService:  jwtService,
	}
}

func (s *Service) Register(username, email, password string) (*user.User, error) {
	if username == "" {
		username = email
	}
	if !user.ValidateEmail(email) {
		return nil, user.ErrInvalidEmail
	}
	if !user.ValidatePassword(password) {
		return nil, user.ErrWeakPassword
	}

	hashed, err := user.HashPassword(password)
	if err != nil {
		return nil, err
	}

	u := user.NewUser(username, email, hashed)
	err = s.userRepo.Create(u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *Service) Login(login, password string) (*TokenPair, error) {
	foundUser, err := s.findUser(login)
	if err != nil {
		return nil, err
	}

	if !user.ComparePassword(foundUser.Password, password) {
		return nil, fmt.Errorf("invalid password")
	}

	accessToken, err := s.jwtService.GenerateAccessToken(
		foundUser.ID,
		foundUser.Username,
		foundUser.Email,
	)
	if err != nil {
		return nil, err
	}

	refreshTokenValue := s.jwtService.GenerateRefreshToken()
	refreshHash := hashToken(refreshTokenValue)
	refreshToken := &RefreshToken{
		UserID:    foundUser.ID,
		TokenHash: refreshHash,
		ExpiresAt: time.Now().Add(refreshTokenTTL),
		CreatedAt: time.Now(),
	}

	err = s.refreshRepo.Save(refreshToken)
	if err != nil {
		return nil, err
	}

	return &TokenPair{AccessToken: accessToken, RefreshToken: refreshTokenValue}, nil
}

func (s *Service) findUser(login string) (*user.User, error) {
	foundUser, err := s.userRepo.FindByUsername(login)
	if err != nil {
		return nil, err
	}

	if foundUser != nil {
		return foundUser, nil
	}

	return s.userRepo.FindByEmail(login)
}

func (s *Service) Refresh(accessToken, refreshTokenValue string) (*TokenPair, error) {
	claims, err := s.jwtService.ParseAccessToken(accessToken)
	if err != nil {
		return nil, err
	}

	newAccessToken, err := s.updateAccessToken(refreshTokenValue, claims)
	if err != nil {
		return nil, err
	}

	newRefreshTokenValue, err := s.updateRefreshToken(claims)
	if err != nil {
		return nil, err
	}

	tokenPair := &TokenPair{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshTokenValue,
	}
	return tokenPair, nil
}

func (s *Service) updateRefreshToken(claims *Claims) (string, error) {
	newRefreshTokenValue := s.jwtService.GenerateRefreshToken()
	newRefreshHash := hashToken(newRefreshTokenValue)

	newToken := &RefreshToken{
		UserID:    claims.Subject,
		TokenHash: newRefreshHash,
		ExpiresAt: time.Now().Add(refreshTokenTTL),
		CreatedAt: time.Now(),
	}

	err := s.refreshRepo.Save(newToken)
	if err != nil {
		return "", err
	}
	return newRefreshTokenValue, nil
}

func (s *Service) updateAccessToken(refreshTokenValue string, claims *Claims) (string, error) {
	refreshHash := hashToken(refreshTokenValue)
	storedToken, err := s.refreshRepo.FindByHash(refreshHash)
	if err != nil {
		return "", err
	}

	if storedToken.Revoked || time.Now().After(storedToken.ExpiresAt) {
		return "", fmt.Errorf("refresh token expired")
	}

	err = s.refreshRepo.RevokeByHash(refreshHash)
	if err != nil {
		return "", err
	}

	newAccessToken, err := s.jwtService.GenerateAccessToken(claims.Subject, claims.Username, claims.Email)
	if err != nil {
		return "", err
	}
	return newAccessToken, nil
}

func (s *Service) Logout(refreshTokenValue string) error {
	refreshHash := hashToken(refreshTokenValue)
	return s.refreshRepo.RevokeByHash(refreshHash)
}

func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
