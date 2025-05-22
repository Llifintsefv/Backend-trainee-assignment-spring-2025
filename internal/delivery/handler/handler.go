package handler

import "github.com/gofiber/fiber/v2"

type AuthHandler interface {
	DummyLoginHandler(c *fiber.Ctx) error
	RefreshToken(c *fiber.Ctx) error
	RegisterUser(c *fiber.Ctx) error
	LoginUser(c *fiber.Ctx) error
}

type PvzHandler interface {
	CreatePvz(c *fiber.Ctx) error
}

type Handler struct {
	AuthHandler AuthHandler
}

func NewHandler(authHandler AuthHandler) Handler {
	return Handler{AuthHandler: authHandler}
}
