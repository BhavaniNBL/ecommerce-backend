package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// var jwtSecret = []byte("secret_key") // Replace with a secure env var or Vault
var jwtSecret = []byte(os.Getenv("JWT_SECRET_KEY"))

// type Claims struct {
// 	UserID   string `json:"user_id"`
// 	UserType string `json:"user_type"` // e.g., "admin", "user"
// 	jwt.RegisteredClaims
// }

type Claims struct {
	Sub  string `json:"sub"`  // User ID
	Role string `json:"role"` // e.g. admin/user
	jwt.RegisteredClaims
}

// JwtMiddleware checks token validity and extracts user info
func JwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Missing Authorization header"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token format"})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token claims"})
			c.Abort()
			return
		}

		// Attach claims to context for access in handlers
		c.Set("user_id", claims.Sub)
		c.Set("user_type", claims.Role)

		c.Next()
	}
}

// AdminOnly ensures the user is an admin
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("user_type")
		if !exists || userType != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"message": "Admins only"})
			c.Abort()
			return
		}
		c.Next()
	}
}
