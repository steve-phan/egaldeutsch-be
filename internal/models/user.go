package models

// User represents a user in the system
type User struct {
	BaseModel
	Name string `json:"name" gorm:"not null"`
	Role string `json:"role" gorm:"not null;check:role IN ('admin', 'user')"`
}
