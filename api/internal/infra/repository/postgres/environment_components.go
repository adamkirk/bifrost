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

type EnvironmentComponentsRepository struct {
	pool *pgxpool.Pool
	l    *slog.Logger
}

func (r *EnvironmentComponentsRepository) ByEnvironmentAndName(environmentID uuid.UUID, name string) (*common.EnvironmentComponent, error) {
	conn := db.New(r.pool)

	row, err := conn.GetEnvironmentComponentByEnvironmentAndName(context.Background(), db.GetEnvironmentComponentByEnvironmentAndNameParams{
		EnvironmentID: pgtype.UUID{
			Bytes: [16]byte(environmentID[:]),
			Valid: true,
		},
		Name: name,
	})

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		r.l.Error("failed to get environment component", "error", err)
		return nil, err
	}

	return &common.EnvironmentComponent{
		ID:            row.ID.Bytes,
		EnvironmentID: row.EnvironmentID.Bytes,
		Name:          row.Name,
		ChartName:     row.ChartName,
		ChartVersion:  row.ChartVersion,
		ChartRegistry: row.ChartRegistry,
	}, nil
}

func (r *EnvironmentComponentsRepository) CountByEnvironment(environmentID uuid.UUID) (int, error) {
	conn := db.New(r.pool)

	count, err := conn.CountEnvironmentComponents(context.Background(), pgtype.UUID{
		Bytes: [16]byte(environmentID[:]),
		Valid: true,
	})
	if err != nil {
		r.l.Error("failed to count environment components", "error", err)
		return 0, err
	}

	return int(count), nil
}

func (r *EnvironmentComponentsRepository) ListByEnvironment(environmentID uuid.UUID, limit, offset int) ([]*common.EnvironmentComponent, error) {
	conn := db.New(r.pool)

	rows, err := conn.ListEnvironmentComponents(context.Background(), db.ListEnvironmentComponentsParams{
		EnvironmentID: pgtype.UUID{
			Bytes: [16]byte(environmentID[:]),
			Valid: true,
		},
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		r.l.Error("failed to list environment components", "error", err)
		return nil, err
	}

	components := make([]*common.EnvironmentComponent, len(rows))
	for i, row := range rows {
		components[i] = &common.EnvironmentComponent{
			ID:            row.ID.Bytes,
			EnvironmentID: row.EnvironmentID.Bytes,
			Name:          row.Name,
			ChartName:     row.ChartName,
			ChartVersion:  row.ChartVersion,
			ChartRegistry: row.ChartRegistry,
		}
	}

	return components, nil
}

func (r *EnvironmentComponentsRepository) Delete(c *common.EnvironmentComponent) error {
	conn := db.New(r.pool)

	err := conn.DeleteEnvironmentComponentByID(context.Background(), pgtype.UUID{
		Bytes: [16]byte(c.ID[:]),
		Valid: true,
	})
	if err != nil {
		r.l.Error("failed to delete environment component", "error", err)
	}

	return err
}

func (r *EnvironmentComponentsRepository) Save(c *common.EnvironmentComponent) error {
	conn := db.New(r.pool)

	_, err := conn.UpsertEnvironmentComponent(context.Background(), db.UpsertEnvironmentComponentParams{
		ID: pgtype.UUID{
			Bytes: [16]byte(c.ID[:]),
			Valid: true,
		},
		EnvironmentID: pgtype.UUID{
			Bytes: [16]byte(c.EnvironmentID[:]),
			Valid: true,
		},
		Name:          c.Name,
		ChartName:     c.ChartName,
		ChartVersion:  c.ChartVersion,
		ChartRegistry: c.ChartRegistry,
	})

	if err != nil {
		r.l.Error("failed to save environment component", "error", err)
	}

	return err
}

func NewEnvironmentComponentsRepository(l *slog.Logger, pool *pgxpool.Pool) *EnvironmentComponentsRepository {
	return &EnvironmentComponentsRepository{
		pool: pool,
		l:    l.With("component", "infra.postgres.environment_components_repository"),
	}
}
