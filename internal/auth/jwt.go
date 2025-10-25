package auth

import (
	"egaldeutsch-be/internal/config"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// UserID represents a validated user identifier.
// Following Go philosophy: make invalid states unrepresentable.
type UserID struct {
	id uuid.UUID
}

// NewUserID creates a UserID from a string, validating it's a proper UUID.
func NewUserID(id string) (UserID, error) {
	if id == "" {
		return UserID{}, fmt.Errorf("user ID cannot be empty")
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return UserID{}, fmt.Errorf("invalid user ID format: %w", err)
	}

	return UserID{id: parsedID}, nil
}

// String returns the string representation of the UserID.
func (u UserID) String() string {
	return u.id.String()
}

// Role represents a validated user role.
// Following Go philosophy: explicit types prevent parameter confusion.
type Role struct {
	value string
}

// NewRole creates a Role from a string, validating it's properly formatted.
func NewRole(role string) (Role, error) {
	if role == "" {
		return Role{}, fmt.Errorf("role cannot be empty")
	}

	// Normalize whitespace
	normalized := strings.TrimSpace(role)
	if normalized == "" {
		return Role{}, fmt.Errorf("role cannot be only whitespace")
	}

	// Optional: strict validation (uncomment if needed)
	// validRoles := map[string]bool{"admin": true, "user": true, "learner": true, "teacher": true}
	// if !validRoles[normalized] {
	//     return Role{}, fmt.Errorf("unknown role: %s", normalized)
	// }

	return Role{value: normalized}, nil
}

// String returns the string representation of the Role.
func (r Role) String() string {
	return r.value
}

// Claims represents JWT token claims with strong typing.
type Claims struct {
	UserId string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// JWTService encapsulates JWT operations with pre-validated configuration.
// Following Go philosophy: encapsulate configuration, provide simple methods.
type JWTService struct {
	config config.JwtConfig
}

// NewJWTService creates a new JWT service with validated configuration.
// The config is assumed to be pre-validated during application startup.
func NewJWTService(jwtConfig config.JwtConfig) *JWTService {
	return &JWTService{
		config: jwtConfig,
	}
}

// CreateAccessToken creates a JWT access token with strongly typed parameters.
// Following Go philosophy: make invalid states unrepresentable, explicit types.
func (j *JWTService) CreateAccessToken(userID UserID, role Role) (string, error) {
	claims := j.createClaims(userID, role)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(j.config.SecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT token: %w", err)
	}

	return signedToken, nil
}

// CreateAccessTokenFromStrings creates a JWT access token from string parameters.
// Convenience function for backward compatibility and cases where you have strings.
func (j *JWTService) CreateAccessTokenFromStrings(userIDStr, roleStr string) (string, error) {
	userID, err := NewUserID(userIDStr)
	if err != nil {
		return "", fmt.Errorf("invalid user ID: %w", err)
	}

	role, err := NewRole(roleStr)
	if err != nil {
		return "", fmt.Errorf("invalid role: %w", err)
	}

	return j.CreateAccessToken(userID, role)
}

// ParseToken parses and validates a JWT token string, returning the claims.
// Following Go philosophy: explicit validation, clear error messages.
func (j *JWTService) ParseToken(tokenString string) (*Claims, error) {
	if tokenString == "" {
		return nil, fmt.Errorf("token string cannot be empty")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %s", t.Method.Alg())
		}
		return []byte(j.config.SecretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims type")
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	return claims, nil
}

// createClaims creates JWT claims from validated parameters.
// Separated from CreateAccessToken following single responsibility principle.
func (j *JWTService) createClaims(userID UserID, role Role) Claims {
	now := time.Now()
	expiresAt := now.Add(time.Duration(j.config.ExpirationHours) * time.Hour)

	return Claims{
		UserId: userID.String(),
		Role:   role.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			Issuer:    j.config.Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
}

// Legacy functions for backward compatibility
// These will be deprecated in favor of the JWTService methods

// Legacy functions for backward compatibility
// These will be deprecated in favor of the JWTService methods

// CreateAccessToken creates a JWT access token with strongly typed parameters.
// Deprecated: Use JWTService.CreateAccessToken instead for better ergonomics.
// Following Go philosophy: make invalid states unrepresentable, explicit types.
// Note: jwtConfig is assumed to be pre-validated during application startup.
func CreateAccessToken(userID UserID, role Role, jwtConfig config.JwtConfig) (string, error) {
	service := NewJWTService(jwtConfig)
	return service.CreateAccessToken(userID, role)
}

// CreateAccessTokenFromStrings creates a JWT access token from string parameters.
// Deprecated: Use JWTService.CreateAccessTokenFromStrings instead for better ergonomics.
// Convenience function for backward compatibility and cases where you have strings.
func CreateAccessTokenFromStrings(userIDStr, roleStr string, jwtConfig config.JwtConfig) (string, error) {
	service := NewJWTService(jwtConfig)
	return service.CreateAccessTokenFromStrings(userIDStr, roleStr)
}

// ParseToken parses and validates a JWT token string, returning the claims.
// Deprecated: Use JWTService.ParseToken instead for better ergonomics.
// Following Go philosophy: explicit validation, clear error messages.
func ParseToken(tokenString string, jwtConfig config.JwtConfig) (*Claims, error) {
	service := NewJWTService(jwtConfig)
	return service.ParseToken(tokenString)
}
