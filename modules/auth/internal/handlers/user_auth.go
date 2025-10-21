package handlers

// UserAuthenticator describes just what the auth module needs from the  user service.
// Keep it minimal and return a primitive (string userId) to avoid importing user models.
type UserAuthenticator interface {
	// AuthenticateUser verifies credentials and returns user Id (UUID string) if valid.
	AuthenticateUser(email, password string) (string, error)
}
