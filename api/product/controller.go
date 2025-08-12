package product

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vatsal-omniful/onboarding-ims/models"
)

type ProductController struct {
	ser  *ProductService
	repo *ProductRepository
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

	if err := ctrl.repo.CreateProduct(&product); err != nil {
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

	product, err := ctrl.repo.GetProductBySkuID(skuId)
	if err != nil || product == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "SKU ID not found"})
		return
	}
	ctx.JSON(http.StatusOK, product)
}
