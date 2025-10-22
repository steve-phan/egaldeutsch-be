package handlers

// UserService combines the minimal user-related operations required by the auth module.
// Keep the surface small and return primitive types (string userId) to avoid importing
// user domain models into the auth module.
type UserService interface {
	// AuthenticateUser verifies credentials and returns user Id (UUID string) if valid.
	AuthenticateUser(email, password string) (string, error)

	// UpdatePassword updates a user's password (the implementation should handle hashing).
	UpdatePassword(userID string, newPassword string) error

	// GetUserIDByEmail returns the user ID string for the provided email.
	GetUserIDByEmail(email string) (string, error)

	// GetUserViewByID returns a small user view (id, name, email, role) to avoid importing full user models.
	GetUserViewByID(userID string) (map[string]interface{}, error)
}
