package routes

import (
	"github.com/gofiber/fiber/v3"

	"github.com/jj.jobo/FGC/internal/handlers"
)

func HealthRoutes(router fiber.Router) {

	router.Get(
		"/health",
		handlers.Health,
	)
}
