package authmodule

import "github.com/gin-gonic/gin"

func (m *Module) RegisterRoutes(rg *gin.RouterGroup) {
	ag := rg.Group("/auth")
	{
		ag.POST("/login", m.Handler.Login)
	}
}
