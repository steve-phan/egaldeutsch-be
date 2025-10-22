package handlers

import (
	"egaldeutsch-be/internal/auth"
	"egaldeutsch-be/internal/config"
	authModels "egaldeutsch-be/modules/auth/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService auth.AuthService
	userService UserService
	jwtCfg      config.JwtConfig
}

func NewAuthHandler(authService auth.AuthService, userService UserService, jwtCfg config.JwtConfig) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userService: userService,
		jwtCfg:      jwtCfg,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req authModels.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userId, role, err := h.userService.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	// include user's role in the access token so middleware can derive role
	token, err := auth.CreateAccessToken(userId, string(role), h.jwtCfg)
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

// RefreshToken rotates a refresh token and returns a new access and refresh token.
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req authModels.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ip := c.ClientIP()
	ua := c.Request.UserAgent()
	access, newRefresh, err := h.authService.RefreshTokens(req.RefreshToken, ip, ua)
	if err != nil {
		if err == auth.ErrRefreshTokenReuse {
			// revoke-all already performed by repo; return 401
			c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token reuse detected"})
			return
		}
		if err == auth.ErrInvalidRefreshToken {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to refresh tokens"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"access_token": access, "refresh_token": newRefresh})
}

// GetCurrentUser returns the user info for the currently authenticated user.
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	// Expect Authorization: Bearer <token>
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
		return
	}
	// Basic parsing
	var token string
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
		return
	}

	claims, err := h.authService.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	// Lookup user view
	userView, err := h.userService.GetUserViewByID(claims.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": userView})
}
