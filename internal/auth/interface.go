package auth

type AuthService interface {
	// Access token helpers
	CreateAccessToken(userID string) (string, error)
	ParseToken(token string) (*Claims, error)

	// Refresh token lifecycle
	// CreateRefreshToken returns the plain (unhashed) refresh token to be given to the client
	CreateRefreshToken(userID string, ip string, userAgent string) (string, error)
	// RefreshTokens rotates the provided refresh token and returns a new access + refresh token.
	// The service will resolve the user's role from the repository as part of rotation, so
	// callers only need to present the old refresh token and client metadata.
	RefreshTokens(oldRefreshToken string, ip string, userAgent string) (newAccess string, newRefresh string, err error)
	// Revoke a single refresh token (by presenting the plain token)
	RevokeRefreshToken(refreshToken string) error
	// RevokeAllRefreshTokens revokes all refresh tokens for a user (logout everywhere)
	RevokeAllRefreshTokens(userID string) error

	// Password reset helpers
	// CreatePasswordResetForUser creates a single-use token for the given user ID and returns the plain token (to be emailed)
	CreatePasswordResetForUser(userID string) (string, error)
	VerifyPasswordResetToken(token string) (userID string, err error)
}

// AuthRepo declares the persistence operations the auth service needs.
// Implementations live in the auth module so the module owns its models.
type AuthRepo interface {
	// Insert a new refresh token record
	InsertRefreshToken(tokenHash string, userID string, expiresAt int64, ip *string, userAgent *string) error
	// RotateRefreshToken atomically rotates the provided oldHash into a newHash.
	// It returns the associated userID, the user's role (as stored in users table),
	// and a boolean indicating whether reuse was detected (old token was already revoked and replaced_by was set).
	RotateRefreshToken(oldHash, newHash string, newExpiresAt int64, ip *string, userAgent *string) (userID string, role string, reused bool, err error)
	RevokeRefreshTokenByHash(hash string, replacedBy *string) error
	RevokeAllForUser(userID string) error

	// Password reset helpers
	InsertPasswordReset(tokenHash string, userID string, expiresAt int64) error
	VerifyAndMarkPasswordReset(tokenHash string) (userID string, err error)
}
