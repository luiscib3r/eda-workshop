package postgres

import (
	"backend/internal/infrastructure/config"
	"context"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

func NewPool(
	lc fx.Lifecycle,
	cfg *config.AppConfig,
) (*pgxpool.Pool, error) {
	ctx := context.Background()

	pgcfg, err := pgxpool.ParseConfig(cfg.Postgres.Uri)
	if err != nil {
		return nil, err
	}

	pgcfg.ConnConfig.Tracer = otelpgx.NewTracer()

	pool, err := pgxpool.NewWithConfig(ctx, pgcfg)
	if err != nil {
		return nil, err
	}

	if err := otelpgx.RecordStats(pool); err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			pool.Close()
			return nil
		},
	})

	return pool, nil
}
