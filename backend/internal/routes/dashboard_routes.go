package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jj.jobo/FGC/internal/handlers"
)

func DashboardRoutes(
	api fiber.Router,
	handler *handlers.DashboardHandler,
) {
	dashboard := api.Group("/dashboard")

	dashboard.Get(
		"/summary",
		handler.GetSummary,
	)
}
