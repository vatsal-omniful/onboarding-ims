package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vatsal-omniful/onboarding-ims/api/product"
)

func InternalRouter(router *gin.Engine) {
	// Define internal routes here
	internalGroup := router.Group("/internal")
	productController := product.ProductController{}

	productGroup := internalGroup.Group("/product")
	{
		productGroup.PATCH("/fulfillOrder", productController.FulfillOrderRequest)
	}
}
