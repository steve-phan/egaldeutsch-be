package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token           string `json:"token" binding:"required"`
	Password        string `json:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" binding:"required,eqfield=Password"`
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Lookup user by email via unified user service
	userService := h.userService
	userID, err := userService.GetUserIDByEmail(req.Email)
	if err == nil {
		// create reset token
		token, err := h.authService.CreatePasswordResetForUser(userID)
		if err != nil {
			logrus.WithError(err).Error("failed to create password reset token")
			// still return accepted to avoid enumeration
			c.Status(http.StatusAccepted)
			return
		}

		// log the reset link (replace with Mailer in future)
		resetLink := "https://your.app/reset-password?token=" + token
		logrus.WithField("email", req.Email).Infof("password reset link: %s", resetLink)
	}

	// Always return accepted to avoid enumeration
	c.Status(http.StatusAccepted)
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := h.authService.VerifyPasswordResetToken(req.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		return
	}

	// Update password via unified user service
	if err := h.userService.UpdatePassword(userID, req.Password); err != nil {
		logrus.WithError(err).Error("failed to update password")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	// Revoke all refresh tokens for this user
	_ = h.authService.RevokeAllRefreshTokens(userID)

	c.Status(http.StatusOK)
}
