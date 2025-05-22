package auth

import (
	"Backend-trainee-assignment-spring-2025/internal/domain/models"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type auth struct {
	logger          *slog.Logger
	secretKey       string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewAuth(logger *slog.Logger, secretKey string, accessTokenTTL time.Duration, refreshTokenTTL time.Duration) Auth {
	return &auth{
		logger:          logger,
		secretKey:       secretKey,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

type Auth interface {
	GenerateAccessToken(userID string, JTI string, clientIP string, role models.Role) (string, error)
	GenerateRefreshToken(userID string, clientIP string) (string, models.RefreshToken, error)
	HashRefreshToken(refreshToken string) (string, error)
	ParseAccessToken(accessToken string) (models.AccessTokenClaims, error)
}

const (
	RefreshTokenSecretLength = 32
)

func (a *auth) GenerateAccessToken(userID string, JTI string, clientIP string, role models.Role) (string, error) {
	claims := models.AccessTokenClaims{
		UserID:    userID,
		IPAddress: clientIP,
		Role:      role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   userID,
			ID:        JTI,
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signetToken, err := accessToken.SignedString([]byte(a.secretKey))
	if err != nil {
		a.logger.Error("failed to sign access token", "error", err)
		return "", err
	}

	return signetToken, nil
}

func (a *auth) GenerateRefreshToken(userID, clientIP string) (string, models.RefreshToken, error) {
	JTI := uuid.New().String()
	secretByte := make([]byte, RefreshTokenSecretLength)
	_, err := rand.Read(secretByte)
	if err != nil {
		a.logger.Error("failed to generate refresh token", "error", err)
		return "", models.RefreshToken{}, err
	}

	secretPart := base64.StdEncoding.EncodeToString(secretByte)

	RefreshTokenHash, err := a.HashRefreshToken(secretPart)
	if err != nil {
		a.logger.Error("failed to hash refresh token", "error", err)
		return "", models.RefreshToken{}, err
	}

	clientToken := fmt.Sprintf("%s.%s", secretPart, JTI)

	return clientToken, models.RefreshToken{
		UserID:        userID,
		JTI:       JTI,
		TokenHash: RefreshTokenHash,
		IPAddress: clientIP,
		ExpiresAt: time.Now().Add(a.refreshTokenTTL),
	}, nil

}

func (a *auth) HashRefreshToken(refreshToken string) (string, error) {
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedToken), nil
}

func (a *auth) ParseAccessToken(accessToken string) (models.AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &models.AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.secretKey), nil
	})
	if err != nil {
		a.logger.Error("failed to parse access token", "error", err)
		return models.AccessTokenClaims{}, err
	}

	claims, ok := token.Claims.(*models.AccessTokenClaims)
	if !ok || !token.Valid {
		err := fmt.Errorf("invalid access token")
		a.logger.Error("invalid access token", "error", err)
		return models.AccessTokenClaims{}, err
	}

	return *claims, nil
}
