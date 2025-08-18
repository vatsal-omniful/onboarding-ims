package product

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewProductController,
	NewProductService,
	NewProductRepository,
)
