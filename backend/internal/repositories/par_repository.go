package repositories

import (
	"math"

	"github.com/jj.jobo/FGC/internal/database"
	"github.com/jj.jobo/FGC/internal/dto"
)

type PARRepository struct {
	*BaseRepository[struct{}]
}

func NewPARRepository() *PARRepository {
	return &PARRepository{
		BaseRepository: NewBaseRepository[struct{}](
			database.DB,
		),
	}
}

// GetPAR returns loans that currently have
// unpaid scheduled amounts past their due date.
func (r *PARRepository) GetPAR(
	search string,
	status string,
	aging string,
	page int,
	limit int,
) (*dto.PARResponse, error) {

	var result dto.PARResponse

	// ========================================
	// PAGINATION DEFAULTS
	// ========================================

	if page < 1 {
		page = 1
	}

	if limit <= 0 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	// ========================================
	// 1. SUMMARY
	// ========================================

	err := database.DB.Raw(`
		WITH past_due AS (
			SELECT
				a.loan_id,

				SUM(
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
				) AS default_amount

			FROM amortization_schedules a

			WHERE
				a.due_date < CURRENT_DATE

			GROUP BY
				a.loan_id

			HAVING
				SUM(
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
				) > 0
		),

		total_portfolio AS (
			SELECT
				COALESCE(
					SUM(
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
					),
					0
				) AS outstanding

			FROM amortization_schedules a
		),

		par_totals AS (
			SELECT
				COUNT(*) AS par_loans,

				COALESCE(
					SUM(default_amount),
					0
				) AS default_amount

			FROM past_due
		)

		SELECT
			par_totals.par_loans,
			par_totals.default_amount,

			CASE
				WHEN total_portfolio.outstanding > 0
				THEN
					(
						par_totals.default_amount
						/
						total_portfolio.outstanding
					) * 100
				ELSE 0
			END AS par_ratio

		FROM par_totals

		CROSS JOIN total_portfolio
	`).Scan(&result.Summary).Error

	if err != nil {
		return nil, err
	}

	// ========================================
	// 2. AGING SUMMARY
	// ========================================

	type agingRow struct {
		Aging         string  `gorm:"column:aging"`
		Loans         int     `gorm:"column:loans"`
		DefaultAmount float64 `gorm:"column:default_amount"`
	}

	var agingRows []agingRow

	err = database.DB.Raw(`
  		  WITH past_due AS (
        SELECT
            a.loan_id,

            MIN(
                a.due_date
            ) FILTER (
                WHERE
                    a.due_date < CURRENT_DATE
                    AND (
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
                    ) > 0
            ) AS oldest_due_date,

            SUM(
                CASE
                    WHEN a.due_date < CURRENT_DATE
                    THEN
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
                    ELSE 0
                END
            ) AS default_amount

        FROM amortization_schedules a

        GROUP BY
            a.loan_id

        HAVING
            SUM(
                CASE
                    WHEN a.due_date < CURRENT_DATE
                    THEN
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
                    ELSE 0
                END
            ) > 0
    ),

    filtered AS (
        SELECT
            p.loan_id,
            p.oldest_due_date,
            p.default_amount

        FROM past_due p

        INNER JOIN loans l
            ON l.id = p.loan_id

        LEFT JOIN clients c
            ON c.id = l.client_id

        WHERE
            (
                ? = ''
                OR l.pn_number ILIKE '%' || ? || '%'
                OR c.first_name ILIKE '%' || ? || '%'
                OR c.last_name ILIKE '%' || ? || '%'
                OR (
                    COALESCE(c.first_name, '')
                    || ' '
                    ||
                    COALESCE(c.last_name, '')
                ) ILIKE '%' || ? || '%'
            )

            AND (
                ? = ''
                OR 'PAR' = ?
            )
    )

    SELECT
        CASE
            WHEN (
                CURRENT_DATE -
                oldest_due_date
            ) BETWEEN 1 AND 30
                THEN '1-30'

            WHEN (
                CURRENT_DATE -
                oldest_due_date
            ) BETWEEN 31 AND 60
                THEN '31-60'

            WHEN (
                CURRENT_DATE -
                oldest_due_date
            ) BETWEEN 61 AND 90
                THEN '61-90'

            ELSE '90+'
        END AS aging,

        COUNT(*) AS loans,

        COALESCE(
            SUM(default_amount),
            0
        ) AS default_amount

    FROM filtered

    GROUP BY
        CASE
            WHEN (
                CURRENT_DATE -
                oldest_due_date
            ) BETWEEN 1 AND 30
                THEN '1-30'

            WHEN (
                CURRENT_DATE -
                oldest_due_date
            ) BETWEEN 31 AND 60
                THEN '31-60'

            WHEN (
                CURRENT_DATE -
                oldest_due_date
            ) BETWEEN 61 AND 90
                THEN '61-90'

            ELSE '90+'
        END

    ORDER BY
        MIN(
            CURRENT_DATE -
            oldest_due_date
        )
`,
		search,
		search,
		search,
		search,
		search,
		status,
		status,
	).Scan(&agingRows).Error

	if err != nil {
		return nil, err
	}

	result.Aging = make(
		[]dto.PARAgingResponse,
		0,
		len(agingRows),
	)

	for _, row := range agingRows {

		result.Aging = append(
			result.Aging,
			dto.PARAgingResponse{
				Aging:         row.Aging,
				Loans:         row.Loans,
				DefaultAmount: row.DefaultAmount,
			},
		)
	}

	// ========================================
	// 3. TOTAL FILTERED PAR LOANS
	// ========================================

	var total int64

	err = database.DB.Raw(`
		WITH past_due AS (
			SELECT
				a.loan_id,

				MIN(
					a.due_date
				) FILTER (
					WHERE
						a.due_date < CURRENT_DATE
						AND (
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
						) > 0
				) AS oldest_due_date,

				SUM(
					CASE
						WHEN
							a.due_date < CURRENT_DATE
						THEN
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
						ELSE 0
					END
				) AS default_amount

			FROM amortization_schedules a

			GROUP BY
				a.loan_id

			HAVING
				SUM(
					CASE
						WHEN
							a.due_date < CURRENT_DATE
						THEN
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
						ELSE 0
					END
				) > 0
		)

		SELECT
			COUNT(*)

		FROM past_due p

		INNER JOIN loans l
			ON l.id = p.loan_id

		LEFT JOIN clients c
			ON c.id = l.client_id

		WHERE

			(
				? = ''
				OR l.pn_number ILIKE '%' || ? || '%'
				OR c.first_name ILIKE '%' || ? || '%'
				OR c.last_name ILIKE '%' || ? || '%'
				OR (
					COALESCE(
						c.first_name,
						''
					)
					||
					' '
					||
					COALESCE(
						c.last_name,
						''
					)
				) ILIKE '%' || ? || '%'
			)

			AND

			(
				? = ''
				OR 'PAR' = ?
			)

			AND

			(
				? = ''

				OR (
					? = '1-30'
					AND (
						CURRENT_DATE -
						p.oldest_due_date
					) BETWEEN 1 AND 30
				)

				OR (
					? = '31-60'
					AND (
						CURRENT_DATE -
						p.oldest_due_date
					) BETWEEN 31 AND 60
				)

				OR (
					? = '61-90'
					AND (
						CURRENT_DATE -
						p.oldest_due_date
					) BETWEEN 61 AND 90
				)

				OR (
					? = '90+'
					AND (
						CURRENT_DATE -
						p.oldest_due_date
					) > 90
				)
			)
	`,
		search,
		search,
		search,
		search,
		search,
		status,
		status,
		aging,
		aging,
		aging,
		aging,
		aging,
	).Scan(&total).Error

	if err != nil {
		return nil, err
	}

	// ========================================
	// 3. PAR LOANS
	// ========================================

	type parRow struct {
		ID            int     `gorm:"column:id"`
		PNNumber      string  `gorm:"column:pn_number"`
		ClientName    string  `gorm:"column:client_name"`
		DueDate       string  `gorm:"column:due_date"`
		DaysPastDue   int     `gorm:"column:days_past_due"`
		DefaultAmount float64 `gorm:"column:default_amount"`
		Status        string  `gorm:"column:status"`
	}

	var rows []parRow

	err = database.DB.Raw(`
		WITH past_due AS (
			SELECT
				a.loan_id,

				MIN(
					a.due_date
				) FILTER (
					WHERE
						a.due_date < CURRENT_DATE
						AND (
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
						) > 0
				) AS oldest_due_date,

				SUM(
					CASE
						WHEN
							a.due_date < CURRENT_DATE
						THEN
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
						ELSE 0
					END
				) AS default_amount

			FROM amortization_schedules a

			GROUP BY
				a.loan_id

			HAVING
				SUM(
					CASE
						WHEN
							a.due_date < CURRENT_DATE
						THEN
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
						ELSE 0
					END
				) > 0
		)

		SELECT
			l.id,

			l.pn_number,

			TRIM(
				COALESCE(
					c.first_name,
					''
				)
				||
				' '
				||
				COALESCE(
					c.last_name,
					''
				)
			) AS client_name,

			TO_CHAR(
				p.oldest_due_date,
				'YYYY-MM-DD'
			) AS due_date,

			(
				CURRENT_DATE -
				p.oldest_due_date
			)::integer AS days_past_due,

			p.default_amount,

			'PAR' AS status

		FROM past_due p

		INNER JOIN loans l
			ON l.id = p.loan_id

		LEFT JOIN clients c
			ON c.id = l.client_id

		WHERE

			(
				? = ''
				OR l.pn_number ILIKE '%' || ? || '%'
				OR c.first_name ILIKE '%' || ? || '%'
				OR c.last_name ILIKE '%' || ? || '%'
				OR (
					COALESCE(
						c.first_name,
						''
					)
					||
					' '
					||
					COALESCE(
						c.last_name,
						''
					)
				) ILIKE '%' || ? || '%'
			)

			AND

			(
				? = ''
				OR 'PAR' = ?
			)

			AND

			(
				? = ''

				OR (
					? = '1-30'
					AND (
						CURRENT_DATE -
						p.oldest_due_date
					) BETWEEN 1 AND 30
				)

				OR (
					? = '31-60'
					AND (
						CURRENT_DATE -
						p.oldest_due_date
					) BETWEEN 31 AND 60
				)

				OR (
					? = '61-90'
					AND (
						CURRENT_DATE -
						p.oldest_due_date
					) BETWEEN 61 AND 90
				)

				OR (
					? = '90+'
					AND (
						CURRENT_DATE -
						p.oldest_due_date
					) > 90
				)
			)

		ORDER BY
			p.oldest_due_date ASC,
			l.id ASC

		LIMIT ?

		OFFSET ?
	`,
		search,
		search,
		search,
		search,
		search,
		status,
		status,
		aging,
		aging,
		aging,
		aging,
		aging,
		limit,
		offset,
	).Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	// ========================================
	// 4. MAP LOANS
	// ========================================

	result.Loans = make(
		[]dto.PARLoanResponse,
		0,
		len(rows),
	)

	for _, row := range rows {

		result.Loans = append(
			result.Loans,
			dto.PARLoanResponse{
				ID:            int64(row.ID),
				PNNumber:      row.PNNumber,
				ClientName:    row.ClientName,
				DueDate:       row.DueDate,
				DaysPastDue:   row.DaysPastDue,
				DefaultAmount: row.DefaultAmount,
				Status:        row.Status,
			},
		)
	}

	// ========================================
	// 5. PAGINATION RESPONSE
	// ========================================

	totalPages := 0

	if total > 0 {
		totalPages = int(
			math.Ceil(
				float64(total) /
					float64(limit),
			),
		)
	}

	result.Pagination =
		dto.PARPaginationResponse{
			Page:       page,
			Limit:      limit,
			Total:      int(total),
			TotalPages: totalPages,
		}

	return &result, nil
}
