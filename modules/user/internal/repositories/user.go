package repositories

import (
	"errors"

	"gorm.io/gorm"

	usermodels "egaldeutsch-be/modules/user/internal/models"
	"egaldeutsch-be/pkg/models"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

// UserRepository handles database operations for users
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user in the database
func (r *UserRepository) Create(user *usermodels.User) error {
	return r.db.Create(user).Error
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id string) (*usermodels.User, error) {
	var user usermodels.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(email string) (*usermodels.User, error) {
	var user usermodels.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// Update updates an existing user
func (r *UserRepository) Update(user *usermodels.User) error {
	return r.db.Save(user).Error
}

// Delete soft deletes a user by ID
func (r *UserRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&usermodels.User{}).Error
}

// List retrieves users with pagination
func (r *UserRepository) List(page, perPage int) ([]usermodels.User, int64, error) {
	var users []usermodels.User
	var total int64

	// Count total records
	if err := r.db.Model(&usermodels.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * perPage
	err := r.db.Offset(offset).Limit(perPage).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetByRole retrieves users by role
func (r *UserRepository) GetByRole(role models.UserRole) ([]usermodels.User, error) {
	var users []usermodels.User
	err := r.db.Where("role = ?", role).Find(&users).Error
	return users, err
}
