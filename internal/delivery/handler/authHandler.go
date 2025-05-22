package handler

import (
	"Backend-trainee-assignment-spring-2025/internal/domain/models"
	"Backend-trainee-assignment-spring-2025/internal/service"
	"Backend-trainee-assignment-spring-2025/pkg/validator"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

type authHandler struct {
	service service.Auth
	logger  *slog.Logger
}

func NewAuthHandler(service service.Auth, logger *slog.Logger) AuthHandler {
	return &authHandler{
		service: service,
		logger:  logger,
	}
}

func (h *authHandler) DummyLoginHandler(c *fiber.Ctx) error {
	ctx := c.UserContext()
	ip := c.IP()

	var request models.DummyLoginRequest

	err := c.BodyParser(&request)
	if err != nil {
		h.logger.Error("failed to parse role", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"reason": "Неверная роль",
		})
	}

	if err := validator.ValidateStruct(request); err != nil {
		h.logger.Error("failed to validate role", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"reason": "Неверная роль",
		})
	}

	TokenPair, err := h.service.DummyLogin(ctx, ip, request.Role)
	if err != nil {
		h.logger.Error("failed to login", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"reason": "Внутренняя ошибка сервера",
		})
	}
	return c.JSON(TokenPair)
}

func (h *authHandler) RefreshToken(c *fiber.Ctx) error {
	ctx := c.UserContext()
	var request models.RefreshRequest

	err := c.BodyParser(&request)
	if err != nil {
		h.logger.Error("failed to parse token", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"reason": "Неверный токен",
		})
	}

	if err := validator.ValidateStruct(request); err != nil {
		h.logger.Error("failed to validate token", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"reason": "Неверный токен",
		})
	}

	TokenPair, err := h.service.RefreshToken(ctx, request.RefreshToken)
	if err != nil {
		h.logger.Error("failed to validate token", "error", err)
		return err
	}

	return c.JSON(TokenPair)
}

func (h *authHandler) RegisterUser(c *fiber.Ctx) error {
	ctx := c.UserContext()
	var request models.RegisterRequest

	err := c.BodyParser(&request)
	if err != nil {
		h.logger.Error("failed to parse registration request", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"reason": "Неверный запрос",
		})
	}

	if err := validator.ValidateStruct(request); err != nil {
		h.logger.Error("failed to validate registration request", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"reason": "Неверные данные",
		})
	}

	if err := h.service.RegisterUser(ctx, request); err != nil {
		if err == models.ErrUserAlreadyExists {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"reason": "Пользователь с таким email уже существует"})
		}
		h.logger.Error("failed to register user", "error", err, "email", request.Email)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"reason": "Внутренняя ошибка сервера"})
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (h *authHandler) LoginUser(c *fiber.Ctx) error {
	ctx := c.UserContext()
	ip := c.IP()
	var request models.LoginRequest

	err := c.BodyParser(&request)
	if err != nil {
		h.logger.Error("failed to parse login request", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"reason": "Неверный запрос",
		})
	}

	if err := validator.ValidateStruct(request); err != nil {
		h.logger.Error("failed to validate login request", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"reason": "Неверные данные",
		})
	}

	tokenPair, err := h.service.LoginUser(ctx, request, ip)
	if err != nil {
		if err == models.ErrInvalidCredentials {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"reason": "Неверные учетные данные"})
		}
		h.logger.Error("failed to login user", "error", err, "email", request.Email)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"reason": "Внутренняя ошибка сервера"})
	}

	return c.JSON(tokenPair)
}
