package dto

type DashboardSummary struct {
	TotalClients     int64   `json:"total_clients"`
	TotalLoans       int64   `json:"total_loans"`
	ActiveLoans      int64   `json:"active_loans"`
	PaidLoans        int64   `json:"paid_loans"`
	TotalPrincipal   float64 `json:"total_principal"`
	TotalPNValue     float64 `json:"total_pn_value"`
	TotalCollected   float64 `json:"total_collected"`
	TotalOutstanding float64 `json:"total_outstanding"`
}
