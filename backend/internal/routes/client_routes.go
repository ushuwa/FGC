package routes

import (
	"github.com/gofiber/fiber/v3"

	"github.com/jj.jobo/FGC/internal/container"
	"github.com/jj.jobo/FGC/internal/middleware"
)

func ClientRoutes(
	router fiber.Router,
	c *container.Container,
) {

	clients := router.Group(
		"/clients",
	)

	clients.Use(
		middleware.AuthMiddleware,
	)

	clients.Get(
		"",
		c.ClientHandler.GetClients,
	)

	clients.Get(
		"/:id/profile",
		c.ClientHandler.GetClientProfile,
	)

	clients.Get(
		"/:id",
		c.ClientHandler.GetClient,
	)

	clients.Post(
		"",
		c.ClientHandler.CreateClient,
	)

	clients.Put(
		"/:id",
		c.ClientHandler.UpdateClient,
	)

	clients.Delete(
		"/:id",
		c.ClientHandler.DeleteClient,
	)

}
