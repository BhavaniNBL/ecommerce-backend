// package middleware

// import (
// 	"net/http"
// 	"strings"

// 	"github.com/dgrijalva/jwt-go"
// 	"github.com/gin-gonic/gin"
// 	"github.com/golang-jwt/jwt/v4"
// )

// func Authenticate() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// Get token from header
// 		tokenString := c.GetHeader("Authorization")
// 		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

// 		if tokenString == "" {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is missing"})
// 			c.Abort()
// 			return
// 		}

// 		// Parse and validate the token
// 		claims := jwt.MapClaims{}
// 		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
// 			return []byte("your_jwt_secret_key"), nil
// 		})

// 		if err != nil || !token.Valid {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
// 			c.Abort()
// 			return
// 		}

// 		// Attach user info to context
// 		c.Set("user_id", claims["sub"])
// 		c.Set("role", claims["role"])

// 		c.Next()
// 	}
// }

package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// var jwtSecret = []byte("your_jwt_secret_key") // Replace with env or Vault in production
var jwtSecret = []byte(os.Getenv("JWT_SECRET_KEY"))

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization must be Bearer token"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.NewValidationError("unexpected signing method", jwt.ValidationErrorSignatureInvalid)
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("user_id", claims["sub"])
			c.Set("role", claims["role"])
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		c.Next()
	}
}
