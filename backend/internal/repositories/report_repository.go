package repositories

import (
	"sort"
	"strings"

	"github.com/jj.jobo/FGC/internal/database"
	"github.com/jj.jobo/FGC/internal/dto"
	"gorm.io/gorm"
)

type ReportRepository struct {
	*BaseRepository[struct{}]
}

func NewReportRepository() *ReportRepository {
	return &ReportRepository{
		BaseRepository: NewBaseRepository[struct{}](
			database.DB,
		),
	}
}

// ========================================
// SUMMARY OF COLLECTION
// ========================================

func (r *ReportRepository) GetCollectionReport(
	dateFrom string,
	dateTo string,
) (*dto.CollectionReportData, error) {

	type collectionRow struct {
		PaymentDate        string  `gorm:"column:payment_date"`
		ReferenceNumber    string  `gorm:"column:reference_number"`
		PNNumber           string  `gorm:"column:pn_number"`
		ClientName         string  `gorm:"column:client_name"`
		PrincipalApplied   float64 `gorm:"column:principal_applied"`
		InterestApplied    float64 `gorm:"column:interest_applied"`
		AmountPaid         float64 `gorm:"column:amount_paid"`
		PaymentChannel     string  `gorm:"column:payment_channel"`
		OutstandingBalance float64 `gorm:"column:outstanding_balance"`
	}

	var rows []collectionRow

	err := database.DB.Raw(`
		SELECT *
		FROM public.fn_report_summary_collection(
			?,
			?
		)
	`,
		dateFrom,
		dateTo,
	).Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	result := &dto.CollectionReportData{
		Rows: make(
			[]dto.CollectionReportRow,
			0,
			len(rows),
		),
	}

	for _, row := range rows {

		result.Rows = append(
			result.Rows,
			dto.CollectionReportRow{
				PaymentDate:        row.PaymentDate,
				ReferenceNumber:    row.ReferenceNumber,
				PNNumber:           row.PNNumber,
				ClientName:         row.ClientName,
				PrincipalApplied:   row.PrincipalApplied,
				InterestApplied:    row.InterestApplied,
				AmountPaid:         row.AmountPaid,
				PaymentChannel:     row.PaymentChannel,
				OutstandingBalance: row.OutstandingBalance,
			},
		)

		result.Summary.TotalTransactions++

		result.Summary.TotalPrincipal +=
			row.PrincipalApplied

		result.Summary.TotalInterest +=
			row.InterestApplied

		result.Summary.TotalCollection +=
			row.AmountPaid

	}

	/*
		Outstanding balance is a snapshot
		per payment, so we should NOT sum it.
	*/

	if len(rows) > 0 {
		result.Summary.TotalOutstanding =
			rows[len(rows)-1].OutstandingBalance
	}

	return result, nil
}

// ========================================
// LOAN ACCOUNT AMORTIZATION
// ========================================

