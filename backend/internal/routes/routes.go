package routes

import (
	"github.com/gofiber/fiber/v3"

	"github.com/jj.jobo/FGC/internal/container"
)

func Setup(
	app *fiber.App,
	c *container.Container,
) {

	api := app.Group("/api")

	v1 := api.Group("/v1")

	HealthRoutes(v1)

	AuthRoutes(
		v1,
		c,
	)

	ClientRoutes(
		v1,
		c,
	)

	LoanRoutes(
		v1,
		c.LoanHandler,
	)
	DashboardRoutes(
		v1,
		c.DashboardHandler,
	)
	PARRoutes(
		v1,
		c.PARHandler,
	)

	RegisterReportRoutes(
		v1,
		c.ReportHandler,
	)
}
