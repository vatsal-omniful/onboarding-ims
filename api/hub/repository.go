package hub

import (
	"github.com/gin-gonic/gin"

	"github.com/omniful/go_commons/db/sql/postgres"
	"github.com/vatsal-omniful/onboarding-ims/models"
)

type HubRepository struct {
	db *postgres.DbCluster
}

func NewHubRepository(db *postgres.DbCluster) *HubRepository {
	return &HubRepository{
		db: db,
	}
}

func (rep *HubRepository) CreateHubRepo(ctx *gin.Context, hub *models.Hub) error {
	return rep.db.GetMasterDB(ctx).Create(hub).Error
}

func (rep *HubRepository) GetHubById(ctx *gin.Context, hubId string) (*models.Hub, error) {
	var hub models.Hub
	if err := rep.db.GetMasterDB(ctx).Where("id = ?", hubId).First(&hub).Error; err != nil {
		return nil, err
	}
	return &hub, nil
}
