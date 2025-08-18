package init

import (
	"context"
	"time"

	"github.com/omniful/go_commons/config"
	opostgres "github.com/omniful/go_commons/db/sql/postgres"
	"github.com/omniful/go_commons/log"
	"github.com/vatsal-omniful/onboarding-ims/pkg/db/postgres"
)

func Initialize(ctx context.Context) {
	initializeLog(ctx)
	initializeDB(ctx)
}

// Initialize logging
func initializeLog(ctx context.Context) {
	err := log.InitializeLogger(
		log.Formatter(config.GetString(ctx, "log.format")),
		log.Level(config.GetString(ctx, "log.level")),
	)
	if err != nil {
		log.WithError(err).Panic("unable to initialise log")
	}
}

func initializeDB(ctx context.Context) {
	maxOpenConnections := config.GetInt(ctx, "postgresql.maxOpenConns")
	maxIdleConnections := config.GetInt(ctx, "postgresql.maxIdleConns")

	database := config.GetString(ctx, "postgresql.database")
	connIdleTimeout := 10 * time.Minute

	// Read Write endpoint config
	postgresWriteServer := config.GetString(ctx, "postgresql.master.host")
	postgresWritePort := config.GetString(ctx, "postgresql.master.port")
	postgresWritePassword := config.GetString(ctx, "postgresql.master.password")
	postgresWriterUsername := config.GetString(ctx, "postgresql.master.username")

	debugMode := config.GetBool(ctx, "postgresql.debugMode")

	// Master config i.e. - Write endpoint
	masterConfig := opostgres.DBConfig{
		Host:               postgresWriteServer,
		Port:               postgresWritePort,
		Username:           postgresWriterUsername,
		Password:           postgresWritePassword,
		Dbname:             database,
		MaxOpenConnections: maxOpenConnections,
		MaxIdleConnections: maxIdleConnections,
		ConnMaxLifetime:    connIdleTimeout,
		DebugMode:          debugMode,
	}

	db := opostgres.InitializeDBInstance(masterConfig, &[]opostgres.DBConfig{})
	log.InfofWithContext(ctx, "Initialized Postgres DB client")
	postgres.SetCluster(db)
}
