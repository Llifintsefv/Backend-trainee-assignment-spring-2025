package repository

import (
	"Backend-trainee-assignment-spring-2025/internal/domain/models"
	"context"
)

type PVZRepository interface {
	CreatePVZ(ctx context.Context, city models.City) (models.PVZ, error)
}

type ReceptionRepository interface {
}

type UserRepository interface {
	SaveUser(ctx context.Context, email string, password string, role models.Role) error
	GetUser(ctx context.Context, email string) (models.User, error)
	GetUserByID(ctx context.Context, id string) (models.User, error)
	DeleteUser(ctx context.Context, email string) error
	ComparePassword(ctx context.Context, password string, hash string) (bool, error)
}

type ProductRepository interface {
}

type JwtRepository interface {
	SaveRefreshToken(ctx context.Context, refreshTokenData models.RefreshToken) error
	FindRefreshTokenByJTI(ctx context.Context, jti string) (models.RefreshToken, error)
	Delete(ctx context.Context, tokenHash string) error
}
