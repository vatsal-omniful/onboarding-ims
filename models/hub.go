package models

import "gorm.io/datatypes"

type Hub struct {
	AbstractCreateUpdateModel
	Name        string         `json:"name"        bson:"name"        bind:"required"`
	Description string         `json:"description" bson:"description"`
	Location    datatypes.JSON `json:"location"    bson:"location"    bind:"required"`
	Status      string         `json:"status"      bson:"status"      bind:"required"`
	TenantId    uint           `json:"tenantId"    bson:"tenantId"    bind:"required"`
}
