package postgres

import (
	"Backend-trainee-assignment-spring-2025/internal/domain/models"
	"Backend-trainee-assignment-spring-2025/internal/repository"
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReceptionRepo struct {
	db     *pgxpool.Pool
	logger slog.Logger
}

func NewReceptionRepo(db *pgxpool.Pool, logger slog.Logger) repository.ReceptionRepository {
	return &ReceptionRepo{
		db:     db,
		logger: logger,
	}
}

func (r *ReceptionRepo) CreateReception(ctx context.Context, PVZID uuid.UUID) (models.Reception, error) {
	query := `INSERT INTO reception (pvz_id) VALUES ($1) RETURNING id,data_time,pvz_id,status`

	var Reception models.Reception

	err := r.db.QueryRow(ctx, query, PVZID).Scan(&Reception.ID, &Reception.DataTime, &Reception.PVZID, &Reception.Status)

	if err != nil {
		r.logger.Error("failed to create reception", "error", err)
		return models.Reception{}, err
	}

	return Reception, nil

}

func (r *ReceptionRepo) GetReception(ctx context.Context, id uuid.UUID) (models.Reception, error) {
	query := `SELECT id,data_time,pvz_id,status FROM reception WHERE id = $1`

	var Reception models.Reception

	err := r.db.QueryRow(ctx, query, id).Scan(&Reception.ID, &Reception.DataTime, &Reception.PVZID, &Reception.Status)

	if err != nil {
		r.logger.Error("failed to get reception", "error", err)
		return models.Reception{}, err
	}

	return Reception, nil
}

func (r *ReceptionRepo) GetAllReception(ctx context.Context) ([]models.Reception, error) {
	query := `SELECT id,data_time,pvz_id,status FROM reception`

	rows, err := r.db.Query(ctx, query)

	if err != nil {
		r.logger.Error("failed to get all reception", "error", err)
		return nil, err
	}

	var Reception []models.Reception

	for rows.Next() {
		var reception models.Reception

		err := rows.Scan(&reception.ID, &reception.DataTime, &reception.PVZID, &reception.Status)

		if err != nil {
			r.logger.Error("failed to scan reception", "error", err)
			return nil, err
		}

		Reception = append(Reception, reception)
	}

	return Reception, nil
}

func (r *ReceptionRepo) UpdateReceptionStatus(ctx context.Context, id uuid.UUID, status models.ReceptionStatus) (models.Reception, error) {
	query := `UPDATE reception SET status = $2 WHERE id = $1 RETURNING id,data_time,pvz_id,status`

	var Reception models.Reception

	err := r.db.QueryRow(ctx, query, id, status).Scan(&Reception.ID, &Reception.DataTime, &Reception.PVZID, &Reception.Status)

	if err != nil {
		r.logger.Error("failed to update reception", "error", err)
		return models.Reception{}, err
	}

	return Reception, nil
}
