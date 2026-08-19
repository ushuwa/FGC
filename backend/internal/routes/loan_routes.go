package routes

import (
	"github.com/gofiber/fiber/v3"

	"github.com/jj.jobo/FGC/internal/handlers"
)

func LoanRoutes(
	router fiber.Router,
	handler *handlers.LoanHandler,
) {

	loans := router.Group("/loans")

	loans.Get(
		"/",
		handler.GetLoans,
	)
	loans.Post(
		"/",
		handler.CreateLoan,
	)

	loans.Get(
		"/:id",
		handler.GetLoan,
	)

	loans.Put(
		"/:id",
		handler.UpdateLoan,
	)

	loans.Get(
		"/:id/payments",
		handler.GetPayments,
	)

	loans.Post(
		"/:id/payments",
		handler.CreatePayment,
	)

	loans.Delete(
		"/:id/payments/:paymentId",
		handler.DeletePayment,
	)

	loans.Post("/:id/rebuild-amortization",
		handler.RebuildAmortization,
	)

}
