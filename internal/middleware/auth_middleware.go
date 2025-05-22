package middleware

import (
	"Backend-trainee-assignment-spring-2025/internal/domain/models"
	"Backend-trainee-assignment-spring-2025/pkg/auth"
	"fmt"
	"log/slog"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware struct {
	TokenParser auth.Auth
	Logger      *slog.Logger
}

func NewAuthMiddleware(tokenParser auth.Auth, logger *slog.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		TokenParser: tokenParser,
		Logger:      logger,
	}
}

const (
	UserCtxKey    = "userID"
	RoleCtxKey    = "userRole"
	authHeader    = "Authorization"
	bearerPrefix  = "Bearer "
	RoleModerator = "moderator"
)

func (am *AuthMiddleware) AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get(authHeader)
		if header == "" {
			am.Logger.Error("authorization header is not provided")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"reason": "Заголовок Authorization отсутствует",
			})
		}

		if !strings.HasPrefix(header, bearerPrefix) {
			am.Logger.Error("invalid authorization header")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"reason": "Неверный заголовок Authorization",
			})
		}

		tokenString := strings.TrimPrefix(header, bearerPrefix)
		if tokenString == "" {
			am.Logger.Error("invalid authorization header")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"reason": "Неверный заголовок Authorization",
			})
		}

		claims, err := am.TokenParser.ParseAccessToken(tokenString)
		if err != nil {
			am.Logger.Error("failed to parse access token", "error", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"reason": "Неверный заголовок Authorization",
			})
		}

		c.Locals(UserCtxKey, claims.UserID)
		c.Locals(RoleCtxKey, claims.Role)
		am.Logger.Debug("extracted role from token", "role", claims.Role)
		am.Logger.Debug("extracted role from token", "rыфывole", claims.UserID)
		fmt.Println(claims.UserID, claims.Role)
		fmt.Println(c.Locals(RoleCtxKey))
		return c.Next()
	}
}

func (am *AuthMiddleware) ModeratorRoleMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Используем "comma-ok" идиому для безопасного приведения типа
		role, ok := c.Locals(RoleCtxKey).(models.Role)

		// Проверяем, успешно ли приведение типа И является ли роль модератором
		if !ok || role != RoleModerator {
			// Логируем фактическую роль (если удалось привести к string) и статус приведения
			am.Logger.Warn("access denied: user is not a moderator or role not found/invalid", "role", role, "ok", ok)
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"reason": "Доступ запрещен. Только модераторы могут выполнять это действие.",
			})
		}
		return c.Next()
	}
}