func (r *ReportRepository) GetAmortizationReport(
	loanID int64,
) (*dto.AmortizationReportData, error) {

	type amortizationRow struct {
		LoanID             int64   `gorm:"column:loan_id"`
		PNNumber           string  `gorm:"column:pn_number"`
		ClientName         string  `gorm:"column:client_name"`
		ContactNumber      string  `gorm:"column:contact_number"`
		LoanType           string  `gorm:"column:loan_type"`
		PrincipalAmount    float64 `gorm:"column:principal_amount"`
		InterestRate       float64 `gorm:"column:interest_rate"`
		LoanInterest       float64 `gorm:"column:loan_interest"`
		PNValue            float64 `gorm:"column:pn_value"`
		LoanTerm           int64   `gorm:"column:loan_term"`
		AmortizationAmount float64 `gorm:"column:amortization_amount"`
		Frequency          string  `gorm:"column:frequency"`
		DisbursementDate   string  `gorm:"column:disbursement_date"`
		MaturityDate       string  `gorm:"column:maturity_date"`
		LoanStatus         string  `gorm:"column:loan_status"`

		ScheduleID         int     `gorm:"column:schedule_id"`
		ScheduleNo         int64   `gorm:"column:schedule_no"`
		DueDate            string  `gorm:"column:due_date"`
		ScheduledPrincipal float64 `gorm:"column:scheduled_principal"`
		ScheduledInterest  float64 `gorm:"column:scheduled_interest"`
		ScheduledTotal     float64 `gorm:"column:scheduled_total"`
		PaidPrincipal      float64 `gorm:"column:paid_principal"`
		PaidInterest       float64 `gorm:"column:paid_interest"`
		PaidTotal          float64 `gorm:"column:paid_total"`
		UnpaidPrincipal    float64 `gorm:"column:unpaid_principal"`
		UnpaidInterest     float64 `gorm:"column:unpaid_interest"`
		UnpaidTotal        float64 `gorm:"column:unpaid_total"`
		ScheduleStatus     string  `gorm:"column:schedule_status"`
	}

	var rows []amortizationRow

	err := database.DB.Raw(`
		SELECT *
		FROM public.fn_report_loan_amortization(?)
	`,
		loanID,
	).Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	result := &dto.AmortizationReportData{
		Rows: make(
			[]dto.AmortizationReportRow,
			0,
			len(rows),
		),
	}

	for _, row := range rows {

		reportRow := dto.AmortizationReportRow{
			LoanID:             row.LoanID,
			PNNumber:           row.PNNumber,
			ClientName:         row.ClientName,
			ContactNumber:      row.ContactNumber,
			LoanType:           row.LoanType,
			PrincipalAmount:    row.PrincipalAmount,
			InterestRate:       row.InterestRate,
			LoanInterest:       row.LoanInterest,
			PNValue:            row.PNValue,
			LoanTerm:           row.LoanTerm,
			AmortizationAmount: row.AmortizationAmount,
			Frequency:          row.Frequency,
			DisbursementDate:   row.DisbursementDate,
			MaturityDate:       row.MaturityDate,
			LoanStatus:         row.LoanStatus,

			ScheduleID:         row.ScheduleID,
			ScheduleNo:         row.ScheduleNo,
			DueDate:            row.DueDate,
			ScheduledPrincipal: row.ScheduledPrincipal,
			ScheduledInterest:  row.ScheduledInterest,
			ScheduledTotal:     row.ScheduledTotal,
			PaidPrincipal:      row.PaidPrincipal,
			PaidInterest:       row.PaidInterest,
			PaidTotal:          row.PaidTotal,
			UnpaidPrincipal:    row.UnpaidPrincipal,
			UnpaidInterest:     row.UnpaidInterest,
			UnpaidTotal:        row.UnpaidTotal,
			ScheduleStatus:     row.ScheduleStatus,
		}

		result.Rows = append(
			result.Rows,
			reportRow,
		)

		// The loan information is repeated
		// on every returned schedule row.
		// Keep one copy for the PDF header.
		if result.Loan.LoanID == 0 {
			result.Loan = reportRow
		}

		result.Summary.TotalScheduledPrincipal +=
			row.ScheduledPrincipal

		result.Summary.TotalScheduledInterest +=
			row.ScheduledInterest

		result.Summary.TotalScheduled +=
			row.ScheduledTotal

		result.Summary.TotalPaidPrincipal +=
			row.PaidPrincipal

		result.Summary.TotalPaidInterest +=
			row.PaidInterest

		result.Summary.TotalPaid +=
			row.PaidTotal

		result.Summary.TotalUnpaidPrincipal +=
			row.UnpaidPrincipal

		result.Summary.TotalUnpaidInterest +=
			row.UnpaidInterest

		result.Summary.TotalUnpaid +=
			row.UnpaidTotal
	}

	return result, nil
}

