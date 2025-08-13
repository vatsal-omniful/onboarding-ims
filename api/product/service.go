package product

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
)

type ProductService struct {
	repo *ProductRepository
}

func (service *ProductService) CheckValidityOfProductHub(
	productId, hubId, sellerId, tenantId uint,
) error {
	if productId == 0 || hubId == 0 || sellerId == 0 {
		return errors.New("product_id, hub_id, and seller_id must be greater than 0")
	}

	product, err := service.repo.GetProductById(productId)
	if err != nil {
		return errors.New("invalid product_id")
	}

	hub, err := service.repo.GetHubById(hubId)
	if err != nil {
		return errors.New("invalid hub_id")
	}

	if _, err := service.repo.GetSellerById(sellerId); err != nil {
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

func (service *ProductService) FulfillOrderRequest(orderRequest *OrderRequest) (int, error) {
	if orderRequest.Quantity <= 0 {
		return http.StatusBadRequest, errors.New("quantity must be greater than 0")
	}

	product, err := service.repo.GetProductById(orderRequest.ProductId)
	if err != nil {
		return http.StatusBadRequest, errors.New("invalid product_id")
	}

	validHubs, err := service.repo.GetValidHubsForProduct(product.ID, orderRequest.Quantity)
	if len(validHubs) == 0 || err != nil {
		return http.StatusConflict, errors.New("no valid hubs available")
	}

	selectedHubId := service.selectValidHubsForProduct(validHubs)

	updateProductHubError := service.repo.UpdateProductHubQuantity(
		orderRequest.ProductId, selectedHubId, orderRequest.Quantity,
	)
	if updateProductHubError != nil {
		return http.StatusInternalServerError, errors.New("failed to fulfill order request")
	}
	return http.StatusOK, nil
}
