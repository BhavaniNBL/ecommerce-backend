// package middleware

// import (
// 	"log"
// 	"net/http"
// 	"strings"

// 	"github.com/BhavaniNBL/ecommerce-backend/shared/util"
// 	"github.com/gin-gonic/gin"
// 	"github.com/golang-jwt/jwt/v5"
// )

// func JWTMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		authHeader := c.GetHeader("Authorization")
// 		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
// 			c.Abort()
// 			return
// 		}

// 		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
// 		token, err := util.ValidateAccessToken(tokenStr)
// 		if err != nil || !token.Valid {
// 			log.Printf("❌ Invalid Access Token: %v", err)
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
// 			c.Abort()
// 			return
// 		}

// 		claims := token.Claims.(jwt.MapClaims)
// 		c.Set("user_id", claims["user_id"])
// 		c.Next()
// 	}
// }

package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/BhavaniNBL/ecommerce-backend/shared/util"
	"github.com/gin-gonic/gin"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := util.ValidateAccessToken(tokenStr)
		if err != nil || !token.Valid {
			log.Printf("❌ Invalid Access Token: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		claims := token.Claims.(*util.Claims)
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}
