package auth

type AuthService interface {
	CreateAccessToken(userID string) (string, error)
	ParseToken(token string) (*Claims, error)
}
