package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/raihankhan/jwt-auth-system/config" // Import config package
)

// AuthMiddleware is a Gin middleware to authenticate requests using JWT.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		config := c.MustGet("config").(*config.AppConfig) // Retrieve config from context

		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1) // Remove "Bearer " prefix

		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			return
		}

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			// Return secret key for signing validation
			return []byte(config.JWT.SecretKey), nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Validate token claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Check token expiration
			if float64(time.Now().Unix()) > claims["expiresAt"].(float64) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
				return
			}

			// Optionally, you can validate other claims like "issuer", "audience", etc. here

			// Extract User ID from claims and set it in context (for handler access)
			userID := uint(claims["userID"].(float64)) // Claims values are often float64
			c.Set("userID", userID)                    // Set userID in Gin context

			c.Next() // Proceed to the next handler (protected endpoint)

		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		}
	}
}
