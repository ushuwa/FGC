package services

import (
	"errors"
	"time"

	"github.com/jj.jobo/FGC/internal/dto"
	"github.com/jj.jobo/FGC/internal/repositories"
)

type ReportService struct {
	reportRepository *repositories.ReportRepository
}

func NewReportService(
	reportRepository *repositories.ReportRepository,
) *ReportService {
	return &ReportService{
		reportRepository: reportRepository,
	}
}

func (s *ReportService) GetCollectionReport(
	dateFrom string,
	dateTo string,
) (*dto.CollectionReportData, error) {

	if dateFrom == "" {
		return nil, errors.New(
			"date_from is required",
		)
	}

	if dateTo == "" {
		return nil, errors.New(
			"date_to is required",
		)
	}

	from, err := time.Parse(
		"2006-01-02",
		dateFrom,
	)

	if err != nil {
		return nil, errors.New(
			"invalid date_from",
		)
	}

	to, err := time.Parse(
		"2006-01-02",
		dateTo,
	)

	if err != nil {
		return nil, errors.New(
			"invalid date_to",
		)
	}

	if from.After(to) {
		return nil, errors.New(
			"date_from cannot be later than date_to",
		)
	}

	return s.reportRepository.GetCollectionReport(
		dateFrom,
		dateTo,
	)
}

// ========================================
// LOAN ACCOUNT AMORTIZATION
// ========================================

func (s *ReportService) GetAmortizationReport(
	loanID int64,
) (*dto.AmortizationReportData, error) {

	if loanID <= 0 {
		return nil, errors.New(
			"invalid loan_id",
		)
	}

	return s.reportRepository.GetAmortizationReport(
		loanID,
	)
}

// ========================================
// PORTFOLIO AT RISK REPORT
// ========================================

func (s *ReportService) GetPARReport() (
	*dto.PARReportData,
	error,
) {
	return s.reportRepository.GetPARReport()
}

// ========================================
// LOAN PORTFOLIO SUMMARY
// ========================================

func (s *ReportService) GetLoanPortfolioReport() (
	*dto.LoanPortfolioReportData,
	error,
) {
	return s.reportRepository.GetLoanPortfolioReport()
}

// ========================================
// LOAN MATURITY / DUE REPORT
// ========================================

func (s *ReportService) GetLoanMaturityReport() (
	*dto.LoanMaturityReportData,
	error,
) {
	return s.reportRepository.GetLoanMaturityReport()
}
