package services

import (
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	usermodels "egaldeutsch-be/modules/user/internal/models"
	"egaldeutsch-be/modules/user/internal/repositories"
	sharedmodels "egaldeutsch-be/pkg/models"
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
func (s *UserService) CreateUser(req *usermodels.CreateUserRequest) (*usermodels.User, error) {
	// Extract fields from request
	name := req.Name
	email := req.Email
	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	role := req.Role
	if role == "" {
		role = "user"
	}

	user := &usermodels.User{
		Name:     name,
		Email:    email,
		Password: string(password),
		Role:     sharedmodels.UserRole(role),
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
		role := sharedmodels.UserRole(updates.Role)
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
func (s *UserService) GetUsersByRole(role sharedmodels.UserRole) ([]usermodels.User, error) {
	if !role.IsValid() {
		return nil, ErrInvalidUserRole
	}

	return s.repo.GetByRole(role)
}

// AuthenticateUser authenticates by email and password and returns the userId
func (s *UserService) AuthenticateUser(email, password string) (string, error) {
	user, err := s.repo.GetByEmail(email)

	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	return user.ID.String(), nil
}

// GetByEmail retrieves a user by email (public wrapper)
func (s *UserService) GetByEmail(email string) (*usermodels.User, error) {
	return s.repo.GetByEmail(email)
}

// GetUserIDByEmail returns user ID (string) by email
func (s *UserService) GetUserIDByEmail(email string) (string, error) {
	u, err := s.repo.GetByEmail(email)
	if err != nil {
		return "", err
	}
	return u.ID.String(), nil
}

// UpdatePassword updates the user's password with a bcrypt hash
func (s *UserService) UpdatePassword(userID string, newPassword string) error {
	user, err := s.GetUserByID(userID)
	if err != nil {
		return err
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(passwordHash)
	return s.repo.Update(user)
}

// GetUserViewByID returns a minimal user representation for external modules.
func (s *UserService) GetUserViewByID(userID string) (*sharedmodels.UserView, error) {
	u, err := s.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	return &sharedmodels.UserView{
		ID:    u.ID.String(),
		Name:  u.Name,
		Email: u.Email,
		Role:  u.Role,
	}, nil
}
