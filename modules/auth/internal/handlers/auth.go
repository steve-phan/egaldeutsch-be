package handlers

import (
	"egaldeutsch-be/internal/auth"
	"egaldeutsch-be/modules/auth/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService auth.AuthService
	userAuth    UserAuthenticator
}

func NewAuthHandler(authService auth.AuthService, userAuth UserAuthenticator) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userAuth:    userAuth,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userId, err := h.userAuth.AuthenticateUser(req.Email, req.Password)
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