func (r *ReportRepository) GetPARReport() (
	*dto.PARReportData,
	error,
) {

	type parRow struct {
		LoanID          int64   `gorm:"column:loan_id"`
		PNNumber        string  `gorm:"column:pn_number"`
		ClientName      string  `gorm:"column:client_name"`
		ContactNumber   string  `gorm:"column:contact_number"`
		LoanType        string  `gorm:"column:loan_type"`
		PrincipalAmount float64 `gorm:"column:principal_amount"`
		PNValue         float64 `gorm:"column:pn_value"`
		LoanStatus      string  `gorm:"column:loan_status"`
		OldestDueDate   string  `gorm:"column:oldest_due_date"`
		DaysPastDue     int     `gorm:"column:days_past_due"`
		DefaultAmount   float64 `gorm:"column:default_amount"`
		AgingBucket     string  `gorm:"column:aging_bucket"`
	}

	var rows []parRow

	err := database.DB.Raw(`
		SELECT *
		FROM public.fn_report_portfolio_at_risk()
	`).Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	result := &dto.PARReportData{
		Rows: make(
			[]dto.PARReportRow,
			0,
			len(rows),
		),
	}

	for _, row := range rows {

		result.Rows = append(
			result.Rows,
			dto.PARReportRow{
				LoanID:          row.LoanID,
				PNNumber:        row.PNNumber,
				ClientName:      row.ClientName,
				ContactNumber:   row.ContactNumber,
				LoanType:        row.LoanType,
				PrincipalAmount: row.PrincipalAmount,
				PNValue:         row.PNValue,
				LoanStatus:      row.LoanStatus,
				OldestDueDate:   row.OldestDueDate,
				DaysPastDue:     row.DaysPastDue,
				DefaultAmount:   row.DefaultAmount,
				AgingBucket:     row.AgingBucket,
			},
		)

		result.Summary.TotalLoans++

		result.Summary.DefaultAmount +=
			row.DefaultAmount

		switch row.AgingBucket {

		case "1-30":
			result.Summary.Aging1To30Loans++
			result.Summary.Aging1To30Amount +=
				row.DefaultAmount

		case "31-60":
			result.Summary.Aging31To60Loans++
			result.Summary.Aging31To60Amount +=
				row.DefaultAmount

		case "61-90":
			result.Summary.Aging61To90Loans++
			result.Summary.Aging61To90Amount +=
				row.DefaultAmount

		case "90+":
			result.Summary.Aging90PlusLoans++
			result.Summary.Aging90PlusAmount +=
				row.DefaultAmount
		}
	}

	return result, nil
}

// ========================================
// LOAN PORTFOLIO SUMMARY
// ========================================

