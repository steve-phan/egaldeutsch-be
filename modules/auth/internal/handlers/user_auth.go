package handlers

import (
	sharedmodels "egaldeutsch-be/pkg/models"
)

// UserService combines the minimal user-related operations required by the auth module.
// Keep the surface small and return a typed UserView for cross-module communication.
type UserService interface {
	// AuthenticateUser verifies credentials and returns user Id (UUID string) and role if valid.
	AuthenticateUser(email, password string) (string, sharedmodels.UserRole, error)

	// UpdatePassword updates a user's password (the implementation should handle hashing).
	UpdatePassword(userID string, newPassword string) error

	// GetUserIDByEmail returns the user ID string for the provided email.
	GetUserIDByEmail(email string) (string, error)

	// GetUserViewByID returns a small user view (id, name, email, role).
	GetUserViewByID(userID string) (*sharedmodels.UserView, error)
}
