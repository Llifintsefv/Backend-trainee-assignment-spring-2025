package service

import (
	"Backend-trainee-assignment-spring-2025/internal/domain/models"
	"context"
)

type Auth interface {
	DummyLogin(ctx context.Context, ip string, role models.Role) (models.TokenPair, error)
	RefreshToken(ctx context.Context, refreshToken string) (models.TokenPair, error)
	RegisterUser(ctx context.Context, registerRequest models.RegisterRequest) error
	LoginUser(ctx context.Context, loginRequest models.LoginRequest, ip string) (models.TokenPair, error)
}

type PvzService interface {
	CreatePvz(ctx context.Context, city models.City) (models.PVZ, error)
}
