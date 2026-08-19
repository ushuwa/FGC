package dto

type LoanListItem struct {
	ID                 int     `json:"id"`
	ClientID           int     `json:"client_id"`
	ClientName         string  `json:"client_name"`
	PNNumber           string  `json:"pn_number"`
	LoanType           *string `json:"loan_type,omitempty"`
	PrincipalAmount    float64 `json:"principal_amount"`
	InterestRate       float64 `json:"interest_rate"`
	LoanInterest       float64 `json:"loan_interest"`
	PNValue            float64 `json:"pn_value"`
	LoanTerm           int     `json:"loan_term"`
	Frequency          string  `json:"frequency"`
	AmortizationAmount float64 `json:"amortization_amount"`
	DisbursementDate   *string `json:"disbursement_date,omitempty"`
	MaturityDate       *string `json:"maturity_date,omitempty"`
	Status             string  `json:"status"`
	TotalPaid          float64 `json:"total_paid"`
	OutstandingBalance float64 `json:"outstanding_balance"`
}

type LoanProfileClient struct {
	ID             int     `json:"id"`
	FirstName      string  `json:"first_name"`
	LastName       string  `json:"last_name"`
	ContactNumber  *string `json:"contact_number,omitempty"`
	Email          *string `json:"email,omitempty"`
	CurrentAddress *string `json:"current_address,omitempty"`
}

type LoanProfileInfo struct {
	ID                 int     `json:"id"`
	ClientID           int     `json:"client_id"`
	PNNumber           string  `json:"pn_number"`
	LoanType           *string `json:"loan_type,omitempty"`
	PrincipalAmount    float64 `json:"principal_amount"`
	InterestRate       float64 `json:"interest_rate"`
	LoanInterest       float64 `json:"loan_interest"`
	PNValue            float64 `json:"pn_value"`
	LoanTerm           int     `json:"loan_term"`
	Frequency          string  `json:"frequency"`
	AmortizationAmount float64 `json:"amortization_amount"`
	DisbursementDate   *string `json:"disbursement_date,omitempty"`
	MaturityDate       *string `json:"maturity_date,omitempty"`
	Status             string  `json:"status"`
}

type LoanProfileSummary struct {
	PrincipalAmount    float64 `json:"principal_amount"`
	PNValue            float64 `json:"pn_value"`
	TotalPaid          float64 `json:"total_paid"`
	OutstandingBalance float64 `json:"outstanding_balance"`
}

type LoanAmortization struct {
	ID                  int     `json:"id"`
	LoanID              int     `json:"loan_id"`
	DueDate             string  `json:"due_date"`
	PrincipalAmount     float64 `json:"principal_amount"`
	InterestAmount      float64 `json:"interest_amount"`
	TotalAmount         float64 `json:"total_amount"`
	PaidPrincipalAmount float64 `json:"paid_principal_amount"`
	PaidInterestAmount  float64 `json:"paid_interest_amount"`
	Status              string  `json:"status"`
}

type LoanPayment struct {
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

type LoanProfileResponse struct {
	Loan          LoanProfileInfo    `json:"loan"`
	Client        LoanProfileClient  `json:"client"`
	Summary       LoanProfileSummary `json:"summary"`
	Amortizations []LoanAmortization `json:"amortizations"`
	Payments      []LoanPayment      `json:"payments"`
}

type CreatePaymentRequest struct {
	PaymentDate     string  `json:"payment_date"`
	AmountPaid      float64 `json:"amount_paid"`
	PaymentChannel  *string `json:"payment_channel"`
	ReferenceNumber *string `json:"reference_number"`
}

type PaymentResponse struct {
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

type CreateLoanRequest struct {
	ClientID           int     `json:"client_id"`
	PNNumber           string  `json:"pn_number"`
	LoanType           *string `json:"loan_type"`
	PrincipalAmount    float64 `json:"principal_amount"`
	InterestRate       float64 `json:"interest_rate"`
	LoanInterest       float64 `json:"loan_interest"`
	PNValue            float64 `json:"pn_value"`
	LoanTerm           int     `json:"loan_term"`
	AmortizationAmount float64 `json:"amortization_amount"`
	DisbursementDate   string  `json:"disbursement_date"`
	MaturityDate       string  `json:"maturity_date"`
	Frequency          string  `json:"frequency"`
	Status             string  `json:"status"`
}

type UpdateLoanRequest struct {
	ClientID           int     `json:"client_id"`
	PNNumber           string  `json:"pn_number"`
	LoanType           *string `json:"loan_type"`
	PrincipalAmount    float64 `json:"principal_amount"`
	InterestRate       float64 `json:"interest_rate"`
	LoanInterest       float64 `json:"loan_interest"`
	PNValue            float64 `json:"pn_value"`
	LoanTerm           int     `json:"loan_term"`
	AmortizationAmount float64 `json:"amortization_amount"`
	DisbursementDate   string  `json:"disbursement_date"`
	MaturityDate       string  `json:"maturity_date"`
	Frequency          string  `json:"frequency"`
	Status             string  `json:"status"`
}
