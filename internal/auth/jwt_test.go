package auth

import (
	"strings"
	"testing"
	"time"

	"egaldeutsch-be/internal/config"

	"github.com/google/uuid"
)

func TestCreateAndParseToken(t *testing.T) {
	cfg := config.JwtConfig{
		SecretKey:                  "this-is-a-very-secure-secret-key-with-32-plus-characters",
		Issuer:                     "egaldeutsch",
		ExpirationHours:            1,
		RefreshTokenExpirationDays: 7,
	}

	userIDStr := uuid.New().String()
	roleStr := "learner"

	// Test the strongly typed version
	userID, err := NewUserID(userIDStr)
	if err != nil {
		t.Fatalf("NewUserID failed: %v", err)
	}

	role, err := NewRole(roleStr)
	if err != nil {
		t.Fatalf("NewRole failed: %v", err)
	}

	token, err := CreateAccessToken(userID, role, cfg)
	if err != nil {
		t.Fatalf("CreateAccessToken failed: %v", err)
	}

	claims, err := ParseToken(token, cfg)
	if err != nil {
		t.Fatalf("ParseToken failed: %v", err)
	}
	if claims.UserId != userID.String() {
		t.Fatalf("expected user id %s got %s", userID.String(), claims.UserId)
	}
	if claims.Role != role.String() {
		t.Fatalf("expected role %s got %s", role.String(), claims.Role)
	}
}

func TestCreateAccessTokenFromStrings(t *testing.T) {
	cfg := config.JwtConfig{
		SecretKey:                  "this-is-a-very-secure-secret-key-with-32-plus-characters",
		Issuer:                     "egaldeutsch",
		ExpirationHours:            1,
		RefreshTokenExpirationDays: 7,
	}

	userIDStr := uuid.New().String()
	roleStr := "learner"

	// Test the convenience string version
	token, err := CreateAccessTokenFromStrings(userIDStr, roleStr, cfg)
	if err != nil {
		t.Fatalf("CreateAccessTokenFromStrings failed: %v", err)
	}

	claims, err := ParseToken(token, cfg)
	if err != nil {
		t.Fatalf("ParseToken failed: %v", err)
	}
	if claims.UserId != userIDStr {
		t.Fatalf("expected user id %s got %s", userIDStr, claims.UserId)
	}
	if claims.Role != roleStr {
		t.Fatalf("expected role %s got %s", roleStr, claims.Role)
	}
}

func TestTokenExpiry(t *testing.T) {
	cfg := config.JwtConfig{
		SecretKey:                  "test-secret-key-that-is-long-enough-for-hs256-algorithm",
		Issuer:                     "egaldeutsch",
		ExpirationHours:            1, // 1 hour expiry
		RefreshTokenExpirationDays: 7,
	}

	userID := uuid.New().String()
	role := "learner"

	// Create token with short expiry (1 hour)
	_, err := CreateAccessTokenFromStrings(userID, role, cfg)
	if err != nil {
		t.Fatalf("failed to create token: %v", err)
	}

	// Create another token
	token, err := CreateAccessTokenFromStrings(userID, role, cfg)
	if err != nil {
		t.Fatalf("failed to create token: %v", err)
	}

	// Parse the token to check claims
	claims, err := ParseToken(token, cfg)
	if err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}

	// Check if token is not expired (it should be valid for 1 hour)
	if time.Until(claims.ExpiresAt.Time) <= 0 {
		t.Fatal("token should not be expired yet")
	}
}

// TestCreateAccessTokenValidation tests parameter validation scenarios
// Note: JWT config validation is now done at startup, so we focus on parameter validation
func TestCreateAccessTokenValidation(t *testing.T) {
	validCfg := config.JwtConfig{
		SecretKey:                  "this-is-a-very-secure-secret-key-with-32-plus-characters",
		Issuer:                     "egaldeutsch",
		ExpirationHours:            1,
		RefreshTokenExpirationDays: 7,
	}
	validUserID := uuid.New().String()
	validRole := "learner"

	tests := []struct {
		name      string
		userID    string
		role      string
		shouldErr bool
		errorMsg  string
	}{
		{
			name:      "valid parameters",
			userID:    validUserID,
			role:      validRole,
			shouldErr: false,
		},
		{
			name:      "empty userID",
			userID:    "",
			role:      validRole,
			shouldErr: true,
			errorMsg:  "invalid user ID",
		},
		{
			name:      "invalid UUID userID",
			userID:    "not-a-uuid",
			role:      validRole,
			shouldErr: true,
			errorMsg:  "invalid user ID",
		},
		{
			name:      "empty role",
			userID:    validUserID,
			role:      "",
			shouldErr: true,
			errorMsg:  "invalid role",
		},
		{
			name:      "whitespace-only role",
			userID:    validUserID,
			role:      "   ",
			shouldErr: true,
			errorMsg:  "invalid role",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CreateAccessTokenFromStrings(tt.userID, tt.role, validCfg)

			if tt.shouldErr {
				if err == nil {
					t.Fatalf("expected error containing '%s', got nil", tt.errorMsg)
				}
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Fatalf("expected error containing '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
			}
		})
	}
}

// TestJWTService tests the new JWT service approach
func TestJWTService(t *testing.T) {
	cfg := config.JwtConfig{
		SecretKey:                  "this-is-a-very-secure-secret-key-with-32-plus-characters",
		Issuer:                     "egaldeutsch",
		ExpirationHours:            1,
		RefreshTokenExpirationDays: 7,
	}

	service := NewJWTService(cfg)

	userIDStr := uuid.New().String()
	roleStr := "learner"

	// Test strongly typed method
	userID, err := NewUserID(userIDStr)
	if err != nil {
		t.Fatalf("failed to create UserID: %v", err)
	}

	role, err := NewRole(roleStr)
	if err != nil {
		t.Fatalf("failed to create Role: %v", err)
	}

	token, err := service.CreateAccessToken(userID, role)
	if err != nil {
		t.Fatalf("failed to create access token: %v", err)
	}

	// Test parsing
	claims, err := service.ParseToken(token)
	if err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}

	if claims.UserId != userIDStr {
		t.Fatalf("expected user ID %s, got %s", userIDStr, claims.UserId)
	}

	if claims.Role != roleStr {
		t.Fatalf("expected role %s, got %s", roleStr, claims.Role)
	}

	// Test string convenience method
	token2, err := service.CreateAccessTokenFromStrings(userIDStr, roleStr)
	if err != nil {
		t.Fatalf("failed to create access token from strings: %v", err)
	}

	claims2, err := service.ParseToken(token2)
	if err != nil {
		t.Fatalf("failed to parse second token: %v", err)
	}

	if claims2.UserId != userIDStr {
		t.Fatalf("expected user ID %s, got %s", userIDStr, claims2.UserId)
	}

	if claims2.Role != roleStr {
		t.Fatalf("expected role %s, got %s", roleStr, claims2.Role)
	}
}
