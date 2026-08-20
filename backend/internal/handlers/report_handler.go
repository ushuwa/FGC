package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v3"

	"github.com/jj.jobo/FGC/internal/services"
)

type ReportHandler struct {
	reportService    *services.ReportService
	reportPDFService *services.ReportPDFService
}

func NewReportHandler(
	reportService *services.ReportService,
	reportPDFService *services.ReportPDFService,
) *ReportHandler {
	return &ReportHandler{
		reportService:    reportService,
		reportPDFService: reportPDFService,
	}
}

func (h *ReportHandler) GenerateCollectionPDF(
	c fiber.Ctx,
) error {

	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	data, err := h.reportService.GetCollectionReport(
		dateFrom,
		dateTo,
	)

	if err != nil {
		return c.Status(
			fiber.StatusBadRequest,
		).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	path, err :=
		h.reportPDFService.GenerateCollectionReport(
			data,
			dateFrom,
			dateTo,
		)

	if err != nil {
		return c.Status(
			fiber.StatusInternalServerError,
		).JSON(fiber.Map{
			"success": false,
			"message": "Failed to generate PDF",
		})
	}

	return c.Download(
		path,
		"summary_collection.pdf",
	)
}

// ========================================
// LOAN ACCOUNT AMORTIZATION PDF
// ========================================

func (h *ReportHandler) GenerateAmortizationPDF(
	c fiber.Ctx,
) error {

	loanIDStr := c.Params("loanID")
	loanID := 0
	_, err := fmt.Sscanf(loanIDStr, "%d", &loanID)

	if err != nil || loanID <= 0 {
		return c.Status(
			fiber.StatusBadRequest,
		).JSON(fiber.Map{
			"success": false,
			"message": "Invalid loan ID",
		})
	}

	data, err :=
		h.reportService.GetAmortizationReport(
			int64(loanID),
		)

	if err != nil {

		return c.Status(
			fiber.StatusNotFound,
		).JSON(fiber.Map{
			"success": false,
			"message": "Loan amortization data not found",
		})
	}

	path, err :=
		h.reportPDFService.GenerateAmortizationReport(
			data,
		)

	if err != nil {

		return c.Status(
			fiber.StatusInternalServerError,
		).JSON(fiber.Map{
			"success": false,
			"message": "Failed to generate amortization PDF",
		})
	}

	return c.Download(
		path,
		"loan_amortization.pdf",
	)
}

// ========================================
// PORTFOLIO AT RISK PDF
// ========================================

func (h *ReportHandler) GeneratePARPDF(
	c fiber.Ctx,
) error {

	data, err :=
		h.reportService.GetPARReport()

	if err != nil {

		fmt.Println(
			"Failed to get PAR report:",
			err,
		)

		return c.Status(
			fiber.StatusInternalServerError,
		).JSON(fiber.Map{
			"success": false,
			"message": "Unable to generate Portfolio at Risk report.",
		})
	}

	path, err :=
		h.reportPDFService.GeneratePARReport(
			data,
		)

	if err != nil {

		fmt.Println(
			"Failed to generate PAR PDF:",
			err,
		)

		return c.Status(
			fiber.StatusInternalServerError,
		).JSON(fiber.Map{
			"success": false,
			"message": "Unable to generate Portfolio at Risk PDF.",
		})
	}

	return c.Download(
		path,
		"portfolio_at_risk.pdf",
	)
}

// ========================================
// LOAN PORTFOLIO SUMMARY PDF
// ========================================

func (h *ReportHandler) GenerateLoanPortfolioPDF(
	c fiber.Ctx,
) error {

	data, err :=
		h.reportService.GetLoanPortfolioReport()

	if err != nil {
		fmt.Println(
			"Failed to get loan portfolio report:",
			err,
		)

		return c.Status(
			fiber.StatusInternalServerError,
		).JSON(fiber.Map{
			"success": false,
			"message": "Unable to generate loan portfolio report.",
		})
	}

	path, err :=
		h.reportPDFService.GenerateLoanPortfolioReport(
			data,
		)

	if err != nil {
		fmt.Println(
			"Failed to generate loan portfolio PDF:",
			err,
		)

		return c.Status(
			fiber.StatusInternalServerError,
		).JSON(fiber.Map{
			"success": false,
			"message": "Unable to generate loan portfolio PDF.",
		})
	}

	return c.Download(
		path,
		"loan_portfolio_summary.pdf",
	)
}

// ========================================
// LOAN MATURITY / DUE PDF
// ========================================

func (h *ReportHandler) GenerateLoanMaturityPDF(
	c fiber.Ctx,
) error {

	data, err :=
		h.reportService.GetLoanMaturityReport()

	if err != nil {

		fmt.Println(
			"Failed to get loan maturity report:",
			err,
		)

		return c.Status(
			fiber.StatusInternalServerError,
		).JSON(fiber.Map{
			"success": false,
			"message": "Unable to generate loan maturity report.",
		})
	}

	path, err :=
		h.reportPDFService.GenerateLoanMaturityReport(
			data,
		)

	if err != nil {

		fmt.Println(
			"Failed to generate loan maturity PDF:",
			err,
		)

		return c.Status(
			fiber.StatusInternalServerError,
		).JSON(fiber.Map{
			"success": false,
			"message": "Unable to generate loan maturity PDF.",
		})
	}

	return c.Download(
		path,
		"loan_maturity_due.pdf",
	)
}
