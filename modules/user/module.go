package user

import (
	"gorm.io/gorm"

	"egaldeutsch-be/modules/user/internal/handlers"
	"egaldeutsch-be/modules/user/internal/models"
	"egaldeutsch-be/modules/user/internal/repositories"
	"egaldeutsch-be/modules/user/internal/services"
)

// Module provides the user module functionality
type Module struct {
	Handler *handlers.UserHandler
	Service *services.UserService
	Repo    *repositories.UserRepository
}

// NewModule creates a new user module with all dependencies
func NewModule(db *gorm.DB) *Module {
	// Initialize repository
	repo := repositories.NewUserRepository(db)

	// Initialize service
	service := services.NewUserService(repo)

	// Initialize handler
	handler := handlers.NewUserHandler(service)

	return &Module{
		Handler: handler,
		Service: service,
		Repo:    repo,
	}
}

// GetModelsForMigration returns models that need to be migrated
func (m *Module) GetModelsForMigration() []interface{} {
	return []interface{}{
		&models.User{},
	}
}
