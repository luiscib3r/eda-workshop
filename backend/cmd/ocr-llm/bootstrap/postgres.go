package bootstrap

import (
	"backend/internal/infrastructure/postgres"
	orcdb "backend/internal/ocr/db"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

var PostgresModule = fx.Module(
	"postgres",
	fx.Provide(postgres.NewPool),
	fx.Provide(func(pool *pgxpool.Pool) *orcdb.Queries {
		return orcdb.New(pool)
	}),
)
