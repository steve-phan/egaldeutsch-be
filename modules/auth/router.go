package authmodule

import (
	"egaldeutsch-be/internal/config"
	"egaldeutsch-be/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers both public and protected routes for the auth module.
// Protected routes (like /me) are registered under the module and are guarded by the
// provided JWT config via the auth middleware.
func (m *Module) RegisterRoutes(rg *gin.RouterGroup, jwtCfg config.JwtConfig) {

	ag := rg.Group("/auth")
	{
		ag.POST("/login", m.Handler.Login)
		ag.POST("/logout", m.Handler.Logout)
		ag.POST("/refresh", m.Handler.RefreshToken)
		ag.POST("/forgot-password", m.Handler.ForgotPassword)
		ag.POST("/reset-password", m.Handler.ResetPassword)
	}

	// Protected routes inside the module
	agProtected := rg.Group("/auth")
	agProtected.Use(middleware.AuthMiddleware(jwtCfg))
	{
		agProtected.GET("/me", middleware.RateLimitMiddleware(middleware.RateLimitConfig{
			RequestsPerMinute: 10,
		}), m.Handler.GetCurrentUser)
	}

}
