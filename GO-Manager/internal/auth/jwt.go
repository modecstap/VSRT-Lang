package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	accessTokenTTL  = 15 * time.Minute
	refreshTokenTTL = 30 * 24 * time.Hour
)

type Claims struct {
	Subject  string `json:"sub"`
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type JWTService struct {
	secret string
}

func NewJWTService(secret string) *JWTService {
	return &JWTService{secret: secret}
}

func (s *JWTService) GenerateAccessToken(userID, username, email string) (string, error) {
	now := time.Now()
	claims := Claims{
		Subject:  userID,
		Username: username,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(accessTokenTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

func (s *JWTService) ParseAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString, &Claims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf(
					"unexpected signing method: %v", 
					token.Header["alg"],
				)
			}
			return []byte(s.secret), nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func (s *JWTService) GenerateRefreshToken() string {
	return fmt.Sprintf("refresh-%d", time.Now().UnixNano())
}
