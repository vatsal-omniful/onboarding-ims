package product

import (
	"errors"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/vatsal-omniful/onboarding-ims/models"
)

type ProductController struct {
	ser  *ProductService
	repo *ProductRepository
}

var (
	ctrlOnce sync.Once
	ctrl     *ProductController
)

func NewProductController(ser *ProductService, repo *ProductRepository) *ProductController {
	ctrlOnce.Do(func() {
		ctrl = &ProductController{
			ser:  ser,
			repo: repo,
		}
	})
	return ctrl
}

func (ctrl *ProductController) validateCreateProductRequest(product *models.Product) error {
	if product.Name == "" {
		return errors.New("product name is required")
	}

	if product.Price <= 0 {
		return errors.New("product price must be positive")
	}

	if product.SkuId == "" {
		return errors.New("product SKU ID is required")
	}

	if strings.Contains(product.SkuId, " ") {
		return errors.New("product SKU ID can not contain spaces")
	}

	return nil
}

func (ctrl *ProductController) CreateProduct(ctx *gin.Context) {
	var product models.Product
	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.validateCreateProductRequest(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.repo.CreateProduct(ctx, &product); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	ctx.JSON(http.StatusCreated, product)
}

func (ctrl *ProductController) GetProduct(ctx *gin.Context) {
	skuId := ctx.Param("skuId")

	if skuId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "SKU ID is needed"})
		return
	}

	product, err := ctrl.repo.GetProductBySkuID(ctx, skuId)
	if err != nil || product == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "SKU ID not found"})
		return
	}
	ctx.JSON(http.StatusOK, product)
}

type ProductInflow struct {
	ProductId uint `json:"product_id"         binding:"required"`
	HubId     uint `json:"hub_id"             binding:"required"`
	SellerId  uint `json:"seller_id"          binding:"required"`
	Quantity  uint `json:"quantity_available" binding:"required"`
}

func (ctrl *ProductController) InflowProduct(ctx *gin.Context) {
	var inflow ProductInflow
	if err := ctx.ShouldBindJSON(&inflow); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid inflow data"})
		return
	}

	if inflow.Quantity <= 0 {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Quantity available can never be 0 for an inflow."},
		)
		return
	}

	tenantId, _ := ctx.Get("tenantId")

	requestValidationErr := ctrl.ser.CheckValidityOfProductHub(
		ctx,
		inflow.ProductId,
		inflow.HubId,
		inflow.SellerId,
		tenantId.(uint),
	)

	if requestValidationErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": requestValidationErr.Error()})
		return
	}

	if statusCode, err := ctrl.repo.UpsertProductInflow(ctx, &inflow); statusCode != http.StatusOK {
		ctx.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Product inflow processed successfully"})
}

type OrderRequest struct {
	ProductId uint `json:"product_id" binding:"required"`
	Quantity  uint `json:"quantity"   binding:"required"`
}

func (ctrl *ProductController) FulfillOrderRequest(ctx *gin.Context) {
	var order OrderRequest
	if err := ctx.ShouldBindJSON(&order); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid order request"})
		return
	}

	if order.Quantity <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Quantity must be greater than 0"})
		return
	}

	statusCode, err := ctrl.ser.FulfillOrderRequest(ctx, &order)
	if err != nil {
		ctx.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(statusCode, gin.H{"message": "Order fulfilled successfully"})
}

func (ctrl *ProductController) GetProducts(ctx *gin.Context) {
	tenantId := ctx.Query("tenant_id")
	sellerId := ctx.Query("seller_id")
	sku_codes := ctx.QueryArray("sku_codes")

	filter := make(map[string]any)
	if tenantId != "" {
		filter["tenant_id"] = tenantId
	}
	if sellerId != "" {
		filter["seller_id"] = sellerId
	}
	if len(sku_codes) > 0 {
		filter["sku_id"] = sku_codes
	}
	products, err := ctrl.ser.GetProducts(ctx, filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products"})
		return
	}
	if len(products) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "No products found"})
		return
	}
	ctx.JSON(http.StatusOK, products)
}

func (ctrl *ProductController) GetInventory(ctx *gin.Context) {
	tenantId, _ := ctx.Get("tenantId")

	inventory, err := ctrl.ser.GetInventory(ctx, tenantId.(uint))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve inventory"})
		return
	}
	if len(inventory) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "No inventory found"})
		return
	}
	ctx.JSON(http.StatusOK, inventory)
}
