package user

import "github.com/gin-gonic/gin"

// RegisterRoutes mounts the user module routes onto the provided API group.
// The server should keep responsibility for creating the top-level API group
// (for example: /api/v1) and pass it to modules so they can attach their
// sub-routes (for example: /users).
func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	{
		users.POST("", m.Handler.CreateUser)
		users.GET("/:id", m.Handler.GetUser)
		users.PUT("/:id", m.Handler.UpdateUser)
		users.DELETE("/:id", m.Handler.DeleteUser)
		users.GET("", m.Handler.ListUsers)
	}
}
