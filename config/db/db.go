// // package db

// // import (
// // 	"fmt"
// // 	"log"
// // 	"os"

// // 	"gorm.io/driver/postgres"
// // 	"gorm.io/gorm"
// // )

// // var DB *gorm.DB

// // func InitDB() {
// // 	dsn := os.Getenv("POSTGRES_DSN")
// // 	var err error
// // 	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
// // 	if err != nil {
// // 		log.Fatalf("failed to connect to database: %v", err)
// // 	}
// // 	fmt.Println("🚀 Database connected")
// // }

// package db

// import (
// 	"fmt"
// 	"log"
// 	"os"

// 	"gorm.io/driver/postgres"
// 	"gorm.io/gorm"
// )

// var DB *gorm.DB

// func InitDB() {
// 	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
// 		os.Getenv("DB_HOST"),
// 		os.Getenv("DB_USER"),
// 		os.Getenv("DB_PASSWORD"),
// 		os.Getenv("DB_NAME"),
// 		os.Getenv("DB_PORT"),
// 	)

// 	var err error
// 	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		log.Fatalf("Error connecting to the database: %v", err)
// 	}

// 	log.Println("Database connection established")

// 	// Auto-migrate the User model to keep schema updated
// 	err = DB.AutoMigrate(&User{})
// 	if err != nil {
// 		log.Fatalf("Error during auto migration: %v", err)
// 	}
// }

package db

import (
	"fmt"
	"log"
	"os"

	"github.com/BhavaniNBL/ecommerce-backend/services/user-service/model" // ✅ adjust if path differs depending on your module name

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Error connecting to the database: %v", err)
	}

	log.Println("✅ Database connection established")

	// ✅ Auto-migrate the User model
	err = DB.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatalf("❌ Error during auto migration: %v", err)
	}

}
