package database

import (
	"fmt"
	"log"

	"github.com/vatsal-omniful/onboarding-ims/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := "host=localhost user=user password=password dbname=omsdb port=5432 sslmode=disable TimeZone=Asia/Kolkata"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	fmt.Println("Connected to database successfully!")
}

func Migrate() {
	if err := DB.AutoMigrate(&models.Hub{}, &models.Product{}, &models.ProductHub{}, &models.Seller{}); err != nil {
		log.Fatalf("unable to migrate: %v", err)
	}
	fmt.Println("Database migration completed.")
}
