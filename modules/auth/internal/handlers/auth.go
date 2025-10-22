package handlers

import (
	"egaldeutsch-be/internal/auth"
	authModels "egaldeutsch-be/modules/auth/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService auth.AuthService
	userService UserService
}

func NewAuthHandler(authService auth.AuthService, userService UserService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userService: userService,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req authModels.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userId, err := h.userService.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	token, err := h.authService.CreateAccessToken(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create token"})
		return
	}
	// create refresh token
	ip := c.ClientIP()
	ua := c.Request.UserAgent()
	refresh, err := h.authService.CreateRefreshToken(userId, ip, ua)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create refresh token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"access_token": token, "refresh_token": refresh})

}

// Logout revokes the provided refresh token. Clients should call this when logging out.
func (h *AuthHandler) Logout(c *gin.Context) {
	var req authModels.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authService.RevokeRefreshToken(req.RefreshToken); err != nil {
		if err == auth.ErrInvalidRefreshToken {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"}) //TODO:  treat as success to avoid token fishing?
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to revoke token"})
		return
	}

	c.Status(http.StatusOK)
}
