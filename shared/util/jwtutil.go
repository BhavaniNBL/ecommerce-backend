// package util

// import (
// 	"errors"
// 	"os"
// 	"time"

// 	"github.com/golang-jwt/jwt/v5"
// )

// var (
// 	accessSecret  = []byte(os.Getenv("ACCESS_SECRET"))  // 15-min token
// 	refreshSecret = []byte(os.Getenv("REFRESH_SECRET")) // 7-day token
// )

// func GenerateAccessToken(userID string) (string, error) {
// 	claims := jwt.MapClaims{
// 		"user_id": userID,
// 		"exp":     time.Now().Add(15 * time.Minute).Unix(),
// 	}
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	return token.SignedString(accessSecret)
// }

// func GenerateRefreshToken(userID string) (string, error) {
// 	claims := jwt.MapClaims{
// 		"user_id": userID,
// 		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
// 	}
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	return token.SignedString(refreshSecret)
// }

// func ValidateAccessToken(tokenString string) (*jwt.Token, error) {
// 	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, errors.New("unexpected signing method")
// 		}
// 		return accessSecret, nil
// 	})
// }

// func ValidateRefreshToken(tokenString string) (*jwt.Token, error) {
// 	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, errors.New("unexpected signing method")
// 		}
// 		return refreshSecret, nil
// 	})
// }

package util

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	accessSecret  = []byte(os.Getenv("JWT_SECRET_KEY"))         // Short-lived token
	refreshSecret = []byte(os.Getenv("JWT_REFRESH_SECRET_KEY")) // Long-lived token
)

// Claims structure
type Claims struct {
	UserID string `json:"sub"`
	Role   string `json:"role,omitempty"`
	jwt.RegisteredClaims
}

// GenerateAccessToken creates a short-lived token (15m)
func GenerateAccessToken(userID, role string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(accessSecret)
}

// GenerateRefreshToken creates a long-lived token (7d)
func GenerateRefreshToken(userID string) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(refreshSecret)
}

// ValidateAccessToken parses and validates access token
func ValidateAccessToken(tokenStr string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return accessSecret, nil
	})
}

// ValidateRefreshToken parses and validates refresh token
func ValidateRefreshToken(tokenStr string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return refreshSecret, nil
	})
}
