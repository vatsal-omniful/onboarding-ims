package hub

import (
	"errors"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/vatsal-omniful/onboarding-ims/models"
)

type HubController struct {
	rep *HubRepository
}

var (
	ctrlOnce sync.Once
	ctrl     *HubController
)

func NewHubController(rep *HubRepository) *HubController {
	ctrlOnce.Do(func() {
		ctrl = &HubController{
			rep: rep,
		}
	})
	return ctrl
}

func (c *HubController) validateHub(hub *models.Hub) error {
	if hub.TenantId == 0 {
		return errors.New("tenantId is required")
	}
	if hub.Name == "" {
		return errors.New("name is required")
	}
	if hub.Location == nil {
		return errors.New("location is required")
	}
	return nil
}

func (c *HubController) CreateHub(ctx *gin.Context) {
	var hub models.Hub
	if err := ctx.ShouldBindJSON(&hub); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.validateHub(&hub); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.rep.CreateHubRepo(ctx, &hub); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, hub)
}

type TenantId struct {
	TenantId string `json:"tenant_id" bson:"tenant_id" bind:"required"`
}

func (c *HubController) validateGetHubRequest(ctx *gin.Context) (string, uint, error) {
	hubId := ctx.Param("id")
	if hubId == "" {
		return "", 0, errors.New("hub ID is required")
	}
	tenantId, _ := ctx.Get("tenantId")

	return hubId, uint(tenantId.(float64)), nil
}

func (c *HubController) GetHub(ctx *gin.Context) {
	hubId, tenantId, err := c.validateGetHubRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hub, err := c.rep.GetHubById(ctx, hubId)
	if err != nil || hub == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Hub not found"})
		return
	}

	if hub.TenantId != tenantId {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Tenant ID mismatch"})
		return
	}

	ctx.JSON(http.StatusOK, hub)
}
