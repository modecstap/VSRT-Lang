package middleware

import (
	"VSRT-Lang/internal/auth"
	"VSRT-Lang/internal/user"
	"context"
	"net/http"
	"strings"
)

type Middleware func(http.Handler) http.Handler

type contextKey string

const (
	userIDContextKey contextKey = "user_id"
	emailContextKey  contextKey = "email"
)

func Auth(jwt *auth.JWTService) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				http.Error(w, "invalid authorization header", http.StatusUnauthorized)
				return
			}

			claims, err := jwt.ParseAccessToken(parts[1])
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDContextKey, claims.Subject)
			ctx = context.WithValue(ctx, emailContextKey, claims.Email)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserFromContext(ctx context.Context) (user.UserId, string, bool) {
	userID, ok := ctx.Value(userIDContextKey).(string)
	if !ok || userID == "" {
		return "", "", false
	}

	email, _ := ctx.Value(emailContextKey).(string)
	return user.UserId(userID), email, true
}
