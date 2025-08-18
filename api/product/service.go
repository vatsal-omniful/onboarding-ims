package product

import (
	"errors"
	"math/rand"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/vatsal-omniful/onboarding-ims/models"
)

type ProductService struct {
	repo *ProductRepository
}

var (
	serOnce sync.Once
	ser     *ProductService
)

func NewProductService(repo *ProductRepository) *ProductService {
	serOnce.Do(func() {
		ser = &ProductService{repo: repo}
	})
	return ser
}

func (service *ProductService) CheckValidityOfProductHub(
	ctx *gin.Context, productId, hubId, sellerId, tenantId uint,
) error {
	if productId == 0 || hubId == 0 || sellerId == 0 {
		return errors.New("product_id, hub_id, and seller_id must be greater than 0")
	}

	product, err := service.repo.GetProductById(ctx, productId)
	if err != nil {
		return errors.New("invalid product_id")
	}

	hub, err := service.repo.GetHubById(ctx, hubId)
	if err != nil {
		return errors.New("invalid hub_id")
	}

	if _, err := service.repo.GetSellerById(ctx, sellerId); err != nil {
		return errors.New("invalid seller_id")
	}

	if product.TenantId != hub.TenantId || product.TenantId != tenantId {
		return errors.New(
			"product and hub must belong to the same tenant or the product must belong to the specified tenant",
		)
	}

	return nil
}

func (service *ProductService) selectValidHubsForProduct(hubs []uint) uint {
	return hubs[rand.Intn(len(hubs))]
}

func (service *ProductService) FulfillOrderRequest(
	ctx *gin.Context,
	orderRequest *OrderRequest,
) (int, error) {
	if orderRequest.Quantity <= 0 {
		return http.StatusBadRequest, errors.New("quantity must be greater than 0")
	}

	product, err := service.repo.GetProductById(ctx, orderRequest.ProductId)
	if err != nil {
		return http.StatusBadRequest, errors.New("invalid product_id")
	}

	validHubs, err := service.repo.GetValidHubsForProduct(ctx, product.ID, orderRequest.Quantity)
	if len(validHubs) == 0 || err != nil {
		return http.StatusConflict, errors.New("no valid hubs available")
	}

	selectedHubId := service.selectValidHubsForProduct(validHubs)

	updateProductHubError := service.repo.UpdateProductHubQuantity(
		ctx, orderRequest.ProductId, selectedHubId, orderRequest.Quantity,
	)
	if updateProductHubError != nil {
		return http.StatusInternalServerError, errors.New("failed to fulfill order request")
	}
	return http.StatusOK, nil
}

func (service *ProductService) GetProducts(
	ctx *gin.Context,
	filters map[string]interface{},
) ([]*models.Product, error) {
	var products []*models.Product
	var err error

	if len(filters) == 0 {
		products, err = service.repo.GetAllProducts(ctx)
	} else {
		products, err = service.repo.GetProductsByFilters(ctx, filters)
	}

	if err != nil {
		return nil, err
	}
	return products, nil
}

func (service *ProductService) GetInventory(
	ctx *gin.Context,
	tenantId uint,
) ([]*map[string]interface{}, error) {
	inventory, err := service.repo.GetInventory(ctx, tenantId)
	if err != nil {
		return nil, err
	}

	var result []*map[string]any
	for _, item := range inventory {
		result = append(result, &map[string]any{
			"product_sku_id": item["sku_id"],
			"hub_id":         item["hub_id"],
			"quantity":       item["quantity"],
			"seller_name":    item["seller_name"],
		})
	}
	return result, nil
}
