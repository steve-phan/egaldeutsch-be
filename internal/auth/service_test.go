package auth

import (
	"errors"
	"testing"

	"egaldeutsch-be/internal/config"
)

// fakeRepo implements AuthRepo for testing
type fakeRepo struct {
	inserted         map[string]bool
	nextRotateUserID string
	nextRotateRole   string
	nextRotateReused bool
	rotateErr        error
}

func (f *fakeRepo) InsertRefreshToken(tokenHash string, userID string, expiresAt int64, ip *string, userAgent *string) error {
	if f.inserted == nil {
		f.inserted = map[string]bool{}
	}
	f.inserted[tokenHash] = true
	return nil
}
func (f *fakeRepo) RotateRefreshToken(oldHash, newHash string, newExpiresAt int64, ip *string, userAgent *string) (string, string, bool, error) {
	return f.nextRotateUserID, f.nextRotateRole, f.nextRotateReused, f.rotateErr
}
func (f *fakeRepo) RevokeRefreshTokenByHash(hash string, replacedBy *string) error { return nil }
func (f *fakeRepo) RevokeAllForUser(userID string) error                           { return nil }
func (f *fakeRepo) InsertPasswordReset(tokenHash string, userID string, expiresAt int64) error {
	return nil
}
func (f *fakeRepo) VerifyAndMarkPasswordReset(tokenHash string) (string, error) { return "", nil }

func TestCreateRefreshTokenAndRotateSuccess(t *testing.T) {
	repo := &fakeRepo{}
	cfg := config.JwtConfig{SecretKey: "test", Issuer: "eg", ExpirationHours: 1, RefreshTokenExpirationDays: 30}
	svc := NewService(cfg, repo)

	// CreateRefreshToken should call repo.InsertRefreshToken
	plain, err := svc.CreateRefreshToken("user-1", "1.2.3.4", "ua")
	if err != nil {
		t.Fatalf("CreateRefreshToken failed: %v", err)
	}
	if plain == "" {
		t.Fatalf("expected plain token")
	}

	// Rotate scenario: repo returns userID+role
	repo.nextRotateUserID = "user-1"
	repo.nextRotateRole = "learner"
	repo.nextRotateReused = false

	access, newPlain, err := svc.RefreshTokens("oldtoken", "1.2.3.4", "ua")
	if err != nil {
		t.Fatalf("RefreshTokens failed: %v", err)
	}
	if access == "" || newPlain == "" {
		t.Fatalf("expected access and refresh tokens")
	}
}

func TestRefreshTokens_ReuseDetected(t *testing.T) {
	repo := &fakeRepo{rotateErr: nil, nextRotateReused: true}
	cfg := config.JwtConfig{SecretKey: "test", Issuer: "eg", ExpirationHours: 1, RefreshTokenExpirationDays: 30}
	svc := NewService(cfg, repo)

	_, _, err := svc.RefreshTokens("oldtoken", "1.2.3.4", "ua")
	if !errors.Is(err, ErrRefreshTokenReuse) {
		t.Fatalf("expected ErrRefreshTokenReuse got %v", err)
	}
}
