package middleware

import (
	"net/http"
	"strings"

	"github.com/aeswibon/shepherd/utils"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware is a middleware to authenticate requests
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BearerSchema = "Bearer "
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, BearerSchema) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenString := authHeader[len(BearerSchema):]
		_, err := utils.ValidateJWT(tokenString)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}
