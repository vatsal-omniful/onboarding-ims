package routes

import (
	"context"

	"github.com/omniful/go_commons/http"
	"github.com/vatsal-omniful/onboarding-ims/api/product"
	"github.com/vatsal-omniful/onboarding-ims/pkg/db/postgres"
)

func InternalRoutes(ctx context.Context, server *http.Server) error {
	productController, err := product.Wire(ctx, postgres.GetCluster().DbCluster)
	if err != nil {
		return err
	}

	productGroup := server.Group("/product")
	{
		productGroup.PATCH("/fulfillOrder", productController.FulfillOrderRequest)
	}
	return nil
}
