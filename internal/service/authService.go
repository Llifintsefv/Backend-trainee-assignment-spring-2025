package service

import (
	"Backend-trainee-assignment-spring-2025/internal/domain/models"
	"Backend-trainee-assignment-spring-2025/internal/repository"
	"Backend-trainee-assignment-spring-2025/pkg/auth"
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
)

type authService struct {
	jwtRepo  repository.JwtRepository
	userRepo repository.UserRepository
	logger   *slog.Logger
	auth     auth.Auth
}

func NewAuthService(jwtRepo repository.JwtRepository, logger *slog.Logger, auth auth.Auth, userRepo repository.UserRepository) Auth {
	return &authService{
		jwtRepo:  jwtRepo,
		logger:   logger,
		auth:     auth,
		userRepo: userRepo,
	}
}

func (s *authService) DummyLogin(ctx context.Context, ip string, role models.Role) (models.TokenPair, error) {

	userID := uuid.New().String()

	refreshToken, refreshTokenData, err := s.auth.GenerateRefreshToken(userID, ip)
	if err != nil {
		s.logger.Error("failed to generate refresh token", "error", err)
		return models.TokenPair{}, err
	}

	err = s.jwtRepo.SaveRefreshToken(ctx, refreshTokenData)
	if err != nil {
		s.logger.Error("failed to save refresh token", "error", err)
		return models.TokenPair{}, err
	}

	AccessToken, err := s.auth.GenerateAccessToken(userID, refreshTokenData.JTI, ip, role)
	if err != nil {
		s.logger.Error("failed to generate access token", "error", err)
		return models.TokenPair{}, err
	}

	return models.TokenPair{
		AccessToken:  AccessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshTokenRequest string) (models.TokenPair, error) {
	parts := strings.Split(refreshTokenRequest, ".")
	JTI := parts[1]

	refreshToken, err := s.jwtRepo.FindRefreshTokenByJTI(ctx, JTI)
	if err != nil {
		s.logger.Error("failed to find refresh token", "error", err)
		return models.TokenPair{}, err
	}

	if time.Now().After(refreshToken.ExpiresAt) {
		s.logger.Error("refresh token expired", "error", err)
		return models.TokenPair{}, err
	}

	user, err := s.userRepo.GetUserByID(ctx, refreshToken.UserID)
	if err != nil {
		s.logger.Error("failed to find user", "error", err)
		return models.TokenPair{}, err
	}
	role := user.Role

	err = s.jwtRepo.Delete(ctx, refreshToken.TokenHash)
	if err != nil {
		s.logger.Error("failed to delete refresh token", "error", err)
		return models.TokenPair{}, err
	}

	newRefreshToken, newRefreshTokenData, err := s.auth.GenerateRefreshToken(refreshToken.UserID, refreshToken.IPAddress)
	if err != nil {
		s.logger.Error("failed to generate refresh token", "error", err)
		return models.TokenPair{}, err
	}

	err = s.jwtRepo.SaveRefreshToken(ctx, newRefreshTokenData)
	if err != nil {
		s.logger.Error("failed to save refresh token", "error", err)
		return models.TokenPair{}, err
	}

	AccessToken, err := s.auth.GenerateAccessToken(newRefreshTokenData.UserID, newRefreshTokenData.JTI, newRefreshTokenData.IPAddress, role)
	if err != nil {
		s.logger.Error("failed to generate access token", "error", err)
		return models.TokenPair{}, err
	}

	return models.TokenPair{
		AccessToken:  AccessToken,
		RefreshToken: newRefreshToken,
	}, nil

}

func (s *authService) RegisterUser(ctx context.Context, registerRequest models.RegisterRequest) error {
	_, err := s.userRepo.GetUser(ctx, registerRequest.Email)
	if err == nil {
		s.logger.Warn("user with this email already exists", "email", registerRequest.Email)
		return models.ErrUserAlreadyExists
	} else if err != models.ErrUserNotFound {
		s.logger.Error("failed to check if user exists", "error", err, "email", registerRequest.Email)
		return err
	}

	err = s.userRepo.SaveUser(ctx, registerRequest.Email, registerRequest.Password, registerRequest.Role)
	if err != nil {
		s.logger.Error("failed to save user", "error", err, "email", registerRequest.Email)
		return err
	}

	return nil
}

func (s *authService) LoginUser(ctx context.Context, loginRequest models.LoginRequest, ip string) (models.TokenPair, error) {
	user, err := s.userRepo.GetUser(ctx, loginRequest.Email)
	if err != nil {
		s.logger.Error("failed to get user", "error", err, "email", loginRequest.Email)
		return models.TokenPair{}, err
	}

	valid, err := s.userRepo.ComparePassword(ctx, loginRequest.Password, user.PasswordHash)
	if err != nil {
		s.logger.Error("failed to compare password", "error", err, "email", loginRequest.Email)
		return models.TokenPair{}, err
	}

	if !valid {
		s.logger.Warn("invalid credentials", "email", loginRequest.Email)
		return models.TokenPair{}, models.ErrInvalidCredentials
	}

	tokenPair, err := s.DummyLogin(ctx, ip, user.Role)
	if err != nil {
		s.logger.Error("failed to perform dummy login", "error", err)
		return models.TokenPair{}, err
	}

	return tokenPair, nil
}
