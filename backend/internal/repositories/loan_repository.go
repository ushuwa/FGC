package repositories

import (
	"errors"
	"math"
	"strings"
	"time"

	"github.com/jj.jobo/FGC/internal/database"
	"github.com/jj.jobo/FGC/internal/dto"
	"gorm.io/gorm"
)

type LoanRepository struct {
	*BaseRepository[struct{}]
}

func NewLoanRepository() *LoanRepository {
	return &LoanRepository{
		BaseRepository: NewBaseRepository[struct{}](
			database.DB,
		),
	}
}

// ============================================================
// FIND ALL LOANS
// ============================================================

func (r *LoanRepository) FindAll(
	search string,
	status string,
) ([]dto.LoanListItem, error) {

	var loans []dto.LoanListItem

	search = strings.TrimSpace(search)
	status = strings.TrimSpace(status)

	query := database.DB.Raw(`
		SELECT
			l.id,
			l.client_id,

			TRIM(
				COALESCE(c.first_name, '') ||
				' ' ||
				COALESCE(c.last_name, '')
			) AS client_name,

			l.pn_number,
			l.loan_type,
			l.principal_amount,
			l.interest_rate,
			l.loan_interest,
			l.pn_value,
			l.loan_term,
			l.frequency,
			l.amortization_amount,

			TO_CHAR(
				l.disbursement_date,
				'YYYY-MM-DD'
			) AS disbursement_date,

			TO_CHAR(
				l.maturity_date,
				'YYYY-MM-DD'
			) AS maturity_date,

			l.status,

			COALESCE(
				(
					SELECT SUM(p.amount_paid)
					FROM payments p
					WHERE p.loan_id = l.id
				),
				0
			) AS total_paid,

			l.pn_value
			-
			COALESCE(
				(
					SELECT SUM(p.amount_paid)
					FROM payments p
					WHERE p.loan_id = l.id
				),
				0
			) AS outstanding_balance

		FROM loans l

		LEFT JOIN clients c
			ON c.id = l.client_id

		WHERE
			(
				? = ''
				OR l.pn_number ILIKE '%' || ? || '%'
				OR l.loan_type ILIKE '%' || ? || '%'
				OR c.first_name ILIKE '%' || ? || '%'
				OR c.last_name ILIKE '%' || ? || '%'
				OR (
					COALESCE(c.first_name, '') ||
					' ' ||
					COALESCE(c.last_name, '')
				) ILIKE '%' || ? || '%'
			)

			AND
			(
				? = ''
				OR l.status = ?
			)

		ORDER BY
			l.id DESC
	`,
		search,
		search,
		search,
		search,
		search,
		search,
		status,
		status,
	)

	err := query.Scan(&loans).Error

	if err != nil {
		return nil, err
	}

	if loans == nil {
		loans = []dto.LoanListItem{}
	}

	return loans, nil
}

// ============================================================
// GET LOAN PROFILE
// ============================================================

func (r *LoanRepository) GetProfile(
	id int,
) (*dto.LoanProfileResponse, error) {

	var loan dto.LoanProfileInfo

	err := database.DB.Raw(`
		SELECT
			id,
			client_id,
			pn_number,
			loan_type,
			principal_amount,
			interest_rate,
			loan_interest,
			pn_value,
			loan_term,
			frequency,
			amortization_amount,

			TO_CHAR(
				disbursement_date,
				'YYYY-MM-DD'
			) AS disbursement_date,

			TO_CHAR(
				maturity_date,
				'YYYY-MM-DD'
			) AS maturity_date,

			status

		FROM loans

		WHERE id = ?
	`, id).Scan(&loan).Error

	if err != nil {
		return nil, err
	}

	if loan.ID == 0 {
		return nil, errors.New(
			"loan not found",
		)
	}

	// ========================================================
	// CLIENT
	// ========================================================

	var client dto.LoanProfileClient

	err = database.DB.Raw(`
		SELECT
			c.id,
			c.first_name,
			c.last_name,
			c.contact_number,
			c.email,
			c.current_address

		FROM clients c

		INNER JOIN loans l
			ON l.client_id = c.id

		WHERE l.id = ?
	`, id).Scan(&client).Error

	if err != nil {
		return nil, err
	}

	if client.ID == 0 {
		return nil, errors.New(
			"client not found",
		)
	}

	// ========================================================
	// SUMMARY
	// ========================================================

	var summary dto.LoanProfileSummary

	err = database.DB.Raw(`
		SELECT
			l.principal_amount,
			l.pn_value,

			COALESCE(
				(
					SELECT SUM(p.amount_paid)
					FROM payments p
					WHERE p.loan_id = l.id
				),
				0
			) AS total_paid,

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
					WHERE a.loan_id = l.id
				),
				0
			) AS outstanding_balance

		FROM loans l

		WHERE l.id = ?
	`, id).Scan(&summary).Error

	if err != nil {
		return nil, err
	}

	// ========================================================
	// AMORTIZATION
	// ========================================================

	var amortizations []dto.LoanAmortization

	err = database.DB.Raw(`
		SELECT
			id,
			loan_id,

			TO_CHAR(
				due_date,
				'YYYY-MM-DD'
			) AS due_date,

			principal_amount,
			interest_amount,
			total_amount,
			paid_principal_amount,
			paid_interest_amount,
			status

		FROM amortization_schedules

		WHERE loan_id = ?

		ORDER BY
			due_date ASC,
			id ASC
	`, id).Scan(&amortizations).Error

	if err != nil {
		return nil, err
	}

	if amortizations == nil {
		amortizations = []dto.LoanAmortization{}
	}

	// ========================================================
	// PAYMENTS
	// ========================================================

	var payments []dto.LoanPayment

	err = database.DB.Raw(`
		SELECT
			id,
			loan_id,

			TO_CHAR(
				payment_date,
				'YYYY-MM-DD'
			) AS payment_date,

			amount_paid,
			payment_channel,
			reference_number,

			COALESCE(
				principal_applied,
				0
			) AS principal_applied,

			COALESCE(
				interest_applied,
				0
			) AS interest_applied,

			COALESCE(
				outstanding_balance,
				0
			) AS outstanding_balance

		FROM payments

		WHERE loan_id = ?

		ORDER BY
			payment_date ASC,
			id ASC
	`, id).Scan(&payments).Error

	if err != nil {
		return nil, err
	}

	if payments == nil {
		payments = []dto.LoanPayment{}
	}

	return &dto.LoanProfileResponse{
		Loan:          loan,
		Client:        client,
		Summary:       summary,
		Amortizations: amortizations,
		Payments:      payments,
	}, nil
}

