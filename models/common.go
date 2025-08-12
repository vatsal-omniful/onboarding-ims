package models

import (
	"time"
)

type AbstractCreateUpdateModel struct {
	ID        uint      `json:"id"         bson:"_id"        gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
