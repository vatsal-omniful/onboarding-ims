package models

type Seller struct {
	AbstractCreateUpdateModel
	Name string `json:"name" bson:"name" bind:"required"`
}
