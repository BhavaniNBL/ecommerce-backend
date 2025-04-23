package handler

import (
	"net/http"

	//"github.com/BhavaniNBL/ecommerce-backend/services/user-service/db"
	"github.com/BhavaniNBL/ecommerce-backend/config/db"
	// "github.com/BhavaniNBL/ecommerce-backend/services/user-service/middleware"
	"github.com/BhavaniNBL/ecommerce-backend/services/user-service/model"
	"github.com/BhavaniNBL/ecommerce-backend/services/user-service/service"
	"github.com/BhavaniNBL/ecommerce-backend/shared/middleware"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func SetupRoutes(r *gin.Engine) {
	// Public routes
	r.POST("/signup", SignUp)
	r.POST("/login", Login)
	r.GET("/user/:id", middleware.JWTMiddleware(), GetUser)
	r.GET("/users", middleware.JWTMiddleware(), ListUsers)

}

func SignUp(c *gin.Context) {
	var req model.SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, refreshToken, err := service.SignUp(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":            user.ID,
		"name":          user.Name,
		"email":         user.Email,
		"token":         token,
		"refresh_token": refreshToken,
	})
}

func Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user model.User
	if err := db.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Check if the password matches the hashed password in the database
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password"})
		return
	}

	token, refreshToken, err := service.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":         token,
		"refresh_token": refreshToken,
	})
}

func GetUser(c *gin.Context) {
	id := c.Param("id")

	user, err := service.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func ListUsers(c *gin.Context) {
	users, err := service.ListUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}
