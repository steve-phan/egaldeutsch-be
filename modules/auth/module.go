package authmodule

import (
	"egaldeutsch-be/internal/auth"
	"egaldeutsch-be/modules/auth/internal/handlers"
)

type Module struct {
	Handler *handlers.AuthHandler
}

func NewModule(authService auth.AuthService, userAuth handlers.UserAuthenticator) *Module {
	return &Module{Handler: handlers.NewAuthHandler(authService, userAuth)}
}
