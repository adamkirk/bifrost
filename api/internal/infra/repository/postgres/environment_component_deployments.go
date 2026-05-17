package postgres

import (
	"context"
	"log/slog"

	"github.com/adamkirk/bifrost/api/internal/common"
	"github.com/adamkirk/bifrost/api/internal/infra/repository/postgres/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DeploymentsRepository struct {
	pool *pgxpool.Pool
	l    *slog.Logger
}

func (r *DeploymentsRepository) ByID(id uuid.UUID) (*common.Deployment, error) {
	conn := db.New(r.pool)

	row, err := conn.GetDeploymentByID(context.Background(), pgtype.UUID{
		Bytes: [16]byte(id[:]),
		Valid: true,
	})

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		r.l.Error("failed to get deployment", "error", err)
		return nil, err
	}

	return &common.Deployment{
		ID:                     row.ID.Bytes,
		CreatedAt:              row.CreatedAt.Time,
		Status:                 common.DeploymentStatus(row.Status),
		EnvironmentID:          row.EnvironmentID.Bytes,
		EnvironmentComponentID: row.EnvironmentComponentID.Bytes,
	}, nil
}

func (r *DeploymentsRepository) Save(d *common.Deployment) error {
	conn := db.New(r.pool)

	_, err := conn.UpsertDeployment(context.Background(), db.UpsertDeploymentParams{
		ID: pgtype.UUID{
			Bytes: [16]byte(d.ID[:]),
			Valid: true,
		},
		EnvironmentID: pgtype.UUID{
			Bytes: [16]byte(d.EnvironmentID[:]),
			Valid: true,
		},
		EnvironmentComponentID: pgtype.UUID{
			Bytes: [16]byte(d.EnvironmentComponentID[:]),
			Valid: true,
		},
		CreatedAt: pgtype.Timestamptz{
			Time:  d.CreatedAt,
			Valid: true,
		},
		Status: string(d.Status),
	})

	if err != nil {
		r.l.Error("failed to save deployment", "error", err)
	}

	return err
}

func NewDeploymentsRepository(l *slog.Logger, pool *pgxpool.Pool) *DeploymentsRepository {
	return &DeploymentsRepository{
		pool: pool,
		l:    l.With("component", "infra.postgres.deployments_repository"),
	}
}
