package models

import (
	"time"

	"gorm.io/gorm"
)

type AbstractCreateUpdateModel struct {
	ID        uint      `json:"id"         bson:"_id"        gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

func (a *AbstractCreateUpdateModel) BeforeCreate(tx *gorm.DB) (err error) {
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()
	return
}

func (a *AbstractCreateUpdateModel) BeforeUpdate(tx *gorm.DB) (err error) {
	a.UpdatedAt = time.Now()
	return
}
