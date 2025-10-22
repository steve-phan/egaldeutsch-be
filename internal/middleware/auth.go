package middleware

import (
	"net/http"
	"strings"

	"egaldeutsch-be/internal/auth"
	"egaldeutsch-be/internal/config"
	"egaldeutsch-be/pkg/models"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates the Authorization: Bearer <token> header using the JWT config
// and stores the user id in the request context under the key "user_id".
func AuthMiddleware(jwtCfg config.JwtConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) < 8 || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid authorization header"})
			return
		}
		token := authHeader[7:]
		claims, err := auth.ParseToken(token, jwtCfg)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		c.Set("user_id", claims.UserId)
		// also set role from token claims (if present)
		if claims.Role != "" {
			c.Set("user_role", models.UserRole(claims.Role))
		}
		c.Next()
	}
}

// GetUserIDFromContext extracts the user id (string) from gin context, if present.
func GetUserIDFromContext(c *gin.Context) (string, bool) {
	v, ok := c.Get("user_id")
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}
