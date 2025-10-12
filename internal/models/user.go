package models

// User represents a user in the system
type User struct {
	BaseModel
	Name string `json:"name" gorm:"not null"`
	Role string `json:"role" gorm:"not null;check:role IN ('admin', 'user')"`
}

type CreateUserRequest struct {
	Name string `json:"name" binding:"required,min=2,max=100"`
	Role string `json:"role" binding:"required,oneof=admin user"`
}

type UpdateUserRequest struct {
	Name string `json:"name,omitempty" binding:"omitempty,min=2,max=100"`
	Role string `json:"role,omitempty" binding:"omitempty,oneof=admin user"`
}

type UserIDParam struct {
	ID string `uri:"id" binding:"required,uuid"`
}
