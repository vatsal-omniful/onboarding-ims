package seller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vatsal-omniful/onboarding-ims/models"
)

type SellerController struct {
	repo *SellerRepository
}

func (ctrl *SellerController) CreateSeller(ctx *gin.Context) {
	var seller models.Seller
	if err := ctx.ShouldBindJSON(&seller); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid seller data"})
		return
	}

	if err := ctrl.repo.CreateSeller(&seller); err != nil {
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

	seller, err := ctrl.repo.GetSellerById(sellerId)
	if err != nil || seller == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Seller not found"})
		return
	}

	ctx.JSON(http.StatusOK, seller)
}
