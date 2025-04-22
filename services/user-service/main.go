package main

import (
	"log"

	"github.com/BhavaniNBL/ecommerce-backend/config"
	"github.com/BhavaniNBL/ecommerce-backend/config/db"
	"github.com/BhavaniNBL/ecommerce-backend/services/user-service/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load environment variables
	config.LoadConfig()

	// Initialize the database
	db.InitDB()

	// Create a new Gin router
	r := gin.Default()

	// Set up routes for the User service
	handler.SetupRoutes(r)

	// Start the server on port 8080
	log.Println("User service is running on port 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Error starting the server: %v", err)
	}
}
