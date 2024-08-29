package config

import (
	"fmt"
	"golang-ms/src/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		LoadConfig().Postgres.Host, LoadConfig().Postgres.Port,
		LoadConfig().Postgres.User, LoadConfig().Postgres.Password,
		LoadConfig().Postgres.Database, LoadConfig().Postgres.SSLMode)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	DB.AutoMigrate(&models.User{})

}
