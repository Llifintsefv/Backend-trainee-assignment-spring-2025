package postgres

import (
	"Backend-trainee-assignment-spring-2025/internal/domain/models"
	"Backend-trainee-assignment-spring-2025/internal/repository"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type jwtrepo struct {
	db              *pgxpool.Pool
	logger          *slog.Logger
	secretKey       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func NewJWTRepo(db *pgxpool.Pool, logger *slog.Logger, secretKey string, accessTokenTTL time.Duration, refreshTokenTTL time.Duration) repository.JwtRepository {
	return &jwtrepo{
		db:              db,
		logger:          logger,
		secretKey:       secretKey,
		AccessTokenTTL:  accessTokenTTL,
		RefreshTokenTTL: refreshTokenTTL,
	}
}

func (j *jwtrepo) SaveRefreshToken(ctx context.Context, refreshTokenData models.RefreshToken) error {
	query := `INSERT INTO refresh_tokens (jti, user_id, token_hash, ip, expires_at,created_at) VALUES ($1, $2, $3, $4, $5,$6)`

	_, err := j.db.Exec(ctx, query, refreshTokenData.JTI, refreshTokenData.UserID, refreshTokenData.TokenHash, refreshTokenData.IPAddress, refreshTokenData.ExpiresAt, time.Now())
	if err != nil {
		j.logger.Error("failed to save refresh token", "error", err)
		return err
	}
	fmt.Println(refreshTokenData.JTI)

	return nil
}

func (j *jwtrepo) FindRefreshTokenByJTI(ctx context.Context, jti string) (models.RefreshToken, error) {
	fmt.Println(jti)
	var refreshToken models.RefreshToken

	query := `SELECT * FROM refresh_tokens WHERE jti = $1`
	err := j.db.QueryRow(ctx, query, jti).Scan(&refreshToken.ID,&refreshToken.JTI, &refreshToken.UserID, &refreshToken.TokenHash, &refreshToken.IPAddress, &refreshToken.ExpiresAt, &refreshToken.CreatedAt)
	if err != nil {
		j.logger.Error("failed to find refresh token", "error", err)
		return models.RefreshToken{}, err
	}
	return refreshToken, nil

}

func (j *jwtrepo) Delete(ctx context.Context, tokenHash string) error {
	query := `DELETE FROM refresh_tokens WHERE token_hash = $1`
	_, err := j.db.Exec(ctx, query, tokenHash)
	if err != nil {
		j.logger.Error("failed to delete refresh token", "error", err)
		return err
	}
	return nil

}
