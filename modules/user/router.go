package user

import (
	"egaldeutsch-be/internal/config"
	"egaldeutsch-be/internal/middleware"
	"egaldeutsch-be/pkg/models"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes mounts the user module routes onto the provided API group.
// The server should keep responsibility for creating the top-level API group
// (for example: /api/v1) and pass it to modules so they can attach their
// sub-routes (for example: /users).
//
// Protection rules (recommended):
// - POST /users            : public (user sign-up)
// - GET /users/:id         : protected (authenticated users can fetch their own data, admin can fetch any)
// - PUT /users/:id         : protected (user can update their own, admin can update any)
// - DELETE /users/:id      : role-protected (admin only)
// - GET /users             : role-protected (admin only)
//
// The function accepts the JWT config so it can mount the AuthMiddleware
// consistently with the rest of the app.
func (m *Module) RegisterRoutes(rg *gin.RouterGroup, jwtCfg config.JwtConfig) {
	users := rg.Group("/users")
	{
		// public: signup
		users.POST("", m.Handler.CreateUser)

		// protected routes: require authentication
		usersAuth := users.Group("")
		usersAuth.Use(middleware.AuthMiddleware(jwtCfg))
		{

		}

		// admin-only routes
		usersAdmin := users.Group("")
		usersAdmin.Use(middleware.AuthMiddleware(jwtCfg), middleware.RequireRole(models.UserRoleAdmin))
		{
			usersAdmin.GET(":id", m.Handler.GetUser)
			usersAdmin.PUT(":id", m.Handler.UpdateUser)
			usersAdmin.DELETE(":id", m.Handler.DeleteUser)
			usersAdmin.GET("", m.Handler.ListUsers)
		}
	}
}
