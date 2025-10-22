package authmodule

import (
	"egaldeutsch-be/internal/auth"
	"egaldeutsch-be/internal/config"
	"egaldeutsch-be/modules/auth/internal/handlers"
	"egaldeutsch-be/modules/auth/internal/models"
)

type Module struct {
	Handler *handlers.AuthHandler
}

func NewModule(authService auth.AuthService, userService handlers.UserService, jwtCfg config.JwtConfig) *Module {
	return &Module{Handler: handlers.NewAuthHandler(authService, userService, jwtCfg)}
}

// GetModelsForMigration returns module models that should be auto-migrated.
func GetModelsForMigration() []interface{} {
	return []interface{}{&models.RefreshToken{}, &models.PasswordReset{}}
}
