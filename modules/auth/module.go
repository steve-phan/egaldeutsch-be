package authmodule

import (
	"egaldeutsch-be/internal/auth"
	"egaldeutsch-be/modules/auth/internal/handlers"
	"egaldeutsch-be/modules/auth/internal/models"
)

type Module struct {
	Handler *handlers.AuthHandler
}

func NewModule(authService auth.AuthService, userAuth handlers.UserAuthenticator) *Module {
	return &Module{Handler: handlers.NewAuthHandler(authService, userAuth)}
}

// GetModelsForMigration returns module models that should be auto-migrated.
func GetModelsForMigration() []interface{} {
	return []interface{}{&models.RefreshToken{}, &models.PasswordReset{}}
}
