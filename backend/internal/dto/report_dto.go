package dto

type ReportType string

const (
	ReportSummaryCollection ReportType = "SUMMARY_COLLECTION"
	ReportAmortization      ReportType = "LOAN_AMORTIZATION"
	ReportPAR               ReportType = "AT_RISK_LOANS"
	ReportAging             ReportType = "LOAN_AGING"
	ReportSummaryRelease    ReportType = "SUMMARY_RELEASE"
	ReportPaymentJournal    ReportType = "PAYMENT_JOURNAL"
)

type ReportRequest struct {
	ReportType ReportType `json:"report_type"`
	DateFrom   string     `json:"date_from"`
	DateTo     string     `json:"date_to"`
	LoanID     *int64     `json:"loan_id,omitempty"`
}

// ========================================
// SUMMARY OF COLLECTION
// ========================================

type CollectionReportRow struct {
	PaymentDate        string  `json:"payment_date"`
	ReferenceNumber    string  `json:"reference_number"`
	PNNumber           string  `json:"pn_number"`
	ClientName         string  `json:"client_name"`
	PrincipalApplied   float64 `json:"principal_applied"`
	InterestApplied    float64 `json:"interest_applied"`
	AmountPaid         float64 `json:"amount_paid"`
	PaymentChannel     string  `json:"payment_channel"`
	OutstandingBalance float64 `json:"outstanding_balance"`
}

type CollectionReportSummary struct {
	TotalTransactions int     `json:"total_transactions"`
	TotalPrincipal    float64 `json:"total_principal"`
	TotalInterest     float64 `json:"total_interest"`
	TotalCollection   float64 `json:"total_collection"`
	TotalOutstanding  float64 `json:"total_outstanding"`
}

type CollectionReportData struct {
	Summary CollectionReportSummary `json:"summary"`
	Rows    []CollectionReportRow   `json:"rows"`
}

// ========================================
// LOAN ACCOUNT AMORTIZATION
// ========================================

type AmortizationReportRow struct {
	LoanID             int64   `json:"loan_id"`
	PNNumber           string  `json:"pn_number"`
	ClientName         string  `json:"client_name"`
	ContactNumber      string  `json:"contact_number"`
	LoanType           string  `json:"loan_type"`
	PrincipalAmount    float64 `json:"principal_amount"`
	InterestRate       float64 `json:"interest_rate"`
	LoanInterest       float64 `json:"loan_interest"`
	PNValue            float64 `json:"pn_value"`
	LoanTerm           int64   `json:"loan_term"`
	AmortizationAmount float64 `json:"amortization_amount"`
	Frequency          string  `json:"frequency"`
	DisbursementDate   string  `json:"disbursement_date"`
	MaturityDate       string  `json:"maturity_date"`
	LoanStatus         string  `json:"loan_status"`

	ScheduleID         int     `json:"schedule_id"`
	ScheduleNo         int64   `json:"schedule_no"`
	DueDate            string  `json:"due_date"`
	ScheduledPrincipal float64 `json:"scheduled_principal"`
	ScheduledInterest  float64 `json:"scheduled_interest"`
	ScheduledTotal     float64 `json:"scheduled_total"`
	PaidPrincipal      float64 `json:"paid_principal"`
	PaidInterest       float64 `json:"paid_interest"`
	PaidTotal          float64 `json:"paid_total"`
	UnpaidPrincipal    float64 `json:"unpaid_principal"`
	UnpaidInterest     float64 `json:"unpaid_interest"`
	UnpaidTotal        float64 `json:"unpaid_total"`
	ScheduleStatus     string  `json:"schedule_status"`
}

type AmortizationReportSummary struct {
	TotalScheduledPrincipal float64 `json:"total_scheduled_principal"`
	TotalScheduledInterest  float64 `json:"total_scheduled_interest"`
	TotalScheduled          float64 `json:"total_scheduled"`

	TotalPaidPrincipal float64 `json:"total_paid_principal"`
	TotalPaidInterest  float64 `json:"total_paid_interest"`
	TotalPaid          float64 `json:"total_paid"`

	TotalUnpaidPrincipal float64 `json:"total_unpaid_principal"`
	TotalUnpaidInterest  float64 `json:"total_unpaid_interest"`
	TotalUnpaid          float64 `json:"total_unpaid"`
}

type AmortizationReportData struct {
	Loan    AmortizationReportRow     `json:"loan"`
	Summary AmortizationReportSummary `json:"summary"`
	Rows    []AmortizationReportRow   `json:"rows"`
}

type PARReportRow struct {
	LoanID          int64   `json:"loan_id"`
	PNNumber        string  `json:"pn_number"`
	ClientName      string  `json:"client_name"`
	ContactNumber   string  `json:"contact_number"`
	LoanType        string  `json:"loan_type"`
	PrincipalAmount float64 `json:"principal_amount"`
	PNValue         float64 `json:"pn_value"`
	LoanStatus      string  `json:"loan_status"`

	OldestDueDate string  `json:"oldest_due_date"`
	DaysPastDue   int     `json:"days_past_due"`
	DefaultAmount float64 `json:"default_amount"`
	AgingBucket   string  `json:"aging_bucket"`
}

