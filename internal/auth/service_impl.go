package auth

import (
	"egaldeutsch-be/internal/config"
)

type service struct {
	cfg config.JwtConfig
}

func NewService(cfg config.JwtConfig) AuthService {
	return &service{cfg: cfg}
}

func (s *service) CreateAccessToken(userID string) (string, error) {
	return CreateAccessToken(userID, s.cfg)
}
func (s *service) ParseToken(token string) (*Claims, error) {
	return ParseToken(token, s.cfg)
}
