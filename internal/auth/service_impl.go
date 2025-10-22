package auth

import (
	"errors"
	"time"

	"egaldeutsch-be/internal/config"

	"github.com/google/uuid"
)

var (
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrRefreshTokenReuse   = errors.New("refresh token reuse detected")
)

type service struct {
	cfg  config.JwtConfig
	repo AuthRepo
}

func NewService(cfg config.JwtConfig, repo AuthRepo) AuthService {
	return &service{cfg: cfg, repo: repo}
}

func (s *service) CreateAccessToken(userID string) (string, error) {
	return CreateAccessToken(userID, s.cfg)
}

func (s *service) ParseToken(token string) (*Claims, error) {
	return ParseToken(token, s.cfg)
}

func (s *service) CreateRefreshToken(userID string, ip string, userAgent string) (string, error) {
	// generate plain token
	token, err := genRandomToken(DefaultRefreshTokenBytes)
	if err != nil {
		return "", err
	}
	hash := hashToken(token)

	// Use repo interface to persist
	expiresAtUnix := time.Now().Add(time.Duration(s.cfg.RefreshTokenExpirationDays*24) * time.Hour).Unix()
	if err := s.repo.InsertRefreshToken(hash, userID, expiresAtUnix, &ip, &userAgent); err != nil {
		return "", err
	}

	return token, nil
}

func (s *service) RefreshTokens(oldRefreshToken string, ip string, userAgent string) (string, string, error) {
	// hash provided token
	oldHash := hashToken(oldRefreshToken)

	// Attempt rotation via repo
	newToken, err := genRandomToken(DefaultRefreshTokenBytes)
	if err != nil {
		return "", "", err
	}
	newHash := hashToken(newToken)

	expiresAtUnix := time.Now().Add(time.Duration(s.cfg.RefreshTokenExpirationDays*24) * time.Hour).Unix()
	userID, reused, err := s.repo.RotateRefreshToken(oldHash, newHash, expiresAtUnix, &ip, &userAgent)
	if err != nil {
		return "", "", err
	}
	if reused {
		return "", "", ErrRefreshTokenReuse
	}

	// create access token
	access, err := s.CreateAccessToken(userID)
	if err != nil {
		return "", "", err
	}

	return access, newToken, nil
}

func (s *service) RevokeRefreshToken(refreshToken string) error {
	hash := hashToken(refreshToken)
	return s.repo.RevokeRefreshTokenByHash(hash, nil)
}

func (s *service) RevokeAllRefreshTokens(userID string) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	return s.repo.RevokeAllForUser(uid.String())
}

func (s *service) CreatePasswordResetForEmail(email string) error {
	return errors.New("not implemented")
}

func (s *service) VerifyPasswordResetToken(token string) (string, error) {
	return "", errors.New("not implemented")
}
