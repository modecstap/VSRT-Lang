package auth_test

import (
	"testing"

	"VSRT-Lang/internal/auth"
	"VSRT-Lang/internal/database/memory"
)

func TestRevokeByUserID_BlocksRefreshKeepsAccess(t *testing.T) {
	users := memory.NewUserRepository()
	tokens := memory.NewRefreshTokenRepository()
	jwtSvc := auth.NewJWTService("test-secret")
	svc := auth.NewService(users, tokens, jwtSvc)

	u, err := svc.Register("demo", "demo@example.com", "Password1")
	if err != nil {
		t.Fatalf("Register: %v", err)
	}

	pair, err := svc.Login("demo", "Password1")
	if err != nil {
		t.Fatalf("Login: %v", err)
	}

	if err := tokens.RevokeByUserID(u.ID); err != nil {
		t.Fatalf("RevokeByUserID: %v", err)
	}

	if _, err := svc.Refresh(pair.AccessToken, pair.RefreshToken); err == nil {
		t.Fatal("Refresh after revoke must fail")
	}

	if _, err := jwtSvc.ParseAccessToken(pair.AccessToken); err != nil {
		t.Fatalf("ParseAccessToken after revoke: %v", err)
	}
}