func (r *ReportRepository) GetLoanPortfolioReport() (
	*dto.LoanPortfolioReportData,
	error,
) {

	type portfolioRow struct {
		LoanID            int64   `gorm:"column:loan_id"`
		PNNumber          string  `gorm:"column:pn_number"`
		ClientName        string  `gorm:"column:client_name"`
		LoanType          string  `gorm:"column:loan_type"`
		PrincipalAmount   float64 `gorm:"column:principal_amount"`
		InterestRate      float64 `gorm:"column:interest_rate"`
		LoanInterest      float64 `gorm:"column:loan_interest"`
		PNValue           float64 `gorm:"column:pn_value"`
		LoanTerm          int64   `gorm:"column:loan_term"`
		Frequency         string  `gorm:"column:frequency"`
		DisbursementDate  string  `gorm:"column:disbursement_date"`
		MaturityDate      string  `gorm:"column:maturity_date"`
		ScheduledAmount   float64 `gorm:"column:scheduled_amount"`
		PaidAmount        float64 `gorm:"column:paid_amount"`
		OutstandingAmount float64 `gorm:"column:outstanding_amount"`
		LoanStatus        string  `gorm:"column:loan_status"`
	}

	var rows []portfolioRow

	err := database.DB.Raw(`
		SELECT *
		FROM public.fn_report_loan_portfolio()
	`).Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	result := &dto.LoanPortfolioReportData{
		Rows: make(
			[]dto.LoanPortfolioReportRow,
			0,
			len(rows),
		),
		Summary: dto.LoanPortfolioReportSummary{
			StatusBreakdown: make(
				[]dto.LoanPortfolioStatusSummary,
				0,
			),
			TypeBreakdown: make(
				[]dto.LoanPortfolioTypeSummary,
				0,
			),
		},
	}

	// ========================================
	// BREAKDOWN MAPS
	// ========================================

	statusMap :=
		make(
			map[string]*dto.LoanPortfolioStatusSummary,
		)

	typeMap :=
		make(
			map[string]*dto.LoanPortfolioTypeSummary,
		)

	// ========================================
	// BUILD RESULT
	// ========================================

	for _, row := range rows {

		result.Rows = append(
			result.Rows,
			dto.LoanPortfolioReportRow{
				LoanID:            row.LoanID,
				PNNumber:          row.PNNumber,
				ClientName:        row.ClientName,
				LoanType:          row.LoanType,
				PrincipalAmount:   row.PrincipalAmount,
				InterestRate:      row.InterestRate,
				LoanInterest:      row.LoanInterest,
				PNValue:           row.PNValue,
				LoanTerm:          row.LoanTerm,
				Frequency:         row.Frequency,
				DisbursementDate:  row.DisbursementDate,
				MaturityDate:      row.MaturityDate,
				ScheduledAmount:   row.ScheduledAmount,
				PaidAmount:        row.PaidAmount,
				OutstandingAmount: row.OutstandingAmount,
				LoanStatus:        row.LoanStatus,
			},
		)

		// ========================================
		// TOTALS
		// ========================================

		result.Summary.TotalLoans++

		result.Summary.TotalPrincipal +=
			row.PrincipalAmount

		result.Summary.TotalPNValue +=
			row.PNValue

		result.Summary.TotalScheduled +=
			row.ScheduledAmount

		result.Summary.TotalPaid +=
			row.PaidAmount

		result.Summary.TotalOutstanding +=
			row.OutstandingAmount

		// ========================================
		// STATUS COUNTS
		// ========================================

		switch strings.ToUpper(
			row.LoanStatus,
		) {

		case "ACTIVE":
			result.Summary.ActiveLoans++

		case "COMPLETED":
			result.Summary.CompletedLoans++

		case "DEFAULT":
			result.Summary.DefaultedLoans++
		}

		// ========================================
		// STATUS BREAKDOWN
		// ========================================

		status := strings.TrimSpace(
			row.LoanStatus,
		)

		if status == "" {
			status = "UNKNOWN"
		}

		if _, exists := statusMap[status]; !exists {

			statusMap[status] =
				&dto.LoanPortfolioStatusSummary{
					Status: status,
				}
		}

		statusMap[status].Loans++

		statusMap[status].Principal +=
			row.PrincipalAmount

		statusMap[status].Outstanding +=
			row.OutstandingAmount

		// ========================================
		// LOAN TYPE BREAKDOWN
		// ========================================

		loanType := strings.TrimSpace(
			row.LoanType,
		)

		if loanType == "" {
			loanType = "UNSPECIFIED"
		}

		if _, exists := typeMap[loanType]; !exists {

			typeMap[loanType] =
				&dto.LoanPortfolioTypeSummary{
					LoanType: loanType,
				}
		}

		typeMap[loanType].Loans++

		typeMap[loanType].Principal +=
			row.PrincipalAmount

		typeMap[loanType].Outstanding +=
			row.OutstandingAmount
	}

	// ========================================
	// MAP → SLICE
	// ========================================

	for _, item := range statusMap {

		result.Summary.StatusBreakdown =
			append(
				result.Summary.StatusBreakdown,
				*item,
			)
	}

	for _, item := range typeMap {

		result.Summary.TypeBreakdown =
			append(
				result.Summary.TypeBreakdown,
				*item,
			)
	}

	// ========================================
	// SORT BREAKDOWNS
	// ========================================

	sort.Slice(
		result.Summary.StatusBreakdown,
		func(i, j int) bool {
			return result.Summary.StatusBreakdown[i].Status <
				result.Summary.StatusBreakdown[j].Status
		},
	)

	sort.Slice(
		result.Summary.TypeBreakdown,
		func(i, j int) bool {
			return result.Summary.TypeBreakdown[i].LoanType <
				result.Summary.TypeBreakdown[j].LoanType
		},
	)

	return result, nil
}

