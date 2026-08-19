package dto

type ClientProfileResponse struct {
	Client   ClientProfileClient `json:"client"`
	Summary  ClientSummary       `json:"summary"`
	Loans    []ClientLoanSummary `json:"loans"`
	Payments []ClientPayment     `json:"payments"`
}

type ClientProfileClient struct {
	ID             int     `json:"id"`
	FirstName      string  `json:"first_name"`
	LastName       string  `json:"last_name"`
	ContactNumber  *string `json:"contact_number,omitempty"`
	Email          *string `json:"email,omitempty"`
	CurrentAddress *string `json:"current_address,omitempty"`
}

type ClientSummary struct {
	TotalLoans       int     `json:"total_loans"`
	ActiveLoans      int     `json:"active_loans"`
	TotalPrincipal   float64 `json:"total_principal"`
	TotalPaid        float64 `json:"total_paid"`
	TotalOutstanding float64 `json:"total_outstanding"`
}

type ClientLoanSummary struct {
	ID                 int     `json:"id"`
	PNNumber           string  `json:"pn_number"`
	LoanType           *string `json:"loan_type,omitempty"`
	PrincipalAmount    float64 `json:"principal_amount"`
	InterestRate       float64 `json:"interest_rate"`
	LoanInterest       float64 `json:"loan_interest"`
	PNValue            float64 `json:"pn_value"`
	LoanTerm           int     `json:"loan_term"`
	AmortizationAmount float64 `json:"amortization_amount"`
	DisbursementDate   *string `json:"disbursement_date,omitempty"`
	MaturityDate       *string `json:"maturity_date,omitempty"`
	Status             string  `json:"status"`
	TotalPaid          float64 `json:"total_paid"`
	OutstandingBalance float64 `json:"outstanding_balance"`
}

type ClientPayment struct {
	ID                 int     `json:"id"`
	LoanID             int     `json:"loan_id"`
	PaymentDate        string  `json:"payment_date"`
	AmountPaid         float64 `json:"amount_paid"`
	PaymentChannel     *string `json:"payment_channel,omitempty"`
	ReferenceNumber    *string `json:"reference_number,omitempty"`
	PrincipalApplied   float64 `json:"principal_applied"`
	InterestApplied    float64 `json:"interest_applied"`
	OutstandingBalance float64 `json:"outstanding_balance"`
}
