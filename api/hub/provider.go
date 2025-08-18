package hub

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewHubController,
	NewHubRepository,
)
