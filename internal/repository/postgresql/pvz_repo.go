package postgres

import (
	"Backend-trainee-assignment-spring-2025/internal/domain/models"
	"Backend-trainee-assignment-spring-2025/internal/repository"
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PVZRepo struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewPVZRepo(db *pgxpool.Pool, logger *slog.Logger) repository.PVZRepository {
	return &PVZRepo{
		db:     db,
		logger: logger,
	}
}

func (r *PVZRepo) CreatePVZ(ctx context.Context, city models.City) (models.PVZ, error) {
	query := `INSERT INTO pvz (city) VALUES ($1) RETURNING id,city,registration_date`

	var PVZ models.PVZ

	err := r.db.QueryRow(ctx, query, city).Scan(&PVZ.ID, &PVZ.City, &PVZ.RegistrationsData)

	if err != nil {
		slog.Error("failed to create pvz", "error", err)
		return models.PVZ{}, err
	}

	return PVZ, nil
}

func (r *PVZRepo) GetPVZ(ctx context.Context, id uuid.UUID) (models.PVZ, error) {
	query := `SELECT id,city,registration_date FROM pvz WHERE id = $1`

	var PVZ models.PVZ

	err := r.db.QueryRow(ctx, query, id).Scan(&PVZ.ID, &PVZ.City, &PVZ.RegistrationsData)

	if err != nil {
		slog.Error("failed to get pvz", "error", err)
		return models.PVZ{}, err
	}

	return PVZ, nil
}

func (r *PVZRepo) GetAllPVZ(ctx context.Context) ([]models.PVZ, error) {
	query := `SELECT id,city,registration_date FROM pvz`

	rows, err := r.db.Query(ctx, query)

	if err != nil {
		slog.Error("failed to get all pvz", "error", err)
		return nil, err
	}

	var PVZ []models.PVZ

	for rows.Next() {
		var pvz models.PVZ

		err := rows.Scan(&pvz.ID, &pvz.City, &pvz.RegistrationsData)

		if err != nil {
			slog.Error("failed to scan pvz", "error", err)
			return nil, err
		}

		PVZ = append(PVZ, pvz)
	}

	return PVZ, nil
}

func (r *PVZRepo) UpdatePVZ(ctx context.Context, id uuid.UUID, city models.City) (models.PVZ, error) {
	query := `UPDATE pvz SET city = $2 WHERE id = $1 RETURNING id,city,registration_date`

	var PVZ models.PVZ

	err := r.db.QueryRow(ctx, query, id, city).Scan(&PVZ.ID, &PVZ.City, &PVZ.RegistrationsData)

	if err != nil {
		slog.Error("failed to update pvz", "error", err)
		return models.PVZ{}, err
	}

	return PVZ, nil
}

func (r *PVZRepo) DeletePVZ(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM pvz WHERE id = $1`

	_, err := r.db.Exec(ctx, query, id)

	if err != nil {
		slog.Error("failed to delete pvz", "error", err)
		return err
	}

	return nil
}
