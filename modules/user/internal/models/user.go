package models

import (
	"egaldeutsch-be/pkg/models"
)

// User represents a user in the system
type User struct {
	models.BaseModel
	Name string          `json:"name" gorm:"not null;size:100"`
	Role models.UserRole `json:"role" gorm:"not null;type:varchar(20);check:role IN ('admin', 'user')"`
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}