type PARReportSummary struct {
	TotalLoans    int     `json:"total_loans"`
	DefaultAmount float64 `json:"default_amount"`

	Aging1To30Loans  int     `json:"aging_1_30_loans"`
	Aging1To30Amount float64 `json:"aging_1_30_amount"`

	Aging31To60Loans  int     `json:"aging_31_60_loans"`
	Aging31To60Amount float64 `json:"aging_31_60_amount"`

	Aging61To90Loans  int     `json:"aging_61_90_loans"`
	Aging61To90Amount float64 `json:"aging_61_90_amount"`

	Aging90PlusLoans  int     `json:"aging_90_plus_loans"`
	Aging90PlusAmount float64 `json:"aging_90_plus_amount"`
}

type PARReportData struct {
	Rows    []PARReportRow   `json:"rows"`
	Summary PARReportSummary `json:"summary"`
}

// ========================================
// LOAN PORTFOLIO SUMMARY
// ========================================

type LoanPortfolioReportRow struct {
	LoanID            int64   `json:"loan_id"`
	PNNumber          string  `json:"pn_number"`
	ClientName        string  `json:"client_name"`
	LoanType          string  `json:"loan_type"`
	PrincipalAmount   float64 `json:"principal_amount"`
	InterestRate      float64 `json:"interest_rate"`
	LoanInterest      float64 `json:"loan_interest"`
	PNValue           float64 `json:"pn_value"`
	LoanTerm          int64   `json:"loan_term"`
	Frequency         string  `json:"frequency"`
	DisbursementDate  string  `json:"disbursement_date"`
	MaturityDate      string  `json:"maturity_date"`
	ScheduledAmount   float64 `json:"scheduled_amount"`
	PaidAmount        float64 `json:"paid_amount"`
	OutstandingAmount float64 `json:"outstanding_amount"`
	LoanStatus        string  `json:"loan_status"`
}

type LoanPortfolioStatusSummary struct {
	Status      string  `json:"status"`
	Loans       int     `json:"loans"`
	Principal   float64 `json:"principal"`
	Outstanding float64 `json:"outstanding"`
}

type LoanPortfolioTypeSummary struct {
	LoanType    string  `json:"loan_type"`
	Loans       int     `json:"loans"`
	Principal   float64 `json:"principal"`
	Outstanding float64 `json:"outstanding"`
}

type LoanPortfolioReportSummary struct {
	TotalLoans       int     `json:"total_loans"`
	TotalPrincipal   float64 `json:"total_principal"`
	TotalPNValue     float64 `json:"total_pn_value"`
	TotalScheduled   float64 `json:"total_scheduled"`
	TotalPaid        float64 `json:"total_paid"`
	TotalOutstanding float64 `json:"total_outstanding"`

	ActiveLoans    int `json:"active_loans"`
	CompletedLoans int `json:"completed_loans"`
	DefaultedLoans int `json:"defaulted_loans"`

	StatusBreakdown []LoanPortfolioStatusSummary `json:"status_breakdown"`
	TypeBreakdown   []LoanPortfolioTypeSummary   `json:"type_breakdown"`
}

type LoanPortfolioReportData struct {
	Summary LoanPortfolioReportSummary `json:"summary"`
	Rows    []LoanPortfolioReportRow   `json:"rows"`
}

// ========================================
// LOAN MATURITY / DUE REPORT
// ========================================

type LoanMaturityReportRow struct {
	LoanID        int64  `json:"loan_id"`
	PNNumber      string `json:"pn_number"`
	ClientName    string `json:"client_name"`
	ContactNumber string `json:"contact_number"`
	LoanType      string `json:"loan_type"`

	PrincipalAmount float64 `json:"principal_amount"`
	PNValue         float64 `json:"pn_value"`

	DisbursementDate string `json:"disbursement_date"`
	MaturityDate     string `json:"maturity_date"`

	DaysUntilMaturity int `json:"days_until_maturity"`

	ScheduledAmount   float64 `json:"scheduled_amount"`
	PaidAmount        float64 `json:"paid_amount"`
	OutstandingAmount float64 `json:"outstanding_amount"`

	LoanStatus     string `json:"loan_status"`
	MaturityStatus string `json:"maturity_status"`
}

type LoanMaturityBucketSummary struct {
	Bucket      string  `json:"bucket"`
	Loans       int     `json:"loans"`
	Outstanding float64 `json:"outstanding"`
}

type LoanMaturityReportSummary struct {
	TotalLoans       int     `json:"total_loans"`
	TotalOutstanding float64 `json:"total_outstanding"`

	MaturedLoans  int     `json:"matured_loans"`
	MaturedAmount float64 `json:"matured_amount"`

	DueTodayLoans  int     `json:"due_today_loans"`
	DueTodayAmount float64 `json:"due_today_amount"`

	UpcomingLoans  int     `json:"upcoming_loans"`
	UpcomingAmount float64 `json:"upcoming_amount"`

	BucketBreakdown []LoanMaturityBucketSummary `json:"bucket_breakdown"`
}

type LoanMaturityReportData struct {
	Summary LoanMaturityReportSummary `json:"summary"`
	Rows    []LoanMaturityReportRow   `json:"rows"`
}
