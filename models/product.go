package models

type Product struct {
	AbstractCreateUpdateModel
	Name        string  `json:"name"        bson:"name"        bind:"required"`
	Description string  `json:"description" bson:"description"`
	Price       float64 `json:"price"       bson:"price"       bind:"required"`
	SkuId       string  `json:"sku_id"      bson:"sku_id"      bind:"required" gorm:"unique"`
	TenantId    uint    `json:"tenant_id"   bson:"tenant_id"   bind:"required"`
}

type ProductHub struct {
	ProductId uint `json:"product_id" bson:"product_id" bind:"required" gorm:"primaryKey,uniqueIndex:idx_product_hub"`
	HubId     uint `json:"hub_id"     bson:"hub_id"     bind:"required" gorm:"primaryKey,uniqueIndex:idx_product_hub"`
	SellerId  uint `json:"seller_id"  bson:"seller_id"  bind:"required" gorm:"primaryKey"`
	Quantity  uint `json:"quantity"   bson:"quantity"   bind:"required"`

	Product Product `gorm:"foreignKey:ProductId;references:ID"`
	Hub     Hub     `gorm:"foreignKey:HubId;references:ID"`
	Seller  Seller  `gorm:"foreignKey:SellerId;references:ID"`
}
