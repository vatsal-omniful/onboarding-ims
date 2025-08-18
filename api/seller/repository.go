package seller

import (
	"sync"

	"github.com/vatsal-omniful/onboarding-ims/database"
	"github.com/vatsal-omniful/onboarding-ims/models"
)

type SellerRepository struct{}

var (
	repoOnce sync.Once
	repo     *SellerRepository
)

func NewSellerRepository() *SellerRepository {
	repoOnce.Do(func() {
		repo = &SellerRepository{}
	})
	return repo
}

func (repo *SellerRepository) CreateSeller(seller *models.Seller) error {
	return database.DB.Create(seller).Error
}

func (repo *SellerRepository) GetSellerById(sellerId string) (*models.Seller, error) {
	var seller models.Seller
	if err := database.DB.Where("id = ?", sellerId).First(&seller).Error; err != nil {
		return nil, err
	}
	return &seller, nil
}
