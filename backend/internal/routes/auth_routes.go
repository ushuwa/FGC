package routes

import (
	"github.com/gofiber/fiber/v3"

	"github.com/jj.jobo/FGC/internal/container"
	"github.com/jj.jobo/FGC/internal/middleware"
)

func AuthRoutes(
	router fiber.Router,
	c *container.Container,
) {

	auth := router.Group("/auth")

	auth.Post("/login", c.AuthHandler.Login)

	auth.Get(
		"/me",
		middleware.AuthMiddleware,
		func(ctx fiber.Ctx) error {

			return ctx.JSON(fiber.Map{
				"success": true,
				"data": fiber.Map{
					"id":       ctx.Locals("userID"),
					"username": ctx.Locals("username"),
					"role":     ctx.Locals("role"),
				},
			})
		},
	)
}
