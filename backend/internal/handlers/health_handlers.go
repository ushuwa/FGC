package handlers

import (
	"time"

	"github.com/gofiber/fiber/v3"

	"github.com/jj.jobo/FGC/internal/utils"
)

func Health(c fiber.Ctx) error {

	return utils.Success(
		c,
		"API Running",
		fiber.Map{
			"version": "1.0.0",
			"status":  "healthy",
			"time":    time.Now().UTC(),
		},
	)
}
