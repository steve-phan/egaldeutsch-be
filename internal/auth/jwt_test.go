package auth

import (
	"testing"
	"time"

	"egaldeutsch-be/internal/config"
)

func TestCreateAndParseToken(t *testing.T) {
	cfg := config.JwtConfig{SecretKey: "s3cr3t", Issuer: "egaldeutsch", ExpirationHours: 1}
	userID := "user-123"
	role := "learner"

	token, err := CreateAccessToken(userID, role, cfg)
	if err != nil {
		t.Fatalf("CreateAccessToken failed: %v", err)
	}

	claims, err := ParseToken(token, cfg)
	if err != nil {
		t.Fatalf("ParseToken failed: %v", err)
	}
	if claims.UserId != userID {
		t.Fatalf("expected user id %s got %s", userID, claims.UserId)
	}
	if claims.Role != role {
		t.Fatalf("expected role %s got %s", role, claims.Role)
	}
}

func TestTokenExpiry(t *testing.T) {
	cfg := config.JwtConfig{SecretKey: "s3cr3t", Issuer: "egaldeutsch", ExpirationHours: 0}
	userID := "user-123"
	role := "learner"

	token, err := CreateAccessToken(userID, role, cfg)
	if err != nil {
		t.Fatalf("CreateAccessToken failed: %v", err)
	}

	// wait a tiny bit to ensure token is expired (ExpirationHours == 0)
	time.Sleep(10 * time.Millisecond)

	_, err = ParseToken(token, cfg)
	if err == nil {
		t.Fatalf("expected token to be invalid/expired")
	}
}
