package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vatsal-omniful/onboarding-ims/middleware"

	"github.com/vatsal-omniful/onboarding-ims/api/hub"
	"github.com/vatsal-omniful/onboarding-ims/api/product"
	"github.com/vatsal-omniful/onboarding-ims/api/seller"
)

func PublicRouter(router *gin.Engine) {
	publicRouter := router.Group("/public")

	hubController := hub.HubController{}
	hubRouter := publicRouter.Group("/hub")
	{
		hubRouter.POST("/create", middleware.TenantIDMiddleware(), hubController.CreateHub)
		hubRouter.GET("/:id", middleware.TenantIDMiddleware(), hubController.GetHub)
	}

	productController := product.ProductController{}
	productRouter := publicRouter.Group("/product")
	{
		productRouter.POST(
			"/create",
			middleware.TenantIDMiddleware(),
			productController.CreateProduct,
		)
		productRouter.GET("/:skuId", middleware.TenantIDMiddleware(), productController.GetProduct)
		productRouter.PATCH(
			"/inflow",
			middleware.TenantIDMiddleware(),
			productController.InflowProduct,
		)
		productRouter.GET("/getAll", productController.GetProducts)
		productRouter.GET(
			"/getProductsByFilters",
			middleware.TenantIDMiddleware(),
			productController.GetInventory,
		)
	}

	sellerController := seller.SellerController{}
	sellerRouter := publicRouter.Group("/seller")
	{
		sellerRouter.POST("/create", middleware.TenantIDMiddleware(), sellerController.CreateSeller)
		sellerRouter.GET("/:id", middleware.TenantIDMiddleware(), sellerController.GetSeller)
	}
}
