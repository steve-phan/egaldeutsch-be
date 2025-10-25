package handlers

import (
	sharedmodels "egaldeutsch-be/pkg/models"
)

// UserAuthenticator handles user authentication operations.
// Small interface following Go philosophy: do one thing well.
type UserAuthenticator interface {
	// AuthenticateUser verifies credentials and returns user ID (UUID string) and role if valid.
	AuthenticateUser(email, password string) (string, sharedmodels.UserRole, error)
}

// UserPasswordManager handles password-related operations.
type UserPasswordManager interface {
	// UpdatePassword updates a user's password (the implementation should handle hashing).
	UpdatePassword(userID string, newPassword string) error
}

// UserLookup handles user lookup operations.
type UserLookup interface {
	// GetUserIDByEmail returns the user ID string for the provided email.
	GetUserIDByEmail(email string) (string, error)

	// GetUserViewByID returns a small user view (id, name, email, role).
	GetUserViewByID(userID string) (*sharedmodels.UserView, error)
}

// UserService combines minimal user-related operations required by the auth module.
// Composed interface following Go's composition principle - keep the surface small.
type UserService interface {
	UserAuthenticator
	UserPasswordManager
	UserLookup
}
