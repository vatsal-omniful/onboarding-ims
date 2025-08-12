package product

import (
	"github.com/vatsal-omniful/onboarding-ims/database"
	"github.com/vatsal-omniful/onboarding-ims/models"
)

type ProductRepository struct{}

func (repo *ProductRepository) CreateProduct(product *models.Product) error {
	return database.DB.Create(product).Error
}

func (repo *ProductRepository) GetProductBySkuID(skuId string) (*models.Product, error) {
	var product models.Product
	if err := database.DB.Where("sku_id = ?", skuId).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}
