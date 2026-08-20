package dto

type PARLoanResponse struct {
	ID            int64   `json:"id"`
	PNNumber      string  `json:"pn_number"`
	ClientName    string  `json:"client_name"`
	DueDate       string  `json:"due_date"`
	DaysPastDue   int     `json:"days_past_due"`
	DefaultAmount float64 `json:"default_amount"`
	Status        string  `json:"status"`
}

type PARSummaryResponse struct {
	PARLoans      int     `json:"par_loans"`
	DefaultAmount float64 `json:"default_amount"`
	PARRatio      float64 `json:"par_ratio"`
}

type PARPaginationResponse struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type PARResponse struct {
	Summary    PARSummaryResponse    `json:"summary"`
	Loans      []PARLoanResponse     `json:"loans"`
	Pagination PARPaginationResponse `json:"pagination"`
	Aging      []PARAgingResponse    `json:"aging"`
}

type PARAgingResponse struct {
	Aging         string  `json:"aging"`
	Loans         int     `json:"loans"`
	DefaultAmount float64 `json:"default_amount"`
}
