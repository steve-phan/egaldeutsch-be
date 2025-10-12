package models

// CreateUserRequest represents the request payload for creating a user
type CreateUserRequest struct {
	Name string `json:"name" binding:"required,min=2,max=100"`
	Role string `json:"role" binding:"required,oneof=admin user"`
}

// UpdateUserRequest represents the request payload for updating a user
type UpdateUserRequest struct {
	Name string `json:"name,omitempty" binding:"omitempty,min=2,max=100"`
	Role string `json:"role,omitempty" binding:"omitempty,oneof=admin user"`
}

// UserIDParam represents the URI parameter for user ID
type UserIDParam struct {
	ID string `uri:"id" binding:"required,uuid"`
}

// ListUsersQuery represents query parameters for listing users
type ListUsersQuery struct {
	Page    int `form:"page" binding:"omitempty,min=1"`
	PerPage int `form:"per_page" binding:"omitempty,min=1,max=100"`
}
