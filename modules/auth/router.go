package authmodule

import "github.com/gin-gonic/gin"

func (m *Module) RegisterRoutes(rg *gin.RouterGroup) {

	ag := rg.Group("/auth")
	{
		ag.POST("/login", m.Handler.Login)
		// ag.POST("/logout", m.Handler.Logout)
		// ag.POST("/refresh", m.Handler.RefreshToken)
		// ag.POST("/forgot-password", m.Handler.ForgotPassword)
		// ag.POST("/reset-password", m.Handler.ResetPassword)
		// ag.POST("revoke-all-tokens", m.Handler.RevokeAllTokens)

		// ag.GET("/me", m.Handler.GetCurrentUser)
	}
}
