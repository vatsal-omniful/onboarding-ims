package database

import (
	"context"
	"fmt"
	"log"

	"github.com/vatsal-omniful/onboarding-ims/models"
	"github.com/vatsal-omniful/onboarding-ims/pkg/db/postgres"
)

func Migrate() {
	db := postgres.GetCluster()
	if db == nil {
		log.Fatalf("Database cluster not initialized")
	}

	ctx := context.Background()
	masterDB := db.GetMasterDB(ctx)

	if err := masterDB.AutoMigrate(&models.Hub{}, &models.Product{}, &models.ProductHub{}, &models.Seller{}); err != nil {
		log.Fatalf("unable to migrate: %v", err)
	}
}
