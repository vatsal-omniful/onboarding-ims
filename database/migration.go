package database

import (
	"fmt"
	"log"

	"github.com/vatsal-omniful/onboarding-ims/models"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Migrate() {
	if err := DB.AutoMigrate(&models.Hub{}, &models.Product{}, &models.ProductHub{}, &models.Seller{}); err != nil {
		log.Fatalf("unable to migrate: %v", err)
	}
	fmt.Println("Database migration completed.")
}
