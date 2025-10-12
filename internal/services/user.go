package services

import (
	"errors"

	"egaldeutsch-be/internal/models"
	"egaldeutsch-be/internal/repositories"

	"github.com/google/uuid"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(name, role string) (*models.User, error) {
	// Business logic validation
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}
	if role != "admin" && role != "user" {
		return nil, errors.New("invalid role")
	}

	user := &models.User{
		Name: name,
		Role: role,
	}

	err := s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUserByID(id string) (*models.User, error) {
	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return nil, errors.New("invalid user ID format")
	}

	return s.userRepo.GetByID(id)
}

func (s *UserService) UpdateUserRole(id, newRole string) error {
	// Business logic: only allow role changes to valid roles
	if newRole != "admin" && newRole != "user" {
		return errors.New("invalid role")
	}

	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	user.Role = newRole
	return s.userRepo.Update(user)
}

func (s *UserService) GetUsersByRole(role string, limit, offset int) ([]models.User, error) {
	// Business logic: validate role if provided
	if role != "" && role != "admin" && role != "user" {
		return nil, errors.New("invalid role")
	}

	return s.userRepo.List(limit, offset)
}

func (s *UserService) DeleteUser(id string) error {
	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return errors.New("invalid user ID format")
	}

	return s.userRepo.Delete(id)
}
