package services

import (
	"errors"

	"github.com/google/uuid"

	usermodels "egaldeutsch-be/modules/user/internal/models"
	"egaldeutsch-be/modules/user/internal/repositories"
	"egaldeutsch-be/pkg/models"
)

var (
	ErrInvalidUserRole = errors.New("invalid user role")
	ErrUserValidation  = errors.New("user validation failed")
)

// UserService handles business logic for users
type UserService struct {
	repo *repositories.UserRepository
}

// NewUserService creates a new user service
func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// CreateUser creates a new user with validation
func (s *UserService) CreateUser(name, roleStr string) (*usermodels.User, error) {
	// Validate role
	role := models.UserRole(roleStr)
	if !role.IsValid() {
		return nil, ErrInvalidUserRole
	}

	// Create user
	user := &usermodels.User{
		Name: name,
		Role: role,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(id string) (*usermodels.User, error) {
	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return nil, errors.New("invalid user ID format")
	}

	return s.repo.GetByID(id)
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(id string, updates *usermodels.UpdateUserRequest) (*usermodels.User, error) {
	// Get existing user
	user, err := s.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if updates.Name != "" {
		user.Name = updates.Name
	}
	if updates.Role != "" {
		role := models.UserRole(updates.Role)
		if !role.IsValid() {
			return nil, ErrInvalidUserRole
		}
		user.Role = role
	}

	// Save updates
	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser soft deletes a user
func (s *UserService) DeleteUser(id string) error {
	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return errors.New("invalid user ID format")
	}

	return s.repo.Delete(id)
}

// ListUsers retrieves users with pagination
func (s *UserService) ListUsers(page, perPage int) ([]usermodels.User, int64, error) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 || perPage > 100 {
		perPage = 10
	}

	return s.repo.List(page, perPage)
}

// GetUsersByRole retrieves users by role
func (s *UserService) GetUsersByRole(role models.UserRole) ([]usermodels.User, error) {
	if !role.IsValid() {
		return nil, ErrInvalidUserRole
	}

	return s.repo.GetByRole(role)
}
