package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BaseModel contains common fields for all domain models
// This provides consistent ID generation, timestamps, and soft deletes
type BaseModel struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// BeforeCreate hook to ensure UUID is generated
func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

// UserRole represents the role of a user in the system
type UserRole string

const (
	UserRoleAdmin UserRole = "admin"
	UserRoleUser  UserRole = "user"
)

// IsValid checks if the user role is valid
func (r UserRole) IsValid() bool {
	switch r {
	case UserRoleAdmin, UserRoleUser:
		return true
	default:
		return false
	}
}

// String returns the string representation of the role
func (r UserRole) String() string {
	return string(r)
}

// Scan implements the sql.Scanner interface for database reading
func (r *UserRole) Scan(value interface{}) error {
	if value == nil {
		*r = ""
		return nil
	}

	switch s := value.(type) {
	case string:
		*r = UserRole(s)
	case []byte:
		*r = UserRole(s)
	}
	return nil
}

// Value implements the driver.Valuer interface for database writing
func (r UserRole) Value() (interface{}, error) {
	return string(r), nil
}

// LanguageLevel represents CEFR language proficiency levels
type LanguageLevel string

const (
	LevelA1 LanguageLevel = "A1"
	LevelA2 LanguageLevel = "A2"
	LevelB1 LanguageLevel = "B1"
	LevelB2 LanguageLevel = "B2"
	LevelC1 LanguageLevel = "C1"
)

// IsValid checks if the language level is valid
func (l LanguageLevel) IsValid() bool {
	switch l {
	case LevelA1, LevelA2, LevelB1, LevelB2, LevelC1:
		return true
	default:
		return false
	}
}

// String returns the string representation of the level
func (l LanguageLevel) String() string {
	return string(l)
}
