package product

import (
	"errors"
	"net/http"

	"github.com/vatsal-omniful/onboarding-ims/database"
	"github.com/vatsal-omniful/onboarding-ims/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (repo *ProductRepository) GetProductById(productId uint) (*models.Product, error) {
	var product models.Product
	if err := database.DB.Where("id = ?", productId).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (repo *ProductRepository) GetHubById(hubId uint) (*models.Hub, error) {
	var hub models.Hub
	if err := database.DB.Where("id = ?", hubId).First(&hub).Error; err != nil {
		return nil, err
	}
	return &hub, nil
}

func (repo *ProductRepository) GetSellerById(sellerId uint) (*models.Seller, error) {
	var seller models.Seller
	if err := database.DB.Where("id = ?", sellerId).First(&seller).Error; err != nil {
		return nil, err
	}
	return &seller, nil
}

func (repo *ProductRepository) UpsertProductInflow(productHub *ProductInflow) (int, error) {
	var existingProductHub models.ProductHub

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(
		"product_id = ? AND hub_id = ?", productHub.ProductId, productHub.HubId,
	).First(&existingProductHub).Error

	switch err {
	case nil:
		if productHub.SellerId != existingProductHub.SellerId {
			tx.Rollback()
			return http.StatusBadRequest, errors.New("seller ID mismatch")
		}

		if err := tx.Model(&existingProductHub).Update("quantity", gorm.Expr("quantity + ?", productHub.Quantity)).Error; err != nil {
			tx.Rollback()
			return http.StatusInternalServerError, err
		}
	case gorm.ErrRecordNotFound:
		newProductHub := models.ProductHub{
			ProductId: productHub.ProductId,
			HubId:     productHub.HubId,
			SellerId:  productHub.SellerId,
			Quantity:  productHub.Quantity,
		}
		if err := tx.Create(&newProductHub).Error; err != nil {
			tx.Rollback()
			return http.StatusInternalServerError, err
		}
	default:
		tx.Rollback()
		return http.StatusInternalServerError, err
	}

	tx.Commit()
	return http.StatusOK, nil
}

func (repo *ProductRepository) GetValidHubsForProduct(
	productId uint,
	quantity uint,
) ([]uint, error) {
	var hubs []uint

	err := database.DB.Table("product_hubs").
		Joins("JOIN hubs ON product_hubs.hub_id = hubs.id").
		Select("hub_id").
		Where("product_id = ? AND quantity >= ? AND hubs.status = ?", productId, quantity, "active").
		Scan(&hubs).Error
	if err != nil {
		return nil, err
	}

	return hubs, nil
}

func (repo *ProductRepository) UpdateProductHubQuantity(productId, hubId, quantity uint) error {
	return database.DB.Model(&models.ProductHub{}).
		Joins("JOIN hubs ON product_hubs.hub_id = hubs.id").
		Where("hubs.status = ?", "active").
		Where("product_id = ? AND hub_id = ? AND quantity >= ?", productId, hubId, quantity).
		UpdateColumn("quantity", gorm.Expr("quantity - ?", quantity)).Error
}

func (repo *ProductRepository) GetAllProducts() ([]*models.Product, error) {
	var products []*models.Product
	if err := database.DB.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (repo *ProductRepository) GetProductsByFilters(
	filters map[string]interface{},
) ([]*models.Product, error) {
	var products []*models.Product
	query := database.DB.Model(&models.Product{})

	for key, value := range filters {
		if key != "sku_id" {
			query = query.Where(key+" = ?", value)
		} else {
			query = query.Where("sku_id IN ?", value)
		}
	}

	if err := query.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (repo *ProductRepository) GetInventory(tenantId uint) ([]map[string]any, error) {
	var inventory []map[string]any
	err := database.DB.Table("product_hubs").
		Joins("JOIN products ON product_hubs.product_id = products.id").
		Joins("JOIN hubs ON product_hubs.hub_id = hubs.id").
		Joins("JOIN sellers ON product_hubs.seller_id = sellers.id").
		Select("products.sku_id as sku_id, hubs.id AS hub_id, hubs.name AS hub_name, product_hubs.quantity, sellers.name AS seller_name").
		Where("products.tenant_id = ? and hubs.tenant_id = ?", tenantId, tenantId).
		Scan(&inventory).Error
	if err != nil {
		return nil, err
	}
	return inventory, nil
}
