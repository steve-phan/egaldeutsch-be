package repositories_test

import (
	"testing"
	"time"

	"egaldeutsch-be/internal/auth"
	"egaldeutsch-be/internal/config"
	"egaldeutsch-be/modules/auth/internal/models"
	"egaldeutsch-be/modules/auth/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// test-only structs to avoid Postgres-specific DDL produced by the production models
type testRefreshToken struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	TokenHash  string
	CreatedAt  time.Time
	LastUsedAt *time.Time
	ExpiresAt  time.Time
	Revoked    bool
	ReplacedBy *string
	IP         *string
	UserAgent  *string
}

func (testRefreshToken) TableName() string { return "refresh_tokens" }

type testPasswordReset struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	CreatedAt time.Time
	ExpiresAt time.Time
	Used      bool
}

func (testPasswordReset) TableName() string { return "password_resets" }

type testUser struct {
	ID       uuid.UUID `gorm:"primaryKey"`
	Email    string
	Password string
	Name     string
	Role     string
}

func (testUser) TableName() string { return "users" }

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	// create tables using test-only structs to avoid Postgres-specific DDL
	if err := db.AutoMigrate(&testRefreshToken{}, &testPasswordReset{}, &testUser{}); err != nil {
		t.Fatalf("failed to automigrate: %v", err)
	}

	return db
}

func insertUserWithRole(t *testing.T, db *gorm.DB, id uuid.UUID, role string) {
	t.Helper()
	type user struct {
		ID   uuid.UUID `gorm:"type:uuid;primaryKey"`
		Role string
	}

	u := user{ID: id, Role: role}
	if err := db.Create(&u).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
}

func TestRotateRefreshToken_Success(t *testing.T) {
	db := setupTestDB(t)
	r := repositories.NewRepository(db)

	uid := uuid.New()
	insertUserWithRole(t, db, uid, "learner")

	// insert an existing refresh token
	oldHash := "oldhash"
	rt := &models.RefreshToken{
		ID:        uuid.New(),
		UserID:    uid,
		TokenHash: oldHash,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	if err := db.Create(rt).Error; err != nil {
		t.Fatalf("failed to create refresh token: %v", err)
	}

	newHash := "newhash"
	userID, role, reused, err := r.RotateRefreshToken(oldHash, newHash, time.Now().Add(24*time.Hour).Unix(), nil, nil)
	if err != nil {
		t.Fatalf("rotate failed: %v", err)
	}
	if reused {
		t.Fatalf("expected reused=false")
	}
	if userID != uid.String() {
		t.Fatalf("expected userID %s got %s", uid.String(), userID)
	}
	if role != "learner" {
		t.Fatalf("expected role learner got %s", role)
	}

	// ensure new token exists and old is revoked
	var newRT models.RefreshToken
	if err := db.Where("token_hash = ?", newHash).First(&newRT).Error; err != nil {
		t.Fatalf("new token not found: %v", err)
	}
	var old models.RefreshToken
	if err := db.Where("token_hash = ?", oldHash).First(&old).Error; err != nil {
		t.Fatalf("old token not found: %v", err)
	}
	if !old.Revoked {
		t.Fatalf("old token should be revoked")
	}
}

func TestRotateRefreshToken_ReuseDetected(t *testing.T) {
	db := setupTestDB(t)
	r := repositories.NewRepository(db)

	uid := uuid.New()
	insertUserWithRole(t, db, uid, "teacher")

	// insert a refresh token that is already revoked
	oldHash := "revokedhash"
	replaced := "somehash"
	rt := &models.RefreshToken{
		ID:         uuid.New(),
		UserID:     uid,
		TokenHash:  oldHash,
		ExpiresAt:  time.Now().Add(24 * time.Hour),
		Revoked:    true,
		ReplacedBy: &replaced,
	}
	if err := db.Create(rt).Error; err != nil {
		t.Fatalf("failed to create revoked refresh token: %v", err)
	}

	newHash := "anothernewhash"
	userID, role, reused, err := r.RotateRefreshToken(oldHash, newHash, time.Now().Add(24*time.Hour).Unix(), nil, nil)
	if err != nil {
		t.Fatalf("rotate failed: %v", err)
	}
	if !reused {
		t.Fatalf("expected reused=true")
	}
	if userID != uid.String() {
		t.Fatalf("expected userID %s got %s", uid.String(), userID)
	}
	if role != "teacher" {
		t.Fatalf("expected role teacher got %s", role)
	}
}

// Integration-style test: user login -> create refresh token -> rotate refresh -> verify access token claims
func TestLoginThenRefreshFlow(t *testing.T) {
	db := setupTestDB(t)
	uid := uuid.New()
	insertUserWithRole(t, db, uid, "learner")

	// simulate login: we will not hash for test simplicity; instead insert refresh token directly and call RotateRefreshToken
	oldHash := "oldhashlogin"
	rt := &models.RefreshToken{
		ID:        uuid.New(),
		UserID:    uid,
		TokenHash: oldHash,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}
	if err := db.Create(rt).Error; err != nil {
		t.Fatalf("failed to create refresh token: %v", err)
	}

	// rotate
	newHash := "rotatedhash"
	userID, role, reused, err := repositories.NewRepository(db).RotateRefreshToken(oldHash, newHash, time.Now().Add(24*time.Hour).Unix(), nil, nil)
	if err != nil {
		t.Fatalf("rotate failed: %v", err)
	}
	if reused {
		t.Fatalf("unexpected reuse")
	}
	if userID != uid.String() {
		t.Fatalf("expected user id %s got %s", uid.String(), userID)
	}
	if role != "learner" {
		t.Fatalf("expected role learner got %s", role)
	}

	// create access token using internal auth CreateAccessToken with a minimal jwt config
	jwtCfg := config.JwtConfig{SecretKey: "testsecret", Issuer: "egaldeutsch", ExpirationHours: 24, RefreshTokenExpirationDays: 30}
	token, err := auth.CreateAccessToken(userID, role, jwtCfg)
	if err != nil {
		t.Fatalf("failed to create access token: %v", err)
	}

	// parse token
	claims, err := auth.ParseToken(token, jwtCfg)
	if err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}
	if claims.UserId != userID {
		t.Fatalf("claims user id mismatch: expected %s got %s", userID, claims.UserId)
	}
	if claims.Role != role {
		t.Fatalf("claims role mismatch: expected %s got %s", role, claims.Role)
	}
}
