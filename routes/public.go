package routes

import (
	"context"

	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/pagination"
	"github.com/vatsal-omniful/onboarding-ims/api/hub"
	"github.com/vatsal-omniful/onboarding-ims/api/product"
	"github.com/vatsal-omniful/onboarding-ims/api/seller"
	"github.com/vatsal-omniful/onboarding-ims/middleware"
	"github.com/vatsal-omniful/onboarding-ims/pkg/db/postgres"
)

func Initialize(ctx context.Context, server *http.Server) error {
	server.Use(config.Middleware())
	server.Use(pagination.Middleware())

	hubController, err := hub.Wire(ctx, postgres.GetCluster().DbCluster)
	if err != nil {
		return err
	}

	productController, err := product.Wire(ctx, postgres.GetCluster().DbCluster)
	if err != nil {
		return err
	}

	loggingMiddlewareOptions := http.LoggingMiddlewareOptions{
		Format:      config.GetString(ctx, "log.format"),
		Level:       config.GetString(ctx, "log.level"),
		LogRequest:  config.GetBool(ctx, "log.request"),
		LogResponse: config.GetBool(ctx, "log.response"),
	}

	publicRoute := server.Group(
		"/public",
		http.RequestLogMiddleware(loggingMiddlewareOptions),
	)

	hubRoute := publicRoute.Group("/hub", middleware.TenantIDMiddleware())
	{
		hubRoute.POST("/create", hubController.CreateHub)
		hubRoute.GET("/", hubController.GetHub)
	}

	productRoute := server.Group("/product")
	{
		productRoute.POST(
			"/create",
			middleware.TenantIDMiddleware(),
			productController.CreateProduct,
		)
		productRoute.GET("/:skuId", middleware.TenantIDMiddleware(), productController.GetProduct)
		productRoute.PATCH(
			"/inflow",
			middleware.TenantIDMiddleware(),
			productController.InflowProduct,
		)
		productRoute.GET("/getAll", productController.GetProducts)
		productRoute.GET(
			"/getProductsByFilters",
			middleware.TenantIDMiddleware(),
			productController.GetInventory,
		)
	}

	sellerController, err := seller.Wire(ctx, postgres.GetCluster().DbCluster)
	if err != nil {
		return err
	}

	sellerRoute := server.Group("/seller", middleware.TenantIDMiddleware())
	{
		sellerRoute.POST("/create", sellerController.CreateSeller)
		sellerRoute.GET("/:id", sellerController.GetSeller)
	}

	return nil
}
