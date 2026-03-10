package auth

import (
	"net/http"

	"github.com/wb-go/wbf/ginext"
)

// RequireRoles - middleware для проверки прав доступа по ролям
func RequireRoles(roles ...string) ginext.HandlerFunc {

	return func(c *ginext.Context) {

		roleValue, exists := c.Get("role")

		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, ginext.H{"error": "role missing"})
			return
		}

		role := roleValue.(string)

		for _, r := range roles {
			if r == role {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, ginext.H{"error": "access denied"})
	}
}
