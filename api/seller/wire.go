//go:build wireinject
// +build wireinject

package seller

import (
	"context"

	"github.com/google/wire"
	"github.com/omniful/go_commons/db/sql/postgres"
)

func Wire(ctx context.Context, db *postgres.DbCluster) (*SellerController, error) {
	panic(wire.Build(ProviderSet))
}
