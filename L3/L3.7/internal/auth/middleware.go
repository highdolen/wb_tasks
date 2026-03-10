package auth

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/wb-go/wbf/ginext"
)

// AuthMiddleware - middleware для проверки JWT токена
func AuthMiddleware(secret string) ginext.HandlerFunc {

	return func(c *ginext.Context) {

		header := c.GetHeader("Authorization")

		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ginext.H{"error": "missing token"})
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")

		token, err := jwt.ParseWithClaims(
			tokenStr,
			&Claims{},
			func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			},
		)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ginext.H{"error": "invalid token"})
			return
		}

		claims := token.Claims.(*Claims)

		c.Set("role", claims.Role)
		c.Set("username", claims.Username)

		c.Next()
	}
}
