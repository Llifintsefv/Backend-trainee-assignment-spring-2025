package config

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBConnStr       string
	Port            string
	SecretKey       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func NewConfig() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		slog.Error("DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, DB_SSL_MODE must be set")
		return nil, fmt.Errorf("missing required environment variables")
	}
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSL_MODE")

	if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" || dbSSLMode == "" {
		slog.Error("DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, DB_SSL_MODE must be set")
		return nil, fmt.Errorf("missing required environment variables")
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		slog.Error("SECRET_KEY must be set")
		return nil, fmt.Errorf("missing required environment variables")
	}

	accessTokenTTL := os.Getenv("ACCESS_TOKEN_TTL")
	if accessTokenTTL == "" {
		accessTokenTTL = "15m"
	}

	refreshTokenTTL := os.Getenv("REFRESH_TOKEN_TTL")
	if refreshTokenTTL == "" {
		refreshTokenTTL = "24h"
	}

	return &Config{
		DBConnStr: connStr,
		Port:      port,
		SecretKey: secretKey,
		AccessTokenTTL: func() time.Duration {
			ttl, _ := time.ParseDuration(accessTokenTTL)
			return ttl
		}(),
		RefreshTokenTTL: func() time.Duration {
			ttl, _ := time.ParseDuration(refreshTokenTTL)
			return ttl
		}(),
	}, nil
}
