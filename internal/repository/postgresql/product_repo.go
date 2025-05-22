package postgres

import (
	"Backend-trainee-assignment-spring-2025/internal/domain/models"
	"Backend-trainee-assignment-spring-2025/internal/repository"
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepo struct {
	db     *pgxpool.Pool
	logger slog.Logger
}

func NewProductRepo(db *pgxpool.Pool, logger slog.Logger) repository.ProductRepository {
	return &ProductRepo{
		db:     db,
		logger: logger,
	}
}

func (p *ProductRepo) CreateProduct(ctx context.Context, product_type models.ProductType, reception_id uuid.UUID) (models.Product, error) {
	query := `INSERT INTO products (type,reception_id) VALUES ($1,$2) RETURNING id,data_time,type,reception_id`

	var product models.Product

	err := p.db.QueryRow(ctx, query, product_type, reception_id).Scan(&product.ID, &product.DataTime, &product.Type, &product.ReceptionID)

	if err != nil {
		p.logger.Error("failed to create product", "error", err)
		return models.Product{}, err
	}

	return product, nil
}

func (p *ProductRepo) GetProduct(ctx context.Context, id uuid.UUID) (models.Product, error) {
	query := `SELECT id,data_time,type,reception_id FROM products WHERE id = $1`

	var product models.Product

	err := p.db.QueryRow(ctx, query, id).Scan(&product.ID, &product.DataTime, &product.Type, &product.ReceptionID)

	if err != nil {
		p.logger.Error("failed to get product", "error", err)
		return models.Product{}, err
	}

	return product, nil
}

func (p *ProductRepo) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM products WHERE id = $1`

	_, err := p.db.Exec(ctx, query, id)

	if err != nil {
		p.logger.Error("failed to delete product", "error", err)
		return err
	}

	return nil
}