// ============================================================
// GET PAYMENTS
// ============================================================

func (r *LoanRepository) GetPayments(
	id int,
) ([]dto.LoanPayment, error) {

	var payments []dto.LoanPayment

	err := database.DB.Raw(`
		SELECT
			p.id,
			p.loan_id,

			TO_CHAR(
				p.payment_date,
				'YYYY-MM-DD'
			) AS payment_date,

			p.amount_paid,
			p.payment_channel,
			p.reference_number,

			COALESCE(
				p.principal_applied,
				0
			) AS principal_applied,

			COALESCE(
				p.interest_applied,
				0
			) AS interest_applied,

			COALESCE(
				p.outstanding_balance,
				0
			) AS outstanding_balance

		FROM payments p

		WHERE p.loan_id = ?

		ORDER BY
			p.payment_date DESC,
			p.id DESC
	`,
		id,
	).Scan(&payments).Error

	if err != nil {
		return nil, err
	}

	if payments == nil {
		payments = []dto.LoanPayment{}
	}

	return payments, nil
}

// ============================================================
// CREATE PAYMENT
// ============================================================

func (r *LoanRepository) CreatePayment(
	id int,
	req dto.CreatePaymentRequest,
) (*dto.LoanPayment, error) {

	var result *dto.LoanPayment

	err := database.DB.Transaction(func(tx *gorm.DB) error {

		// ====================================================
		// 1. LOCK LOAN
		// ====================================================

		var loan struct {
			ID               int        `gorm:"column:id"`
			PNValue          float64    `gorm:"column:pn_value"`
			Status           string     `gorm:"column:status"`
			DisbursementDate *time.Time `gorm:"column:disbursement_date"`
		}

		err := tx.Raw(`
			SELECT
				id,
				pn_value,
				status,
				disbursement_date
			FROM loans
			WHERE id = ?
			FOR UPDATE
		`, id).Scan(&loan).Error

		if err != nil {
			return err
		}

		if loan.ID == 0 {
			return errors.New("loan not found")
		}

		if loan.DisbursementDate == nil {
			return errors.New("loan disbursement is required")
		}
		//pang payment date

		paymentDate, err := time.Parse("2006-01-02", req.PaymentDate)
		if err != nil {
			return errors.New("invalid payment date format")
		}

		loanDate := loan.DisbursementDate.Truncate(
			24 * time.Hour,
		)

		paymentDate = paymentDate.Truncate(
			24 * time.Hour,
		)

		if paymentDate.Before(loanDate) {
			return errors.New(
				"payment date cannot be before loan disbursement date",
			)
		}

		// ====================================================
		// 2. PREVENT PAYMENT AFTER FULL PAYMENT
		// ====================================================

		if strings.EqualFold(
			loan.Status,
			"PAID",
		) {
			return errors.New(
				"loan is already fully paid",
			)
		}

		// ====================================================
		// 3. LOCK AMORTIZATION SCHEDULES
		// ====================================================

		type ScheduleRow struct {
			ID              int
			LoanID          int
			DueDate         time.Time
			PrincipalAmount float64
			InterestAmount  float64
			TotalAmount     float64
			PaidPrincipal   float64
			PaidInterest    float64
			Status          string
		}

		var schedules []ScheduleRow

		err = tx.Raw(`
			SELECT
				id,
				loan_id,
				due_date,
				principal_amount,
				interest_amount,
				total_amount,

				COALESCE(
					paid_principal_amount,
					0
				) AS paid_principal,

				COALESCE(
					paid_interest_amount,
					0
				) AS paid_interest,

				status

			FROM amortization_schedules

			WHERE loan_id = ?

			ORDER BY
				due_date ASC,
				id ASC

			FOR UPDATE
		`, id).Scan(&schedules).Error

		if err != nil {
			return err
		}

		if len(schedules) == 0 {
			return errors.New(
				"loan has no amortization schedule",
			)
		}

		// ====================================================
		// 4. VALIDATE PAYMENT
		// ====================================================

		remainingPayment :=
			money(req.AmountPaid)

		if remainingPayment <= 0 {
			return errors.New(
				"payment amount must be greater than zero",
			)
		}

		// ====================================================
		// 5. CALCULATE CURRENT OUTSTANDING
		// ====================================================

		var totalOutstanding float64

		for _, schedule := range schedules {

			remainingInterest := money(
				schedule.InterestAmount -
					schedule.PaidInterest,
			)

			remainingPrincipal := money(
				schedule.PrincipalAmount -
					schedule.PaidPrincipal,
			)

			if remainingInterest < 0 {
				remainingInterest = 0
			}

			if remainingPrincipal < 0 {
				remainingPrincipal = 0
			}

			totalOutstanding = money(
				totalOutstanding +
					remainingInterest +
					remainingPrincipal,
			)
		}

		totalOutstanding =
			money(totalOutstanding)

		if totalOutstanding <= 0 {
			return errors.New(
				"loan is already fully paid",
			)
		}

		// ====================================================
		// 6. PREVENT OVERPAYMENT
		// ====================================================

		if remainingPayment >
			totalOutstanding {

			return errors.New(
				"payment exceeds loan outstanding balance",
			)
		}

		// ====================================================
		// 7. APPLY PAYMENT
		// ====================================================

		var totalInterestApplied float64
		var totalPrincipalApplied float64

		for i := range schedules {

			if remainingPayment <= 0 {
				break
			}

			schedule :=
				&schedules[i]

			// ------------------------------------------------
			// REMAINING INTEREST
			// ------------------------------------------------

			remainingInterest := money(
				schedule.InterestAmount -
					schedule.PaidInterest,
			)

			if remainingInterest < 0 {
				remainingInterest = 0
			}

			// ------------------------------------------------
			// INTEREST FIRST
			// ------------------------------------------------

			interestApplied := money(
				math.Min(
					remainingPayment,
					remainingInterest,
				),
			)

			if interestApplied > 0 {

				schedule.PaidInterest =
					money(
						schedule.PaidInterest +
							interestApplied,
					)

				remainingPayment =
					money(
						remainingPayment -
							interestApplied,
					)

				totalInterestApplied =
					money(
						totalInterestApplied +
							interestApplied,
					)
			}

			// ------------------------------------------------
			// REMAINING PRINCIPAL
			// ------------------------------------------------

			remainingPrincipal := money(
				schedule.PrincipalAmount -
					schedule.PaidPrincipal,
			)

			if remainingPrincipal < 0 {
				remainingPrincipal = 0
			}

			// ------------------------------------------------
			// PRINCIPAL SECOND
			// ------------------------------------------------

			principalApplied := money(
				math.Min(
					remainingPayment,
					remainingPrincipal,
				),
			)

			if principalApplied > 0 {

				schedule.PaidPrincipal =
					money(
						schedule.PaidPrincipal +
							principalApplied,
					)

				remainingPayment =
					money(
						remainingPayment -
							principalApplied,
					)

				totalPrincipalApplied =
					money(
						totalPrincipalApplied +
							principalApplied,
					)
			}

			// ------------------------------------------------
			// UPDATE SCHEDULE STATUS
			// ------------------------------------------------

			remainingInterest = money(
				schedule.InterestAmount -
					schedule.PaidInterest,
			)

			remainingPrincipal = money(
				schedule.PrincipalAmount -
					schedule.PaidPrincipal,
			)

			var status string

			switch {

			case remainingInterest <= 0 &&
				remainingPrincipal <= 0:

				status = "PAID"

			case schedule.PaidInterest > 0 ||
				schedule.PaidPrincipal > 0:

				status = "PARTIAL"

			default:

				status = "PENDING"
			}

			schedule.Status = status

			// ------------------------------------------------
			// SAVE SCHEDULE
			// ------------------------------------------------

			err = tx.Exec(`
				UPDATE amortization_schedules
				SET
					paid_principal_amount = ?,
					paid_interest_amount = ?,
					status = ?
				WHERE id = ?
			`,
				schedule.PaidPrincipal,
				schedule.PaidInterest,
				schedule.Status,
				schedule.ID,
			).Error

			if err != nil {
				return err
			}
		}

		// ====================================================
		// 8. SAFETY CHECK
		// ====================================================

		remainingPayment =
			money(remainingPayment)

		if remainingPayment > 0 {
			return errors.New(
				"payment could not be fully allocated",
			)
		}

		// ====================================================
		// 9. NEW OUTSTANDING
		// ====================================================

		newOutstanding := money(
			totalOutstanding -
				req.AmountPaid,
		)

		if newOutstanding < 0 {
			newOutstanding = 0
		}

		// ====================================================
		// 10. DETERMINE LOAN STATUS
		// ====================================================

		newStatus := "ACTIVE"

		if newOutstanding <= 0.01 {
			newOutstanding = 0
			newStatus = "PAID"
		}

		// ====================================================
		// 11. INSERT PAYMENT
		// ====================================================

		var payment dto.LoanPayment

		err = tx.Raw(`
			INSERT INTO payments (
				loan_id,
				payment_date,
				amount_paid,
				payment_channel,
				reference_number,
				principal_applied,
				interest_applied,
				outstanding_balance
			)
			VALUES (
				?,
				?::date,
				?,
				?,
				?,
				?,
				?,
				?
			)
			RETURNING
				id,
				loan_id,
				TO_CHAR(
					payment_date,
					'YYYY-MM-DD'
				) AS payment_date,
				amount_paid,
				payment_channel,
				reference_number,
				principal_applied,
				interest_applied,
				outstanding_balance
		`,
			id,
			req.PaymentDate,
			money(req.AmountPaid),
			req.PaymentChannel,
			req.ReferenceNumber,
			money(totalPrincipalApplied),
			money(totalInterestApplied),
			newOutstanding,
		).Scan(&payment).Error

		if err != nil {
			return err
		}

		// ====================================================
		// 12. UPDATE LOAN STATUS
		// ====================================================

		err = tx.Exec(`
			UPDATE loans
			SET status = ?
			WHERE id = ?
		`,
			newStatus,
			id,
		).Error

		if err != nil {
			return err
		}

		result = &payment

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *LoanRepository) DeletePayment(
	id int,
) error {

	return database.DB.Transaction(func(tx *gorm.DB) error {

		// ========================================
		// 1. FIND PAYMENT + LOCK LOAN
		// ========================================

		var payment struct {
			ID     int `gorm:"column:id"`
			LoanID int `gorm:"column:loan_id"`
		}

		err := tx.Raw(`
			SELECT
				p.id,
				p.loan_id
			FROM payments p
			WHERE p.id = ?
			FOR UPDATE
		`, id).Scan(&payment).Error

		if err != nil {
			return err
		}

		if payment.ID == 0 {
			return errors.New("payment not found")
		}

		if payment.LoanID <= 0 {
			return errors.New("payment has no loan")
		}

		// ========================================
		// 2. LOCK LOAN
		// ========================================

		var loanID int

		err = tx.Raw(`
			SELECT id
			FROM loans
			WHERE id = ?
			FOR UPDATE
		`,
			payment.LoanID,
		).Scan(&loanID).Error

		if err != nil {
			return err
		}

		if loanID == 0 {
			return errors.New("loan not found")
		}

		// ========================================
		// 3. DELETE PAYMENT
		// ========================================

		result := tx.Exec(`
			DELETE FROM payments
			WHERE id = ?
		`, id)

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return errors.New("payment not found")
		}

		// ========================================
		// 4. REBUILD AMORTIZATION
		// ========================================

		return rebuildAmortizationTx(
			tx,
			payment.LoanID,
		)
	})
}

// ============================================================
// REBUILD AMORTIZATION
// ============================================================

func (r *LoanRepository) RebuildAmortization(
	id int,
) error {

	return database.DB.Transaction(
		func(tx *gorm.DB) error {
			return rebuildAmortizationTx(
				tx,
				id,
			)
		},
	)
}

// ============================================================
// REBUILD AMORTIZATION TRANSACTION
// ============================================================

func rebuildAmortizationTx(
	tx *gorm.DB,
	id int,
) error {

	// ========================================================
	// 1. LOCK LOAN
	// ========================================================

	var loan struct {
		ID               int        `gorm:"column:id"`
		PrincipalAmount  float64    `gorm:"column:principal_amount"`
		LoanInterest     float64    `gorm:"column:loan_interest"`
		PNValue          float64    `gorm:"column:pn_value"`
		LoanTerm         int        `gorm:"column:loan_term"`
		Frequency        string     `gorm:"column:frequency"`
		DisbursementDate *time.Time `gorm:"column:disbursement_date"`
	}

	err := tx.Raw(`
		SELECT
			id,
			principal_amount,
			loan_interest,
			pn_value,
			loan_term,
			frequency,
			disbursement_date
		FROM loans
		WHERE id = ?
		FOR UPDATE
	`, id).Scan(&loan).Error

	if err != nil {
		return err
	}

	if loan.ID == 0 {
		return errors.New(
			"loan not found",
		)
	}

	if loan.LoanTerm <= 0 {
		return errors.New(
			"loan term must be greater than zero",
		)
	}

	if loan.DisbursementDate == nil {
		return errors.New(
			"loan disbursement date is required",
		)
	}

	frequency := strings.ToUpper(
		strings.TrimSpace(
			loan.Frequency,
		),
	)

	if frequency != "MONTHLY" &&
		frequency != "SEMI-MONTHLY" {

		return errors.New(
			"frequency must be MONTHLY or SEMI-MONTHLY",
		)
	}

	// ========================================================
	// 2. VALIDATE PN VALUE
	// ========================================================

	expectedPNValue := money(
		loan.PrincipalAmount +
			loan.LoanInterest,
	)

	if money(loan.PNValue) != expectedPNValue {
		return errors.New(
			"PN value must equal principal amount plus loan interest",
		)
	}

	// ========================================================
	// 3. READ EXISTING PAYMENTS
	// ========================================================

	type ExistingPayment struct {
		ID          int       `gorm:"column:id"`
		PaymentDate time.Time `gorm:"column:payment_date"`
		AmountPaid  float64   `gorm:"column:amount_paid"`
	}

	var payments []ExistingPayment

	err = tx.Raw(`
		SELECT
			id,
			payment_date,
			amount_paid
		FROM payments
		WHERE loan_id = ?
		ORDER BY
			payment_date ASC,
			id ASC
		FOR UPDATE
	`, id).Scan(&payments).Error

	if err != nil {
		return err
	}

	// ========================================================
	// 4. DELETE OLD AMORTIZATION
	// ========================================================

	err = tx.Exec(`
		DELETE FROM amortization_schedules
		WHERE loan_id = ?
	`, id).Error

	if err != nil {
		return err
	}

	// ========================================================
	// 5. DETERMINE NUMBER OF PERIODS
	// ========================================================

	periods := loan.LoanTerm

	if frequency == "SEMI-MONTHLY" {
		periods = loan.LoanTerm * 2
	}

	if periods <= 0 {
		return errors.New(
			"invalid number of amortization periods",
		)
	}

	// ========================================================
	// 6. PREPARE SCHEDULE
	// ========================================================

	type Schedule struct {
		DueDate             time.Time
		PrincipalAmount     float64
		InterestAmount      float64
		TotalAmount         float64
		PaidPrincipalAmount float64
		PaidInterestAmount  float64
		Status              string
	}

	schedules := make(
		[]Schedule,
		0,
		periods,
	)

	principalPerPeriod := money(
		loan.PrincipalAmount /
			float64(periods),
	)

	interestPerPeriod := money(
		loan.LoanInterest /
			float64(periods),
	)

	var principalAllocated float64
	var interestAllocated float64

	// ========================================================
	// 7. GENERATE DUE DATES
	// ========================================================

	currentDate :=
		*loan.DisbursementDate

	for i := 0; i < periods; i++ {

		principal := principalPerPeriod
		interest := interestPerPeriod

		// ----------------------------------------------------
		// LAST PERIOD GETS ROUNDING DIFFERENCE
		// ----------------------------------------------------

		if i == periods-1 {

			principal = money(
				loan.PrincipalAmount -
					principalAllocated,
			)

			interest = money(
				loan.LoanInterest -
					interestAllocated,
			)
		}

		principal = money(principal)
		interest = money(interest)

		total := money(
			principal + interest,
		)

		principalAllocated = money(
			principalAllocated + principal,
		)

		interestAllocated = money(
			interestAllocated + interest,
		)

		// ----------------------------------------------------
		// CALCULATE NEXT DUE DATE
		// ----------------------------------------------------

		dueDate := nextDueDate(
			currentDate,
			frequency,
		)

		currentDate = dueDate

		schedules = append(
			schedules,
			Schedule{
				DueDate:             dueDate,
				PrincipalAmount:     principal,
				InterestAmount:      interest,
				TotalAmount:         total,
				PaidPrincipalAmount: 0,
				PaidInterestAmount:  0,
				Status:              "PENDING",
			},
		)
	}

	// ========================================================
	// 8. REPLAY EXISTING PAYMENTS
	// ========================================================

	type PaymentAllocation struct {
		ID                 int
		PrincipalApplied   float64
		InterestApplied    float64
		OutstandingBalance float64
	}

	allocations := make(
		[]PaymentAllocation,
		0,
		len(payments),
	)

	totalLoanAmount := money(
		loan.PrincipalAmount +
			loan.LoanInterest,
	)

	var totalPaid float64

	for _, payment := range payments {

		remaining := money(
			payment.AmountPaid,
		)

		if remaining <= 0 {
			return errors.New(
				"existing payment amount must be greater than zero",
			)
		}

		var paymentInterest float64
		var paymentPrincipal float64

		// ====================================================
		// APPLY TO EARLIEST SCHEDULE
		// INTEREST FIRST
		// THEN PRINCIPAL
		// ====================================================

		for i := range schedules {

			if remaining <= 0 {
				break
			}

			schedule :=
				&schedules[i]

			// ------------------------------------------------
			// REMAINING INTEREST
			// ------------------------------------------------

			remainingInterest := money(
				schedule.InterestAmount -
					schedule.PaidInterestAmount,
			)

			if remainingInterest < 0 {
				remainingInterest = 0
			}

			interestApplied := money(
				math.Min(
					remaining,
					remainingInterest,
				),
			)

			if interestApplied > 0 {

				schedule.PaidInterestAmount =
					money(
						schedule.PaidInterestAmount +
							interestApplied,
					)

				remaining = money(
					remaining -
						interestApplied,
				)

				paymentInterest = money(
					paymentInterest +
						interestApplied,
				)
			}

			// ------------------------------------------------
			// REMAINING PRINCIPAL
			// ------------------------------------------------

			remainingPrincipal := money(
				schedule.PrincipalAmount -
					schedule.PaidPrincipalAmount,
			)

			if remainingPrincipal < 0 {
				remainingPrincipal = 0
			}

			principalApplied := money(
				math.Min(
					remaining,
					remainingPrincipal,
				),
			)

			if principalApplied > 0 {

				schedule.PaidPrincipalAmount =
					money(
						schedule.PaidPrincipalAmount +
							principalApplied,
					)

				remaining = money(
					remaining -
						principalApplied,
				)

				paymentPrincipal = money(
					paymentPrincipal +
						principalApplied,
				)
			}

			// ------------------------------------------------
			// UPDATE SCHEDULE STATUS
			// ------------------------------------------------

			remainingInterest = money(
				schedule.InterestAmount -
					schedule.PaidInterestAmount,
			)

			remainingPrincipal = money(
				schedule.PrincipalAmount -
					schedule.PaidPrincipalAmount,
			)

			switch {

			case remainingInterest <= 0 &&
				remainingPrincipal <= 0:

				schedule.Status = "PAID"

			case schedule.PaidInterestAmount > 0 ||
				schedule.PaidPrincipalAmount > 0:

				schedule.Status = "PARTIAL"

			default:

				schedule.Status = "PENDING"
			}
		}

		// ====================================================
		// PAYMENT MUST BE FULLY ALLOCATED
		// ====================================================

		remaining = money(remaining)

		if remaining > 0 {
			return errors.New(
				"existing payments exceed the generated amortization schedule",
			)
		}

		totalPaid = money(
			totalPaid +
				payment.AmountPaid,
		)

		outstanding := money(
			totalLoanAmount -
				totalPaid,
		)

		if outstanding < 0 {
			outstanding = 0
		}

		allocations = append(
			allocations,
			PaymentAllocation{
				ID: payment.ID,

				PrincipalApplied: money(
					paymentPrincipal,
				),

				InterestApplied: money(
					paymentInterest,
				),

				OutstandingBalance: outstanding,
			},
		)
	}

	// ========================================================
	// 9. INSERT NEW SCHEDULES
	// ========================================================

	for _, schedule := range schedules {

		err = tx.Exec(`
			INSERT INTO amortization_schedules (
				loan_id,
				due_date,
				principal_amount,
				interest_amount,
				total_amount,
				paid_principal_amount,
				paid_interest_amount,
				status
			)
			VALUES (
				?,
				?,
				?,
				?,
				?,
				?,
				?,
				?
			)
		`,
			id,
			schedule.DueDate,
			schedule.PrincipalAmount,
			schedule.InterestAmount,
			schedule.TotalAmount,
			schedule.PaidPrincipalAmount,
			schedule.PaidInterestAmount,
			schedule.Status,
		).Error

		if err != nil {
			return err
		}
	}

	// ========================================================
	// 10. UPDATE PAYMENT HISTORY
	// ========================================================

	for _, allocation := range allocations {

		err = tx.Exec(`
			UPDATE payments
			SET
				principal_applied = ?,
				interest_applied = ?,
				outstanding_balance = ?
			WHERE id = ?
		`,
			allocation.PrincipalApplied,
			allocation.InterestApplied,
			allocation.OutstandingBalance,
			allocation.ID,
		).Error

		if err != nil {
			return err
		}
	}

	// ========================================================
	// 11. CALCULATE FINAL OUTSTANDING
	// ========================================================

	var finalOutstanding float64

	for _, schedule := range schedules {

		remainingPrincipal := money(
			schedule.PrincipalAmount -
				schedule.PaidPrincipalAmount,
		)

		remainingInterest := money(
			schedule.InterestAmount -
				schedule.PaidInterestAmount,
		)

		if remainingPrincipal < 0 {
			remainingPrincipal = 0
		}

		if remainingInterest < 0 {
			remainingInterest = 0
		}

		finalOutstanding = money(
			finalOutstanding +
				remainingPrincipal +
				remainingInterest,
		)
	}

	// ========================================================
	// 12. UPDATE LOAN STATUS
	// ========================================================

	loanStatus := "ACTIVE"

	if finalOutstanding <= 0.01 {
		loanStatus = "PAID"
		finalOutstanding = 0
	}

	err = tx.Exec(`
		UPDATE loans
		SET status = ?
		WHERE id = ?
	`,
		loanStatus,
		id,
	).Error

	if err != nil {
		return err
	}

	return nil
}

// ============================================================
// NEXT DUE DATE
// ============================================================

func nextDueDate(
	date time.Time,
	frequency string,
) time.Time {

	switch strings.ToUpper(
		strings.TrimSpace(frequency),
	) {

	case "MONTHLY":

		return date.AddDate(
			0,
			1,
			0,
		)

	case "SEMI-MONTHLY":

		return nextSemiMonthlyDate(
			date,
		)

	default:

		return date.AddDate(
			0,
			1,
			0,
		)
	}
}

// ============================================================
// NEXT SEMI-MONTHLY DATE
// ============================================================

func nextSemiMonthlyDate(
	current time.Time,
) time.Time {

	year := current.Year()
	month := current.Month()
	location := current.Location()

	// ========================================================
	// BEFORE 15TH
	// ========================================================

	if current.Day() < 15 {

		return time.Date(
			year,
			month,
			15,
			0,
			0,
			0,
			0,
			location,
		)
	}

	// ========================================================
	// 15TH -> LAST DAY OF CURRENT MONTH
	// ========================================================

	if current.Day() == 15 {

		return time.Date(
			year,
			month+1,
			0,
			0,
			0,
			0,
			0,
			location,
		)
	}

	// ========================================================
	// LAST DAY -> 15TH OF NEXT MONTH
	// ========================================================

	return time.Date(
		year,
		month+1,
		15,
		0,
		0,
		0,
		0,
		location,
	)
}

// ============================================================
// CREATE LOAN
// ============================================================

func (r *LoanRepository) CreateLoan(
	req dto.CreateLoanRequest,
) error {

	return database.DB.Transaction(func(tx *gorm.DB) error {

		// ========================================
		// 1. INSERT LOAN
		// ========================================

		var loan struct {
			ID int `gorm:"column:id"`
		}

		err := tx.Raw(`
			INSERT INTO loans (
				client_id,
				pn_number,
				loan_type,
				principal_amount,
				interest_rate,
				loan_interest,
				pn_value,
				loan_term,
				amortization_amount,
				disbursement_date,
				maturity_date,
				frequency,
				status
			)
			VALUES (
				?,
				?,
				?,
				?,
				?,
				?,
				?,
				?,
				?,
				?::date,
				?::date,
				?,
				?
			)
			RETURNING id
		`,
			req.ClientID,
			req.PNNumber,
			req.LoanType,
			req.PrincipalAmount,
			req.InterestRate,
			req.LoanInterest,
			req.PNValue,
			req.LoanTerm,
			req.AmortizationAmount,
			req.DisbursementDate,
			req.MaturityDate,
			req.Frequency,
			req.Status,
		).Scan(&loan).Error

		if err != nil {
			return err
		}

		if loan.ID == 0 {
			return errors.New(
				"failed to create loan",
			)
		}

		// ========================================
		// 2. GENERATE INITIAL AMORTIZATION
		// ========================================

		err = rebuildLoanScheduleTx(
			tx,
			loan.ID,
		)

		if err != nil {
			return err
		}

		return nil
	})
}

func rebuildLoanScheduleTx(
	tx *gorm.DB,
	loanID int,
) error {

	// ========================================
	// 1. GET LOAN
	// ========================================

	var loan struct {
		ID               int        `gorm:"column:id"`
		PrincipalAmount  float64    `gorm:"column:principal_amount"`
		LoanInterest     float64    `gorm:"column:loan_interest"`
		PNValue          float64    `gorm:"column:pn_value"`
		LoanTerm         int        `gorm:"column:loan_term"`
		DisbursementDate *time.Time `gorm:"column:disbursement_date"`
		Frequency        string     `gorm:"column:frequency"`
	}

	err := tx.Raw(`
		SELECT
			id,
			principal_amount,
			loan_interest,
			pn_value,
			loan_term,
			disbursement_date,
			frequency
		FROM loans
		WHERE id = ?
		FOR UPDATE
	`, loanID).Scan(&loan).Error

	if err != nil {
		return err
	}

	if loan.ID == 0 {
		return errors.New(
			"loan not found",
		)
	}

	// ========================================
	// 2. VALIDATE LOAN
	// ========================================

	if loan.LoanTerm <= 0 {
		return errors.New(
			"loan term must be greater than zero",
		)
	}

	if loan.DisbursementDate == nil {
		return errors.New(
			"loan disbursement date is required",
		)
	}

	frequency := strings.ToUpper(
		strings.TrimSpace(
			loan.Frequency,
		),
	)

	if frequency != "MONTHLY" &&
		frequency != "SEMI-MONTHLY" {

		return errors.New(
			"frequency must be MONTHLY or SEMI-MONTHLY",
		)
	}

	// ========================================
	// 3. VALIDATE PN VALUE
	// ========================================

	expectedPNValue := money(
		loan.PrincipalAmount +
			loan.LoanInterest,
	)

	if money(loan.PNValue) != expectedPNValue {
		return errors.New(
			"PN value must equal principal amount plus loan interest",
		)
	}

	// ========================================
	// 4. DELETE EXISTING SCHEDULE
	// ========================================

	err = tx.Exec(`
		DELETE FROM amortization_schedules
		WHERE loan_id = ?
	`, loanID).Error

	if err != nil {
		return err
	}

	// ========================================
	// 5. DETERMINE PERIODS
	// ========================================

	periods := loan.LoanTerm

	if frequency == "SEMI-MONTHLY" {
		periods = loan.LoanTerm * 2
	}

	if periods <= 0 {
		return errors.New(
			"invalid number of amortization periods",
		)
	}

	// ========================================
	// 6. CALCULATE AMOUNTS
	// ========================================

	principalPerPeriod := money(
		loan.PrincipalAmount /
			float64(periods),
	)

	interestPerPeriod := money(
		loan.LoanInterest /
			float64(periods),
	)

	var principalAllocated float64
	var interestAllocated float64

	// ========================================
	// 7. GENERATE SCHEDULE
	// ========================================

	currentDate :=
		*loan.DisbursementDate

	for period := 1; period <= periods; period++ {

		principal :=
			principalPerPeriod

		interest :=
			interestPerPeriod

		// ----------------------------------------
		// LAST PERIOD GETS ROUNDING DIFFERENCE
		// ----------------------------------------

		if period == periods {

			principal = money(
				loan.PrincipalAmount -
					principalAllocated,
			)

			interest = money(
				loan.LoanInterest -
					interestAllocated,
			)
		}

		principal = money(principal)
		interest = money(interest)

		total := money(
			principal + interest,
		)

		principalAllocated =
			money(
				principalAllocated +
					principal,
			)

		interestAllocated =
			money(
				interestAllocated +
					interest,
			)

		// ----------------------------------------
		// CALCULATE DUE DATE
		// ----------------------------------------

		var dueDate time.Time

		if frequency == "MONTHLY" {

			dueDate = currentDate.AddDate(
				0,
				1,
				0,
			)

		} else {

			dueDate = nextSemiMonthlyDate(
				currentDate,
			)
		}

		currentDate = dueDate

		// ----------------------------------------
		// INSERT SCHEDULE
		// ----------------------------------------

		err = tx.Exec(`
			INSERT INTO amortization_schedules (
				loan_id,
				due_date,
				principal_amount,
				interest_amount,
				total_amount,
				paid_principal_amount,
				paid_interest_amount,
				status
			)
			VALUES (
				?,
				?,
				?,
				?,
				?,
				0,
				0,
				'PENDING'
			)
		`,
			loanID,
			dueDate,
			principal,
			interest,
			total,
		).Error

		if err != nil {
			return err
		}
	}

	return nil
}

// ============================================================
// UPDATE LOAN
// ============================================================

func (r *LoanRepository) UpdateLoan(
	id int,
	req dto.UpdateLoanRequest,
) error {

	return database.DB.Transaction(
		func(tx *gorm.DB) error {

			// =================================================
			// 1. GET CURRENT LOAN
			// =================================================

			var oldLoan struct {
				PrincipalAmount  float64   `gorm:"column:principal_amount"`
				LoanInterest     float64   `gorm:"column:loan_interest"`
				LoanTerm         int       `gorm:"column:loan_term"`
				Frequency        string    `gorm:"column:frequency"`
				DisbursementDate time.Time `gorm:"column:disbursement_date"`
				Status           string    `gorm:"column:status"`
			}

			err := tx.Raw(`
				SELECT
					principal_amount,
					loan_interest,
					loan_term,
					frequency,
					disbursement_date,
					status
				FROM loans
				WHERE id = ?
				FOR UPDATE
			`, id).Scan(&oldLoan).Error

			if err != nil {
				return err
			}

			// =================================================
			// 2. CHECK LOAN EXISTS
			// =================================================

			var exists int

			err = tx.Raw(`
				SELECT id
				FROM loans
				WHERE id = ?
			`, id).Scan(&exists).Error

			if err != nil {
				return err
			}

			if exists == 0 {
				return errors.New(
					"loan not found",
				)
			}

			// =================================================
			// 3. NORMALIZE FREQUENCY
			// =================================================

			newFrequency := strings.ToUpper(
				strings.TrimSpace(
					req.Frequency,
				),
			)

			if newFrequency != "MONTHLY" &&
				newFrequency != "SEMI-MONTHLY" {

				return errors.New(
					"frequency must be MONTHLY or SEMI-MONTHLY",
				)
			}

			// =================================================
			// 4. CHECK AMORTIZATION CHANGE
			// =================================================

			oldFrequency := strings.ToUpper(
				strings.TrimSpace(
					oldLoan.Frequency,
				),
			)

			oldDate := ""

			if !oldLoan.DisbursementDate.IsZero() {
				oldDate =
					oldLoan.DisbursementDate.Format(
						"2006-01-02",
					)
			}

			newDate := strings.TrimSpace(
				req.DisbursementDate,
			)

			amortizationChanged :=
				money(
					oldLoan.PrincipalAmount,
				) != money(
					req.PrincipalAmount,
				) ||

					money(
						oldLoan.LoanInterest,
					) != money(
						req.LoanInterest,
					) ||

					oldLoan.LoanTerm !=
						req.LoanTerm ||

					oldFrequency !=
						newFrequency ||

					oldDate !=
						newDate

			// =================================================
			// 5. UPDATE LOAN
			// =================================================

			result := tx.Exec(`
				UPDATE loans
				SET
					client_id = ?,
					pn_number = ?,
					loan_type = ?,
					principal_amount = ?,
					interest_rate = ?,
					loan_interest = ?,
					pn_value = ?,
					loan_term = ?,
					amortization_amount = ?,
					disbursement_date = ?::date,
					maturity_date = ?::date,
					frequency = ?
				WHERE id = ?
			`,
				req.ClientID,
				req.PNNumber,
				req.LoanType,
				req.PrincipalAmount,
				req.InterestRate,
				req.LoanInterest,
				req.PNValue,
				req.LoanTerm,
				req.AmortizationAmount,
				req.DisbursementDate,
				req.MaturityDate,
				newFrequency,
				id,
			)

			if result.Error != nil {
				return result.Error
			}

			if result.RowsAffected == 0 {
				return errors.New(
					"loan not found",
				)
			}

			// =================================================
			// 6. REBUILD IF AMORTIZATION CHANGED
			// =================================================

			if amortizationChanged {

				return rebuildAmortizationTx(
					tx,
					id,
				)
			}

			// =================================================
			// 7. HANDLE MANUAL STATUS
			// =================================================
			//
			// PAID should not be manually forced.
			// ACTIVE can be restored manually.
			// CLOSED / DEFAULTED remain manual states.
			//
			// If current status is PAID and no
			// amortization change happened, don't
			// blindly overwrite it.
			// =================================================

			requestedStatus :=
				strings.ToUpper(
					strings.TrimSpace(
						req.Status,
					),
				)

			switch requestedStatus {

			case "ACTIVE",
				"CLOSED",
				"DEFAULTED":

				err = tx.Exec(`
					UPDATE loans
					SET status = ?
					WHERE id = ?
				`,
					requestedStatus,
					id,
				).Error

				if err != nil {
					return err
				}

			case "PAID":

				// ---------------------------------------------
				// Don't allow the UI to manually mark
				// an unpaid loan as PAID.
				// ---------------------------------------------

				var outstanding float64

				err = tx.Raw(`
					SELECT COALESCE(
						SUM(
							GREATEST(
								principal_amount -
								COALESCE(
									paid_principal_amount,
									0
								),
								0
							)
							+
							GREATEST(
								interest_amount -
								COALESCE(
									paid_interest_amount,
									0
								),
								0
							)
						),
						0
					)
					FROM amortization_schedules
					WHERE loan_id = ?
				`,
					id,
				).Scan(&outstanding).Error

				if err != nil {
					return err
				}

				if money(outstanding) > 0.01 {
					return errors.New(
						"loan cannot be marked PAID while it has outstanding balance",
					)
				}

				err = tx.Exec(`
					UPDATE loans
					SET status = 'PAID'
					WHERE id = ?
				`,
					id,
				).Error

				if err != nil {
					return err
				}

			default:

				return errors.New(
					"invalid loan status",
				)
			}

			return nil
		},
	)
}

// ============================================================
// MONEY
// ============================================================

func money(
	value float64,
) float64 {

	return math.Round(
		value*100,
	) / 100
}
