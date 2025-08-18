package seller

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/log"
	"github.com/vatsal-omniful/onboarding-ims/models"
)

type SellerController struct {
	repo *SellerRepository
}

var (
	ctrlOnce sync.Once
	ctrl     *SellerController
)

func NewSellerController(repo *SellerRepository) *SellerController {
	ctrlOnce.Do(func() {
		ctrl = &SellerController{repo: repo}
	})
	return ctrl
}

func (ctrl *SellerController) CreateSeller(ctx *gin.Context) {
	var seller models.Seller
	if err := ctx.ShouldBindJSON(&seller); err != nil {
		log.Errorf("Error binding seller data: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid seller data"})
		return
	}

	if seller.Name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Seller name is required"})
		return
	}

	if err := ctrl.repo.CreateSeller(ctx, &seller); err != nil {
		log.Errorf("Error creating seller: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create seller"})
		return
	}

	ctx.JSON(http.StatusCreated, seller)
}

func (ctrl *SellerController) GetSeller(ctx *gin.Context) {
	sellerId := ctx.Param("id")
	if sellerId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Seller ID is required"})
		return
	}

	seller, err := ctrl.repo.GetSellerById(ctx, sellerId)
	if err != nil || seller == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Seller not found"})
		return
	}

	ctx.JSON(http.StatusOK, seller)
}
