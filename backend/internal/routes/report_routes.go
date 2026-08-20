package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jj.jobo/FGC/internal/handlers"
)

func RegisterReportRoutes(
	router fiber.Router,
	handler *handlers.ReportHandler,
) {

	reports := router.Group(
		"/reports",
	)

	reports.Get(
		"/summary-collection/pdf",
		handler.GenerateCollectionPDF,
	)

	reports.Get(
		"/amortization/:loanID/pdf",
		handler.GenerateAmortizationPDF,
	)

	reports.Get(
		"/portfolio-at-risk/pdf",
		handler.GeneratePARPDF,
	)

	reports.Get(
		"/loan-portfolio/pdf",
		handler.GenerateLoanPortfolioPDF,
	)

	reports.Get(
		"/loan-maturity/pdf",
		handler.GenerateLoanMaturityPDF,
	)
}
