package product

import (
	"errors"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/db/sql/postgres"
	"github.com/vatsal-omniful/onboarding-ims/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductRepository struct {
	db *postgres.DbCluster
}

var (
	repoOnce sync.Once
	repo     *ProductRepository
)

func NewProductRepository(db *postgres.DbCluster) *ProductRepository {
	repoOnce.Do(func() {
		repo = &ProductRepository{db: db}
	})
	return repo
}

func (repo *ProductRepository) CreateProduct(ctx *gin.Context, product *models.Product) error {
	return repo.db.GetMasterDB(ctx).Create(product).Error
}

func (repo *ProductRepository) GetProductBySkuID(
	ctx *gin.Context,
	skuId string,
) (*models.Product, error) {
	var product models.Product
	if err := repo.db.GetMasterDB(ctx).Where("sku_id = ?", skuId).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (repo *ProductRepository) GetProductById(
	ctx *gin.Context,
	productId uint,
) (*models.Product, error) {
	var product models.Product
	if err := repo.db.GetMasterDB(ctx).Where("id = ?", productId).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (repo *ProductRepository) GetHubById(ctx *gin.Context, hubId uint) (*models.Hub, error) {
	var hub models.Hub
	if err := repo.db.GetMasterDB(ctx).Where("id = ?", hubId).First(&hub).Error; err != nil {
		return nil, err
	}
	return &hub, nil
}

func (repo *ProductRepository) GetSellerById(
	ctx *gin.Context,
	sellerId uint,
) (*models.Seller, error) {
	var seller models.Seller
	if err := repo.db.GetMasterDB(ctx).Where("id = ?", sellerId).First(&seller).Error; err != nil {
		return nil, err
	}
	return &seller, nil
}

func (repo *ProductRepository) UpsertProductInflow(
	ctx *gin.Context,
	productHub *ProductInflow,
) (int, error) {
	var existingProductHub models.ProductHub

	tx := repo.db.GetMasterDB(ctx).Begin()
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
	ctx *gin.Context,
	productId uint,
	quantity uint,
) ([]uint, error) {
	var hubs []uint

	err := repo.db.GetMasterDB(ctx).Table("product_hubs").
		Joins("JOIN hubs ON product_hubs.hub_id = hubs.id").
		Select("hub_id").
		Where("product_id = ? AND quantity >= ? AND hubs.status = ?", productId, quantity, "active").
		Scan(&hubs).Error
	if err != nil {
		return nil, err
	}

	return hubs, nil
}

func (repo *ProductRepository) UpdateProductHubQuantity(
	ctx *gin.Context,
	productId, hubId, quantity uint,
) error {
	db := repo.db.GetMasterDB(ctx)

	result := db.Model(&models.ProductHub{}).
		Where("product_id = ? AND hub_id = ? AND quantity >= ?", productId, hubId, quantity).
		Where("hub_id IN (SELECT id FROM hubs WHERE id = ? AND status = 'active')", hubId).
		UpdateColumn("quantity", gorm.Expr("quantity - ?", quantity))

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows updated: insufficient quantity or inactive hub")
	}
	return nil
}

func (repo *ProductRepository) GetAllProducts(ctx *gin.Context) ([]*models.Product, error) {
	var products []*models.Product
	if err := repo.db.GetMasterDB(ctx).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (repo *ProductRepository) GetProductsByFilters(
	ctx *gin.Context,
	filters map[string]interface{},
) ([]*models.Product, error) {
	var products []*models.Product
	query := repo.db.GetMasterDB(ctx).Model(&models.Product{})

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

func (repo *ProductRepository) GetInventory(
	ctx *gin.Context,
	tenantId uint,
) ([]map[string]any, error) {
	var inventory []map[string]any
	err := repo.db.GetMasterDB(ctx).Table("product_hubs").
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
