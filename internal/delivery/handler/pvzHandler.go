package handler

import (
	"Backend-trainee-assignment-spring-2025/internal/domain/models"
	"Backend-trainee-assignment-spring-2025/internal/service"
	"Backend-trainee-assignment-spring-2025/pkg/validator"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

type pvzHandler struct {
	service service.PvzService
	logger  *slog.Logger
}

func NewPvzHandler(service service.PvzService, logger *slog.Logger) PvzHandler {
	return &pvzHandler{
		service: service,
		logger:  logger,
	}
}

func (h *pvzHandler) CreatePvz(c *fiber.Ctx) error {
	ctx := c.UserContext()

	var request models.PVZRequest

	err := c.BodyParser(&request)
	if err != nil {
		h.logger.Error("failed to parse request", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"reason": "Неверный запрос",
		})

	}

	if err := validator.ValidateStruct(request); err != nil {
		h.logger.Error("failed to validate request", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"reason": "Неверные данные",
		})

	}

	pvz, err := h.service.CreatePvz(ctx, request.City)
	if err != nil {
		h.logger.Error("failed to create pvz", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"reason": "Внутренняя ошибка сервера",
		})

	}

	return c.JSON(pvz)
}
