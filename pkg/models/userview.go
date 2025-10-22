package models

// UserView is a small DTO used across modules where a full user model is not required.
type UserView struct {
	ID    string      `json:"id"`
	Name  string      `json:"name"`
	Email string      `json:"email"`
	Role  interface{} `json:"role"`
}
