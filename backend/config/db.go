package config

import (
	"fmt"
	"log"

	"github.com/aeswibon/helmdeploy/backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is the global variable to access the database
var DB *gorm.DB

// ConnectDB connects to the database
func ConnectDB() {
	LoadEnv()
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", GetEnv("DB_HOST"), GetEnv("DB_USER"), GetEnv("DB_PASSWORD"), GetEnv("DB_NAME"), GetEnv("DB_PORT"))
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}
	log.Println("Connected to database")
	DB.AutoMigrate(&models.User{}, &models.Application{}, &models.Log{})
}
