package auth

import "testing"

func TestJWTService_GenerateAndParseAccessToken(t *testing.T) {
	svc := NewJWTService("test-secret")

	token, err := svc.GenerateAccessToken("user-123", "johndoe", "user@example.com")
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	claims, err := svc.ParseAccessToken(token)
	if err != nil {
		t.Fatalf("parse token: %v", err)
	}

	if claims.Subject != "user-123" {
		t.Fatalf("expected subject user-123, got %s", claims.Subject)
	}

	if claims.Email != "user@example.com" {
		t.Fatalf("expected email user@example.com, got %s", claims.Email)
	}

	if claims.Username != "johndoe" {
		t.Fatalf("expected username johndoe, got %s", claims.Username)
	}
}
