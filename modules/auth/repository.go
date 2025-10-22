package authmodule

import (
	"gorm.io/gorm"

	"egaldeutsch-be/internal/auth"
	repos "egaldeutsch-be/modules/auth/internal/repositories"
)

// NewRepository constructs the concrete repository implementing auth.AuthRepo
// This factory lives in the module package so callers outside the module can
// obtain an AuthRepo without importing module internals.
func NewRepository(db *gorm.DB) auth.AuthRepo {
	return repos.NewRepository(db)
}
