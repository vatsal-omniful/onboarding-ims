package seller

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewSellerController,
	NewSellerRepository,
)
