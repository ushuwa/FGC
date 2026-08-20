package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jj.jobo/FGC/internal/services"
)

type DashboardHandler struct {
	service *services.DashboardService
}

func NewDashboardHandler(
	service *services.DashboardService,
) *DashboardHandler {
	return &DashboardHandler{
		service: service,
	}
}

func (h *DashboardHandler) GetSummary(
	c fiber.Ctx,
) error {

	summary, err := h.service.GetSummary()

	if err != nil {
		return c.Status(
			fiber.StatusInternalServerError,
		).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.Status(
		fiber.StatusOK,
	).JSON(fiber.Map{
		"success": true,
		"data":    summary,
	})
}
