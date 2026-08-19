package handlers

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"

	"github.com/jj.jobo/FGC/internal/dto"
	"github.com/jj.jobo/FGC/internal/services"
)

type LoanHandler struct {
	service *services.LoanService
}

func NewLoanHandler(
	service *services.LoanService,
) *LoanHandler {
	return &LoanHandler{
		service: service,
	}
}

func (
	h *LoanHandler,
) GetLoans(
	c fiber.Ctx,
) error {

	search := strings.TrimSpace(
		c.Query("search"),
	)

	status := strings.TrimSpace(
		c.Query("status"),
	)

	loans, err :=
		h.service.GetLoans(
			search,
			status,
		)

	if err != nil {
		return c.Status(
			fiber.StatusInternalServerError,
		).JSON(
			fiber.Map{
				"success": false,
				"message": "Failed to retrieve loans",
				"error":   err.Error(),
			},
		)
	}

	return c.JSON(
		fiber.Map{
			"success": true,
			"message": "Loans retrieved successfully",
			"data":    loans,
		},
	)
}

func (
	h *LoanHandler,
) GetLoan(
	c fiber.Ctx,
) error {

	idString := c.Params("id")

	id, err := strconv.Atoi(idString)

	if err != nil || id <= 0 {
		return c.Status(
			fiber.StatusBadRequest,
		).JSON(
			fiber.Map{
				"success": false,
				"message": "Invalid loan ID",
			},
		)
	}

	loan, err :=
		h.service.GetLoan(id)

	if err != nil {

		if err.Error() ==
			"loan not found" {

			return c.Status(
				fiber.StatusNotFound,
			).JSON(
				fiber.Map{
					"success": false,
					"message": "Loan not found",
				},
			)
		}

		return c.Status(
			fiber.StatusInternalServerError,
		).JSON(
			fiber.Map{
				"success": false,
				"message": "Failed to retrieve loan",
				"error":   err.Error(),
			},
		)
	}

	return c.JSON(
		fiber.Map{
			"success": true,
			"message": "Loan retrieved successfully",
			"data":    loan,
		},
	)
}

func (
	h *LoanHandler,
) GetPayments(
	c fiber.Ctx,
) error {

	idString := c.Params("id")

	id, err := strconv.Atoi(
		idString,
	)

	if err != nil || id <= 0 {
		return c.Status(
			fiber.StatusBadRequest,
		).JSON(
			fiber.Map{
				"success": false,
				"message": "Invalid loan ID",
			},
		)
	}

	payments, err :=
		h.service.GetPayments(id)

	if err != nil {
		return c.Status(
			fiber.StatusInternalServerError,
		).JSON(
			fiber.Map{
				"success": false,
				"message": "Failed to retrieve payment history",
				"error":   err.Error(),
			},
		)
	}

	return c.JSON(
		fiber.Map{
			"success": true,
			"message": "Payment history retrieved successfully",
			"data":    payments,
		},
	)
}
func (
	h *LoanHandler,
) CreatePayment(
	c fiber.Ctx,
) error {

	idString :=
		c.Params("id")

	id, err :=
		strconv.Atoi(idString)

	if err != nil || id <= 0 {
		return c.Status(
			fiber.StatusBadRequest,
		).JSON(
			fiber.Map{
				"success": false,
				"message": "Invalid loan ID",
			},
		)
	}

	var req dto.CreatePaymentRequest

	if err :=
		c.Bind().Body(&req); err != nil {

		return c.Status(
			fiber.StatusBadRequest,
		).JSON(
			fiber.Map{
				"success": false,
				"message": "Invalid request body",
			},
		)
	}

	payment, err :=
		h.service.CreatePayment(
			id,
			req,
		)

	if err != nil {

		status :=
			fiber.StatusInternalServerError

		switch err.Error() {

		case "invalid loan ID",
			"payment date is required",
			"payment amount must be greater than zero":

			status =
				fiber.StatusBadRequest

		case "loan not found",
			"loan has no amortization schedule":

			status =
				fiber.StatusNotFound

		case "payment exceeds loan outstanding balance":

			status =
				fiber.StatusBadRequest

		case "payment could not be fully allocated":

			status =
				fiber.StatusBadRequest
		}

		return c.Status(status).JSON(
			fiber.Map{
				"success": false,
				"message": err.Error(),
			},
		)
	}

	return c.Status(
		fiber.StatusCreated,
	).JSON(
		fiber.Map{
			"success": true,
			"message": "Payment created successfully",
			"data":    payment,
		},
	)
}

func (h *LoanHandler) DeletePayment(
	c fiber.Ctx,
) error {

	id, err := strconv.Atoi(
		c.Params("paymentId"),
	)

	if err != nil || id <= 0 {
		return c.Status(
			fiber.StatusBadRequest,
		).JSON(fiber.Map{
			"success": false,
			"message": "Invalid payment ID",
		})
	}

	if err := h.service.DeletePayment(id); err != nil {
		return c.Status(
			fiber.StatusBadRequest,
		).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.Status(
		fiber.StatusOK,
	).JSON(fiber.Map{
		"success": true,
		"message": "Payment deleted successfully",
	})
}

func (h *LoanHandler) RebuildAmortization(
	c fiber.Ctx,
) error {

	id, err := strconv.Atoi(
		c.Params("id"),
	)

	if err != nil || id <= 0 {
		return c.Status(
			fiber.StatusBadRequest,
		).JSON(fiber.Map{
			"success": false,
			"message": "invalid loan ID",
		})
	}

	err = h.service.RebuildAmortization(
		id,
	)

	if err != nil {
		return c.Status(
			fiber.StatusBadRequest,
		).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "amortization schedule rebuilt successfully",
	})
}

func (h *LoanHandler) CreateLoan(
	c fiber.Ctx,
) error {

	var req dto.CreateLoanRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(
			fiber.StatusBadRequest,
		).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	if err := h.service.CreateLoan(req); err != nil {
		return c.Status(
			fiber.StatusBadRequest,
		).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.Status(
		fiber.StatusCreated,
	).JSON(fiber.Map{
		"success": true,
		"message": "Loan created successfully",
	})
}
func (h *LoanHandler) UpdateLoan(
	c fiber.Ctx,
) error {

	id, err := strconv.Atoi(
		c.Params("id"),
	)

	if err != nil || id <= 0 {
		return c.Status(
			fiber.StatusBadRequest,
		).JSON(fiber.Map{
			"success": false,
			"message": "Invalid loan ID",
		})
	}

	var req dto.UpdateLoanRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(
			fiber.StatusBadRequest,
		).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	if err := h.service.UpdateLoan(
		id,
		req,
	); err != nil {

		status := fiber.StatusBadRequest

		if err.Error() == "loan not found" {
			status = fiber.StatusNotFound
		}

		return c.Status(status).JSON(
			fiber.Map{
				"success": false,
				"message": err.Error(),
			},
		)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Loan updated successfully",
	})
}
