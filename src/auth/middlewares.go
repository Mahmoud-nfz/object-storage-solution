package auth

import (
	"data-storage/src/config"
	"log"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

func EnsureUserAuthenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		// Extract the token from the header
		re := regexp.MustCompile(`^Bearer (.+)$`)
		matches := re.FindStringSubmatch(authHeader)
		if len(matches) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			return
		}
		tokenString := matches[1]

		// Parse and validate the token
		claims, err := verify(tokenString)
		if err != nil {
			log.Println("Invalid or expired token:", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		// Set claims in context to use them
		c.Set("claims", claims)

		c.Next()
	}
}

func EnsureBackendAuthenticated() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		apiKey := ctx.GetHeader("X-API-Key")

		if apiKey != config.Env.APIKey {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		ctx.Next()
	}
}
