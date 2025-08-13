package main

import (
	"log"

	"github.com/vatsal-omniful/onboarding-ims/database"
	"github.com/vatsal-omniful/onboarding-ims/routes"
)

func main() {
	// Connect to the database
	database.Connect()

	// Perform database migration
	database.Migrate()

	// Setup and run the Gin router
	router := routes.SetupRouter()
	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("couldn't run server")
	}
}
