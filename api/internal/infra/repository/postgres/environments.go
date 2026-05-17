package postgres

import (
	"context"
	"log/slog"

	"github.com/adamkirk/bifrost/api/internal/common"
	"github.com/adamkirk/bifrost/api/internal/infra/repository/postgres/db"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EnvironmentsRepository struct {
	pool *pgxpool.Pool
	l    *slog.Logger
}

func (r *EnvironmentsRepository) Create(env *common.Environment) error {
	conn := db.New(r.pool)

	err := conn.InsertEnvironment(context.Background(), db.InsertEnvironmentParams{
		Name: env.Name,
		ID: pgtype.UUID{
			Bytes: [16]byte(env.ID[:]),
			Valid: true,
		},
	})

	if err != nil {
		r.l.Error("failed to insert environment", "error", err)
	}

	return err
}

func NewEnvironmentsRepository(l *slog.Logger, pool *pgxpool.Pool) *EnvironmentsRepository {
	return &EnvironmentsRepository{
		pool: pool,
		l:    l.With("component", "infra.postgres.environments_repository"),
	}
}
