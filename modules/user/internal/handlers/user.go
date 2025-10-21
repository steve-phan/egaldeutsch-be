package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"egaldeutsch-be/internal/auth"
	"egaldeutsch-be/internal/config"
	"egaldeutsch-be/modules/user/internal/models"
	"egaldeutsch-be/modules/user/internal/services"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
	userService *services.UserService
	jwtCfg      config.JwtConfig
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *services.UserService, jwtCfg config.JwtConfig) *UserHandler {
	return &UserHandler{userService: userService, jwtCfg: jwtCfg}
}

// CreateUser handles POST /api/v1/users
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.CreateUser(req.Name, req.Role)
	if err != nil {
		logrus.WithError(err).Error("Failed to create user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate JWT token
	token, err := auth.CreateAccessToken(user.ID.String(), h.jwtCfg)
	if err != nil {
		logrus.WithError(err).Error("Failed to generate access token")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user, "token": token})
}

// GetUser handles GET /api/v1/users/:id
func (h *UserHandler) GetUser(c *gin.Context) {
	var params models.UserIDParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.GetUserByID(params.ID)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		logrus.WithError(err).Error("Failed to get user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser handles PUT /api/v1/users/:id
func (h *UserHandler) UpdateUser(c *gin.Context) {
	var params models.UserIDParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.UpdateUser(params.ID, &req)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		logrus.WithError(err).Error("Failed to update user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser handles DELETE /api/v1/users/:id
func (h *UserHandler) DeleteUser(c *gin.Context) {
	var params models.UserIDParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userService.DeleteUser(params.ID); err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		logrus.WithError(err).Error("Failed to delete user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// ListUsers handles GET /api/v1/users
func (h *UserHandler) ListUsers(c *gin.Context) {
	page := 1
	perPage := 10

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if pp := c.Query("per_page"); pp != "" {
		if parsed, err := strconv.Atoi(pp); err == nil && parsed > 0 && parsed <= 100 {
			perPage = parsed
		}
	}

	users, total, err := h.userService.ListUsers(page, perPage)
	if err != nil {
		logrus.WithError(err).Error("Failed to list users")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	totalPages := (int(total) + perPage - 1) / perPage

	c.JSON(http.StatusOK, gin.H{
		"items":       users,
		"page":        page,
		"per_page":    perPage,
		"total_items": total,
		"total_pages": totalPages,
	})
}