func (r *ReportRepository) GetLoanMaturityReport() (
	*dto.LoanMaturityReportData,
	error,
) {

	type maturityRow struct {
		LoanID        int64  `gorm:"column:loan_id"`
		PNNumber      string `gorm:"column:pn_number"`
		ClientName    string `gorm:"column:client_name"`
		ContactNumber string `gorm:"column:contact_number"`
		LoanType      string `gorm:"column:loan_type"`

		PrincipalAmount float64 `gorm:"column:principal_amount"`
		PNValue         float64 `gorm:"column:pn_value"`

		DisbursementDate string `gorm:"column:disbursement_date"`
		MaturityDate     string `gorm:"column:maturity_date"`

		DaysUntilMaturity int `gorm:"column:days_until_maturity"`

		ScheduledAmount   float64 `gorm:"column:scheduled_amount"`
		PaidAmount        float64 `gorm:"column:paid_amount"`
		OutstandingAmount float64 `gorm:"column:outstanding_amount"`

		LoanStatus     string `gorm:"column:loan_status"`
		MaturityStatus string `gorm:"column:maturity_status"`
	}

	var rows []maturityRow

	err := database.DB.Raw(`
		SELECT *
		FROM public.fn_report_loan_maturity()
	`).Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	result := &dto.LoanMaturityReportData{
		Rows: make(
			[]dto.LoanMaturityReportRow,
			0,
			len(rows),
		),
		Summary: dto.LoanMaturityReportSummary{
			BucketBreakdown: make(
				[]dto.LoanMaturityBucketSummary,
				0,
			),
		},
	}

	// ========================================
	// BUCKET MAP
	// ========================================

	bucketMap :=
		make(
			map[string]*dto.LoanMaturityBucketSummary,
		)

	// ========================================
	// PROCESS ROWS
	// ========================================

	for _, row := range rows {

		result.Rows = append(
			result.Rows,
			dto.LoanMaturityReportRow{
				LoanID:        row.LoanID,
				PNNumber:      row.PNNumber,
				ClientName:    row.ClientName,
				ContactNumber: row.ContactNumber,
				LoanType:      row.LoanType,

				PrincipalAmount: row.PrincipalAmount,
				PNValue:         row.PNValue,

				DisbursementDate: row.DisbursementDate,
				MaturityDate:     row.MaturityDate,

				DaysUntilMaturity: row.DaysUntilMaturity,

				ScheduledAmount:   row.ScheduledAmount,
				PaidAmount:        row.PaidAmount,
				OutstandingAmount: row.OutstandingAmount,

				LoanStatus:     row.LoanStatus,
				MaturityStatus: row.MaturityStatus,
			},
		)

		// ========================================
		// TOTALS
		// ========================================

		result.Summary.TotalLoans++

		result.Summary.TotalOutstanding +=
			row.OutstandingAmount

		// ========================================
		// MATURITY STATUS
		// ========================================

		switch row.MaturityStatus {

		case "MATURED":

			result.Summary.MaturedLoans++

			result.Summary.MaturedAmount +=
				row.OutstandingAmount

		case "DUE TODAY":

			result.Summary.DueTodayLoans++

			result.Summary.DueTodayAmount +=
				row.OutstandingAmount

		case "1-30 DAYS",
			"31-60 DAYS",
			"61-90 DAYS",
			"90+ DAYS":

			result.Summary.UpcomingLoans++

			result.Summary.UpcomingAmount +=
				row.OutstandingAmount
		}

		// ========================================
		// BUCKET BREAKDOWN
		// ========================================

		bucket :=
			strings.TrimSpace(
				row.MaturityStatus,
			)

		if bucket == "" {
			bucket = "UNKNOWN"
		}

		if _, exists :=
			bucketMap[bucket]; !exists {

			bucketMap[bucket] =
				&dto.LoanMaturityBucketSummary{
					Bucket: bucket,
				}
		}

		bucketMap[bucket].Loans++

		bucketMap[bucket].Outstanding +=
			row.OutstandingAmount
	}

	// ========================================
	// MAP → SLICE
	// ========================================

	for _, item := range bucketMap {

		result.Summary.BucketBreakdown =
			append(
				result.Summary.BucketBreakdown,
				*item,
			)
	}

	// ========================================
	// SORT BUCKETS
	// ========================================

	bucketOrder := map[string]int{
		"MATURED":          1,
		"DUE TODAY":        2,
		"1-30 DAYS":        3,
		"31-60 DAYS":       4,
		"61-90 DAYS":       5,
		"90+ DAYS":         6,
		"NO MATURITY DATE": 7,
		"UNKNOWN":          8,
	}

	sort.Slice(
		result.Summary.BucketBreakdown,
		func(i, j int) bool {

			return bucketOrder[result.Summary.BucketBreakdown[i].Bucket] <
				bucketOrder[result.Summary.BucketBreakdown[j].Bucket]
		},
	)

	return result, nil
}
