package repositories

import (
	"github.com/jj.jobo/FGC/internal/database"
	"github.com/jj.jobo/FGC/internal/dto"
)

type DashboardRepository struct {
	*BaseRepository[struct{}]
}

func NewDashboardRepository() *DashboardRepository {
	return &DashboardRepository{
		BaseRepository: NewBaseRepository[struct{}](
			database.DB,
		),
	}
}

func (r *DashboardRepository) GetSummary() (
	*dto.DashboardSummary,
	error,
) {
	var summary dto.DashboardSummary

	err := database.DB.Raw(`
		SELECT
			(
				SELECT COUNT(*)
				FROM clients
			) AS total_clients,

			(
				SELECT COUNT(*)
				FROM loans
			) AS total_loans,

			(
				SELECT COUNT(*)
				FROM loans
				WHERE status = 'ACTIVE'
			) AS active_loans,

			(
				SELECT COUNT(*)
				FROM loans
				WHERE status = 'PAID'
			) AS paid_loans,

			COALESCE(
				(
					SELECT SUM(principal_amount)
					FROM loans
				),
				0
			) AS total_principal,

			COALESCE(
				(
					SELECT SUM(pn_value)
					FROM loans
				),
				0
			) AS total_pn_value,

			COALESCE(
				(
					SELECT SUM(amount_paid)
					FROM payments
				),
				0
			) AS total_collected,

			COALESCE(
				(
					SELECT SUM(
						GREATEST(
							a.principal_amount -
							COALESCE(
								a.paid_principal_amount,
								0
							),
							0
						)
						+
						GREATEST(
							a.interest_amount -
							COALESCE(
								a.paid_interest_amount,
								0
							),
							0
						)
					)
					FROM amortization_schedules a
				),
				0
			) AS total_outstanding
	`).Scan(&summary).Error

	if err != nil {
		return nil, err
	}

	return &summary, nil
}
