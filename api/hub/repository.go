package hub

import (
	"github.com/vatsal-omniful/onboarding-ims/database"
	"github.com/vatsal-omniful/onboarding-ims/models"
)

type HubRepository struct{}

func (rep *HubRepository) CreateHub(hub *models.Hub) error {
	return database.DB.Create(hub).Error
}

func (rep *HubRepository) GetHubById(hubId string) (*models.Hub, error) {
	var hub models.Hub
	if err := database.DB.Where("id = ?", hubId).First(&hub).Error; err != nil {
		return nil, err
	}
	return &hub, nil
}
