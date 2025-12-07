package bootstrap

import (
	"backend/internal/infrastructure/config"
	"backend/internal/infrastructure/postgres"
	storagedb "backend/internal/storage/db"
	"backend/migrations"
	"context"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

var PostgresModule = fx.Module(
	"postgres",
	fx.Provide(postgres.NewPool),
	fx.Provide(func(pool *pgxpool.Pool) *storagedb.Queries {
		return storagedb.New(pool)
	}),
	fx.Invoke(RunStorageMigrations),
)

func RunStorageMigrations(
	lc fx.Lifecycle,
	cfg *config.AppConfig,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			source, err := iofs.New(migrations.MigrationsFS, "storage")
			if err != nil {
				return err
			}

			m, err := migrate.NewWithSourceInstance(
				"iofs",
				source,
				cfg.Postgres.Dsn,
			)
			if err != nil {
				return err
			}
			defer m.Close()

			if err := m.Up(); err != nil && err != migrate.ErrNoChange {
				return err
			}

			return nil
		},
	})
}
