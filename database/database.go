package database

import (
	"fmt"
	"golang-final-project/model"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func StartdB() {
	// get database variable from environment
	var err error
	err = godotenv.Load()
	if err != nil {
		fmt.Println("cannot get enc file :")
	}
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	// connect to mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", username, password, host, port, dbName)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		return
	}
	return
}

func AutoMigrate() {
	// Auto-migrate the User model
	if err := db.AutoMigrate(&model.User{}); err != nil {
		fmt.Println("Failed to migrate tables user:", err)
		return
	}
	// Auto-migrate the Photo model
	if err := db.AutoMigrate(&model.Photo{}); err != nil {
		fmt.Println("Failed to migrate tables photo:", err)
		return
	}
	// Auto-migrate the SocialMedia model
	if err := db.AutoMigrate(&model.SocialMedia{}); err != nil {
		fmt.Println("Failed to migrate tables social media:", err)
		return
	}
	// Auto-migrate the Comment model
	if err := db.AutoMigrate(&model.Comment{}); err != nil {
		fmt.Println("Failed to migrate tables comment:", err)
		return
	}
	// Auto-migrate the Nutrition Data model
	if err := db.AutoMigrate(&model.NutritionData{}); err != nil {
		fmt.Println("Failed to migrate tables comment:", err)
		return
	}
	// Auto-migrate the User Nutrition model
	if err := db.AutoMigrate(&model.UserNutrition{}); err != nil {
		fmt.Println("Failed to migrate tables comment:", err)
		return
	}
	return
}

func GetDB() *gorm.DB {
	return db
}
