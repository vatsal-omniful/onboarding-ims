package seller

import (
	"context"
	"sync"

	"github.com/vatsal-omniful/onboarding-ims/models"
	"github.com/vatsal-omniful/onboarding-ims/pkg/db/postgres"
)

type SellerRepository struct {
	db *postgres.Db
}

var (
	repoOnce sync.Once
	repo     *SellerRepository
)

func NewSellerRepository() *SellerRepository {
	repoOnce.Do(func() {
		repo = &SellerRepository{
			db: postgres.GetCluster(),
		}
	})
	return repo
}

func (repo *SellerRepository) CreateSeller(ctx context.Context, seller *models.Seller) error {
	return repo.db.GetMasterDB(ctx).Create(seller).Error
}

func (repo *SellerRepository) GetSellerById(ctx context.Context, sellerId string) (*models.Seller, error) {
	var seller models.Seller
	if err := repo.db.GetMasterDB(ctx).Where("id = ?", sellerId).First(&seller).Error; err != nil {
		return nil, err
	}
	return &seller, nil
}
