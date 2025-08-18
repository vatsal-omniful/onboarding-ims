//go:build wireinject
// +build wireinject

package product

import (
	"context"

	"github.com/google/wire"
	"github.com/omniful/go_commons/db/sql/postgres"
)

func Wire(ctx context.Context, db *postgres.DbCluster) (*ProductController, error) {
	panic(wire.Build(ProviderSet))
}
