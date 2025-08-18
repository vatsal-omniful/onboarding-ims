package main

import (
	"context"
	"time"

	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/shutdown"
	"github.com/vatsal-omniful/onboarding-ims/database"
	appinit "github.com/vatsal-omniful/onboarding-ims/init"
	"github.com/vatsal-omniful/onboarding-ims/routes"
)

func main() {
	err := config.Init(time.Second * 10)
	if err != nil {
		log.Panicf("Error while initializing config, err: %v", err)
		panic(err)
	}

	ctx, err := config.TODOContext()
	if err != nil {
		log.Panicf("Error while getting context from config, err: %v", err)
		panic(err)
	}

	appinit.Initialize(ctx)

	database.Migrate()

	runHttpServer(ctx)
}

func runHttpServer(ctx context.Context) {
	server := http.InitializeServer(
		config.GetString(ctx, "server.port"),
		10*time.Second,
		10*time.Second,
		70*time.Second,
		true,
	)

	// Initialize middlewares and routes
	err := routes.Initialize(ctx, server)
	if err != nil {
		log.Errorf(err.Error())
		panic(err)
	}

	err = routes.InternalRoutes(ctx, server)
	if err != nil {
		log.Errorf(err.Error())
		panic(err)
	}

	log.Infof("Starting server on port" + config.GetString(ctx, "server.port"))

	err = server.StartServer("Onboarding-ims")
	if err != nil {
		log.Errorf(err.Error())
		panic(err)
	}

	<-shutdown.GetWaitChannel()
}
