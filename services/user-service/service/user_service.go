// package service

// import (
// 	"fmt"
// 	"time"

// 	"github.com/BhavaniNBL/ecommerce-backend/services/user-service/model"
// 	"github.com/BhavaniNBL/ecommerce-backend/services/user-service/repository"
// 	"github.com/BhavaniNBL/ecommerce-backend/services/user-service/util"

// 	"github.com/dgrijalva/jwt-go"
// 	"golang.org/x/crypto/bcrypt"
// )

// var userRepo = repository.NewUserRepository()

// func SignUp(req model.SignUpRequest) (*model.User, string, string, error) {
// 	// Hash password
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return nil, "", "", err
// 	}

// 	user := &model.User{
// 		ID:       util.GenerateUUID(),
// 		Name:     req.Name,
// 		Email:    req.Email,
// 		Password: string(hashedPassword),
// 		UserType: "customer", // Default user type
// 	}

// 	// Save user to DB
// 	err = userRepo.CreateUser(user)
// 	if err != nil {
// 		return nil, "", "", err
// 	}

// 	// Generate JWT and refresh tokens
// 	token, err := generateJWT(user)
// 	if err != nil {
// 		return nil, "", "", err
// 	}

// 	refreshToken, err := generateRefreshToken(user)
// 	if err != nil {
// 		return nil, "", "", err
// 	}

// 	return user, token, refreshToken, nil
// }

// func Login(req model.LoginRequest) (string, string, error) {
// 	// Find user by email
// 	user, err := userRepo.GetUserByEmail(req.Email)
// 	if err != nil {
// 		return "", "", err
// 	}

// 	// Compare password
// 	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
// 	if err != nil {
// 		return "", "", fmt.Errorf("invalid credentials")
// 	}

// 	// Generate JWT and refresh tokens
// 	token, err := generateJWT(user)
// 	if err != nil {
// 		return "", "", err
// 	}

// 	refreshToken, err := generateRefreshToken(user)
// 	if err != nil {
// 		return "", "", err
// 	}

// 	return token, refreshToken, nil
// }

// func GetUserByID(id string) (*model.User, error) {
// 	return userRepo.GetUserByID(id)
// }

// func ListUsers() ([]model.User, error) {
// 	return userRepo.ListUsers()
// }

// func generateJWT(user *model.User) (string, error) {
// 	claims := jwt.MapClaims{
// 		"sub":  user.ID,
// 		"role": user.UserType,
// 		"exp":  time.Now().Add(time.Hour * 24).Unix(), // Expiration time: 24 hours
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	return token.SignedString([]byte("your_jwt_secret_key"))
// }

// func generateRefreshToken(user *model.User) (string, error) {
// 	claims := jwt.MapClaims{
// 		"sub": user.ID,
// 		"exp": time.Now().Add(time.Hour * 72).Unix(), // Expiration time: 72 hours
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	return token.SignedString([]byte("your_jwt_secret_key"))
// }

package service

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/BhavaniNBL/ecommerce-backend/services/user-service/model"
	"github.com/BhavaniNBL/ecommerce-backend/services/user-service/repository"
	"github.com/BhavaniNBL/ecommerce-backend/services/user-service/util"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var userRepo = repository.NewUserRepository()

func SignUp(req model.SignUpRequest) (*model.User, string, string, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to hash password: %v", err)
	}

	// Validate user type (optional, in case you want only "admin" or "customer")
	if req.UserType != "admin" && req.UserType != "customer" {
		req.UserType = "customer" // default to customer if invalid
	}

	user := &model.User{
		ID:       util.GenerateUUID(),
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		//UserType: "customer", // Default user type
		UserType: req.UserType,
	}

	// Save user to DB
	err = userRepo.CreateUser(user)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to create user: %v", err)
	}

	// Generate JWT and refresh tokens
	token, err := generateJWT(user)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate JWT: %v", err)
	}

	refreshToken, err := generateRefreshToken(user)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate refresh token: %v", err)
	}

	log.Printf("👉 Signing up user: %s", req.Email)
	log.Printf("👉 Attempting login for: %s", req.Email)
	log.Printf("Hashed password: %s", string(hashedPassword))

	return user, token, refreshToken, nil
}

func Login(req model.LoginRequest) (string, string, error) {
	// Find user by email
	user, err := userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return "", "", fmt.Errorf("user not found: %v", err)
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", "", fmt.Errorf("invalid credentials")
	}

	// Generate JWT and refresh tokens
	token, err := generateJWT(user)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate JWT: %v", err)
	}

	refreshToken, err := generateRefreshToken(user)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %v", err)
	}

	return token, refreshToken, nil
}

func GetUserByID(id string) (*model.User, error) {
	return userRepo.GetUserByID(id)
}

func ListUsers() ([]model.User, error) {
	return userRepo.ListUsers()
}

func generateJWT(user *model.User) (string, error) {
	// Get the JWT secret key from environment variable
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		secretKey = "default_jwt_secret_key" // Use a default key if env var is not set
	}

	claims := jwt.MapClaims{
		"sub":  user.ID,
		"role": user.UserType,
		"exp":  time.Now().Add(time.Hour * 24).Unix(), // Expiration time: 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func generateRefreshToken(user *model.User) (string, error) {
	// Get the JWT secret key from environment variable
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		secretKey = "default_jwt_secret_key" // Use a default key if env var is not set
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 72).Unix(), // Expiration time: 72 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}
