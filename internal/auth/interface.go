package auth

// TokenValidator handles JWT token validation and parsing.
// Small interface following Go philosophy: "interfaces should be small".
type TokenValidator interface {
	ParseToken(token string) (*Claims, error)
}

// TokenCreator handles JWT access token creation.
type TokenCreator interface {
	CreateAccessToken(userID string) (string, error)
}

// RefreshTokenManager handles refresh token lifecycle operations.
type RefreshTokenManager interface {
	// CreateRefreshToken returns the plain (unhashed) refresh token to be given to the client
	CreateRefreshToken(userID string, ip string, userAgent string) (string, error)

	// RefreshTokens rotates the provided refresh token and returns a new access + refresh token.
	// The service will resolve the user's role from the repository as part of rotation, so
	// callers only need to present the old refresh token and client metadata.
	RefreshTokens(oldRefreshToken string, ip string, userAgent string) (newAccess string, newRefresh string, err error)

	// RevokeRefreshToken revokes a single refresh token (by presenting the plain token)
	RevokeRefreshToken(refreshToken string) error

	// RevokeAllRefreshTokens revokes all refresh tokens for a user (logout everywhere)
	RevokeAllRefreshTokens(userID string) error
}

// PasswordResetManager handles password reset token operations.
type PasswordResetManager interface {
	// CreatePasswordResetForUser creates a single-use token for the given user ID and returns the plain token (to be emailed)
	CreatePasswordResetForUser(userID string) (string, error)

	// VerifyPasswordResetToken verifies and consumes a password reset token, returning the user ID
	VerifyPasswordResetToken(token string) (userID string, err error)
}

// AuthService combines all authentication-related operations.
// Composed of smaller interfaces following Go's composition principle.
type AuthService interface {
	TokenValidator
	TokenCreator
	RefreshTokenManager
	PasswordResetManager
}

// RefreshTokenRepo handles refresh token persistence operations.
type RefreshTokenRepo interface {
	// InsertRefreshToken creates a new refresh token record
	InsertRefreshToken(tokenHash string, userID string, expiresAt int64, ip *string, userAgent *string) error

	// RotateRefreshToken atomically rotates the provided oldHash into a newHash.
	// Returns: userID, role (from users table), reused (true if token was already revoked), error
	RotateRefreshToken(oldHash, newHash string, newExpiresAt int64, ip *string, userAgent *string) (userID string, role string, reused bool, err error)

	// RevokeRefreshTokenByHash marks a refresh token as revoked
	RevokeRefreshTokenByHash(hash string, replacedBy *string) error

	// RevokeAllForUser revokes all refresh tokens for a specific user
	RevokeAllForUser(userID string) error
}

// PasswordResetRepo handles password reset token persistence operations.
type PasswordResetRepo interface {
	// InsertPasswordReset creates a new password reset token record
	InsertPasswordReset(tokenHash string, userID string, expiresAt int64) error

	// VerifyAndMarkPasswordReset atomically verifies and marks a password reset token as used
	VerifyAndMarkPasswordReset(tokenHash string) (userID string, err error)
}

// AuthRepo combines all authentication repository operations.
// Composed interface following Go's interface composition principle.
type AuthRepo interface {
	RefreshTokenRepo
	PasswordResetRepo
}
