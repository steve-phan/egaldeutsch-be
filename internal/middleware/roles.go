package middleware

import (
	"net/http"

	"egaldeutsch-be/pkg/models"

	"github.com/gin-gonic/gin"
)

// RequireRole returns a middleware that allows only users with the given role.
// It expects `user_id` to already be set in the context by AuthMiddleware and
// the user role to be present in the user view stored on the context under
// "user_role" (if you prefer resolving user role lazily, adapt accordingly).
func RequireRole(role models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		// The simplest approach: expect the role to be present on the context
		v, ok := c.Get("user_role")
		if !ok {
			//log
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		if r, ok := v.(models.UserRole); ok {
			if r != role {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
				return
			}
			c.Next()
			return
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
	}
}
