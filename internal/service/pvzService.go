package service

import (
	"Backend-trainee-assignment-spring-2025/internal/domain/models"
	"Backend-trainee-assignment-spring-2025/internal/repository"
	"context"
	"log/slog"
)

type pvzService struct {
	repo   repository.PVZRepository
	logger *slog.Logger
}

func NewPvzService(repo repository.PVZRepository, logger *slog.Logger) PvzService {
	return &pvzService{
		repo:   repo,
		logger: logger,
	}
}

func (s *pvzService) CreatePvz(ctx context.Context, city models.City) (models.PVZ, error) {

	pvz, err := s.repo.CreatePVZ(ctx, city)
	if err != nil {
		s.logger.Error("failed to create pvz", "error", err)
		return models.PVZ{}, err
	}

	return pvz, nil
}
