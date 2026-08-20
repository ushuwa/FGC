package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/jj.jobo/FGC/internal/services"
)

type PARHandler struct {
	parService *services.PARService
}

func NewPARHandler(
	parService *services.PARService,
) *PARHandler {
	return &PARHandler{
		parService: parService,
	}
}

func (h *PARHandler) GetPAR(
	c fiber.Ctx,
) error {

	// ========================================
	// FILTERS
	// ========================================

	search := c.Query("search")
	status := c.Query("status")
	aging := c.Query("aging")

	// ========================================
	// PAGINATION
	// ========================================

	page := 1
	limit := 10

	if pageParam := c.Query("page"); pageParam != "" {

		parsedPage, err :=
			strconv.Atoi(pageParam)

		if err != nil || parsedPage < 1 {

			return c.Status(
				fiber.StatusBadRequest,
			).JSON(fiber.Map{
				"success": false,
				"message": "Invalid page parameter",
			})

		}

		page = parsedPage
	}

	if limitParam := c.Query("limit"); limitParam != "" {

		parsedLimit, err :=
			strconv.Atoi(limitParam)

		if err != nil || parsedLimit < 1 {

			return c.Status(
				fiber.StatusBadRequest,
			).JSON(fiber.Map{
				"success": false,
				"message": "Invalid limit parameter",
			})

		}

		if parsedLimit > 100 {
			parsedLimit = 100
		}

		limit = parsedLimit
	}

	// ========================================
	// SERVICE
	// ========================================

	result, err :=
		h.parService.GetPAR(
			search,
			status,
			aging,
			page,
			limit,
		)

	if err != nil {

		return c.Status(
			fiber.StatusInternalServerError,
		).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})

	}

	// ========================================
	// RESPONSE
	// ========================================

	return c.Status(
		fiber.StatusOK,
	).JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}
