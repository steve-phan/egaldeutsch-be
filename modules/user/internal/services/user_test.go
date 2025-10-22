package services_test

import (
	"testing"
	"time"

	"egaldeutsch-be/modules/user/internal/models"
	"egaldeutsch-be/modules/user/internal/repositories"
	"egaldeutsch-be/modules/user/internal/services"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type testUser struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Email     string
	Password  string
	Name      string
	Role      string
}

func (testUser) TableName() string { return "users" }

func setupUserDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}

	// Create minimal users table for repository
	if err := db.AutoMigrate(&testUser{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func TestAuthenticateUser_Success(t *testing.T) {
	db := setupUserDB(t)
	repo := repositories.NewUserRepository(db)
	svc := services.NewUserService(repo)

	// create user via repo.Create
	password := "supersecret"
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("bcrypt failed: %v", err)
	}
	u := &models.User{Email: "bob@example.com", Password: string(hashed), Name: "Bob", Role: "learner"}
	if err := repo.Create(u); err != nil {
		t.Fatalf("create user: %v", err)
	}

	// Now authenticate
	id, role, err := svc.AuthenticateUser("bob@example.com", password)
	if err != nil {
		t.Fatalf("authenticate user failed: %v", err)
	}
	if id == "" {
		t.Fatalf("expected non-empty id")
	}
	if string(role) != "learner" {
		t.Fatalf("expected role learner got %s", role)
	}
}
