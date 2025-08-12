package models

type Product struct {
	AbstractCreateUpdateModel
	Name        string  `json:"name"        bson:"name"        bind:"required"`
	Description string  `json:"description" bson:"description"`
	Price       float64 `json:"price"       bson:"price"       bind:"required"`
	SkuId       string  `json:"sku_id"      bson:"sku_id"      bind:"required" gorm:"unique"`
	TenantId    string  `json:"tenant_id"   bson:"tenant_id"   bind:"required"`
}

type ProductHub struct {
	ProductID uint `json:"product_id" bson:"product_id" bind:"required" gorm:"primaryKey,uniqueIndex:idx_product_hub"`
	HubID     uint `json:"hub_id"     bson:"hub_id"     bind:"required" gorm:"primaryKey,uniqueIndex:idx_product_hub"`
	SellerID  uint `json:"seller_id"  bson:"seller_id"  bind:"required" gorm:"primaryKey"`
	Quantity  uint `json:"quantity"   bson:"quantity"   bind:"required"`

	Product Product `gorm:"foreignKey:ProductID;references:ID"`
	Hub     Hub     `gorm:"foreignKey:HubID;references:ID"`
	Seller  Seller  `gorm:"foreignKey:SellerID;references:ID"`
}
