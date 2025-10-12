package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BaseModel contains common fields for all models
type BaseModel struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// // BeforeCreate hook to generate UUID if not set
// func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
// 	if b.ID == uuid.Nil {
// 		b.ID = uuid.New()
// 	}
// 	return nil
// }
