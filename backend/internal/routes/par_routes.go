package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jj.jobo/FGC/internal/handlers"
)

func PARRoutes(
	api fiber.Router,
	handler *handlers.PARHandler,
) {
	par := api.Group("/portfolio-at-risk")

	par.Get(
		"/",
		handler.GetPAR,
	)
}
