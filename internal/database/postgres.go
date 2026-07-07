package database

import (
	"fmt"
	"log"

	"github.com/testgithubanjali/url-shortener/internal/config"
	"github.com/testgithubanjali/url-shortener/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.AppConfig.DBHost,
		config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
		config.AppConfig.DBName,
		config.AppConfig.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Failed to connect to PostgreSQL:", err)
	}

	DB = db

	// Automatically create/update database tables
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal(" Failed to migrate database:", err)
	}

	log.Println("PostgreSQL connected successfully")
}
