package models

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID         uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	TokenHash  string     `gorm:"type:text;not null;uniqueIndex" json:"token_hash"`
	CreatedAt  time.Time  `gorm:"not null;default:now()" json:"created_at"`
	LastUsedAt *time.Time `json:"last_used_at"`
	ExpiresAt  time.Time  `gorm:"not null" json:"expires_at"`
	Revoked    bool       `gorm:"not null;default:false" json:"revoked"`
	ReplacedBy *string    `gorm:"type:text" json:"replaced_by"`
	IP         *string    `gorm:"type:text" json:"ip"`
	UserAgent  *string    `gorm:"type:text" json:"user_agent"`
}

type PasswordReset struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	TokenHash string    `gorm:"type:text;not null;uniqueIndex" json:"token_hash"`
	CreatedAt time.Time `gorm:"not null;default:now()" json:"created_at"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	Used      bool      `gorm:"not null;default:false" json:"used"`
}
