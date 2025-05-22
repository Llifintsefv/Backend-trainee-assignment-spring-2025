// internal/router/router.go

package router

import (
	"Backend-trainee-assignment-spring-2025/internal/delivery/handler"
	mw "Backend-trainee-assignment-spring-2025/internal/middleware" // Используем mw как псевдоним для middleware

	"github.com/gofiber/fiber/v2"
)

func NewApp(authHandler handler.AuthHandler, authMw *mw.AuthMiddleware, pvzHandler handler.PvzHandler) *fiber.App {
	app := fiber.New()

	app.Post("/dummyLogin", authHandler.DummyLoginHandler)
	app.Post("/register", authHandler.RegisterUser)
	app.Post("/login", authHandler.LoginUser)

	protectedApiV1 := app.Group("/")
	protectedApiV1.Use(authMw.AuthMiddleware())
	protectedApiV1.Post("/test", func(c *fiber.Ctx) error { return c.SendString("OK") })
	protectedApiV1.Post("/refresh", authHandler.RefreshToken)

	protectedApiV1.Post("/pvz", authMw.ModeratorRoleMiddleware(), pvzHandler.CreatePvz)

	return app
}
