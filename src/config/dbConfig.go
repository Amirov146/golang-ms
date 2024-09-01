package config

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang-ms/src/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

var DB *gorm.DB
var client *mongo.Client
var TokenCollection *mongo.Collection

func ConnectPostgresDB() {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		LoadConfig().Postgres.Host, LoadConfig().Postgres.Port,
		LoadConfig().Postgres.User, LoadConfig().Postgres.Password,
		LoadConfig().Postgres.Database, LoadConfig().Postgres.SSLMode)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	DB.AutoMigrate(&models.Role{})
	DB.AutoMigrate(&models.UsersRoles{})
	DB.AutoMigrate(&models.User{})

	DB.Exec("INSERT INTO roles (id, name) VALUES (?, ?), (?, ?)", 1, "USER", 2, "ADMIN")
}

func ConnectMongoDB() {
	var login string
	if len(LoadConfig().MongoDB.User) > 0 && len(LoadConfig().MongoDB.Password) > 0 {
		login = fmt.Sprintf("%s:%s@", LoadConfig().MongoDB.User, LoadConfig().MongoDB.Password)
	}
	uri := fmt.Sprintf("mongodb://%s%s:%s/?maxPoolSize=%d&w=majority",
		login,
		LoadConfig().MongoDB.Host,
		LoadConfig().MongoDB.Port,
		LoadConfig().MongoDB.MaxPoolSize)

	var err error
	client, err = mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Failed to create MongoDB client:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB server:", err)
	}
	TokenCollection = client.Database(LoadConfig().MongoDB.Database).Collection("tokens_cache")
	fmt.Println("Successfully connected and pinged MongoDB server")
}

func DisconnectMongoDB() {
	if client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = client.Disconnect(ctx)
	}
}
