package services

import (
	"errors"
	"strings"

	"github.com/jj.jobo/FGC/internal/dto"
	"github.com/jj.jobo/FGC/internal/repositories"
)

type LoanService struct {
	loanRepository *repositories.LoanRepository
}

func NewLoanService(
	loanRepository *repositories.LoanRepository,
) *LoanService {
	return &LoanService{
		loanRepository: loanRepository,
	}
}

// ========================================
// LOANS
// ========================================

func (s *LoanService) GetLoans(
	search string,
	status string,
) ([]dto.LoanListItem, error) {

	search = strings.TrimSpace(search)

	status = strings.ToUpper(
		strings.TrimSpace(status),
	)

	return s.loanRepository.FindAll(
		search,
		status,
	)
}

func (s *LoanService) GetLoan(
	id int,
) (*dto.LoanProfileResponse, error) {

	if id <= 0 {
		return nil, errors.New(
			"invalid loan id",
		)
	}

	return s.loanRepository.GetProfile(id)
}

// ========================================
// PAYMENTS
// ========================================

func (s *LoanService) GetPayments(
	id int,
) ([]dto.LoanPayment, error) {

	return s.loanRepository.GetPayments(id)
}

func (s *LoanService) CreatePayment(
	id int,
	req dto.CreatePaymentRequest,
) (*dto.LoanPayment, error) {

	if id <= 0 {
		return nil, errors.New(
			"invalid loan ID",
		)
	}

	req.PaymentDate =
		strings.TrimSpace(
			req.PaymentDate,
		)

	if req.PaymentDate == "" {
		return nil, errors.New(
			"payment date is required",
		)
	}

	if req.AmountPaid <= 0 {
		return nil, errors.New(
			"payment amount must be greater than zero",
		)
	}

	if req.PaymentChannel != nil {

		value :=
			strings.TrimSpace(
				*req.PaymentChannel,
			)

		if value == "" {
			req.PaymentChannel = nil
		} else {
			req.PaymentChannel = &value
		}
	}

	if req.ReferenceNumber != nil {

		value :=
			strings.TrimSpace(
				*req.ReferenceNumber,
			)

		if value == "" {
			req.ReferenceNumber = nil
		} else {
			req.ReferenceNumber = &value
		}
	}

	return s.loanRepository.CreatePayment(
		id,
		req,
	)
}

func (s *LoanService) DeletePayment(
	id int,
) error {

	if id <= 0 {
		return errors.New(
			"invalid payment ID",
		)
	}

	return s.loanRepository.DeletePayment(id)
}

// ========================================
// REBUILD
// ========================================

func (s *LoanService) RebuildAmortization(
	id int,
) error {

	if id <= 0 {
		return errors.New(
			"invalid loan ID",
		)
	}

	return s.loanRepository.RebuildAmortization(id)
}

// ========================================
// CREATE LOAN
// ========================================

func (s *LoanService) CreateLoan(
	req dto.CreateLoanRequest,
) error {

	// ----------------------------------------
	// NORMALIZE
	// ----------------------------------------

	req.PNNumber =
		strings.TrimSpace(
			req.PNNumber,
		)

	req.Frequency =
		strings.ToUpper(
			strings.TrimSpace(
				req.Frequency,
			),
		)

	// ----------------------------------------
	// BASIC VALIDATION
	// ----------------------------------------

	if req.ClientID <= 0 {
		return errors.New(
			"invalid client ID",
		)
	}

	if req.PNNumber == "" {
		return errors.New(
			"PN number is required",
		)
	}

	if req.PrincipalAmount <= 0 {
		return errors.New(
			"principal amount must be greater than zero",
		)
	}

	if req.InterestRate < 0 {
		return errors.New(
			"interest rate cannot be negative",
		)
	}

	if req.LoanInterest < 0 {
		return errors.New(
			"loan interest cannot be negative",
		)
	}

	if req.LoanTerm <= 0 {
		return errors.New(
			"loan term must be greater than zero",
		)
	}

	if req.AmortizationAmount <= 0 {
		return errors.New(
			"amortization amount must be greater than zero",
		)
	}

	// ----------------------------------------
	// FREQUENCY
	// ----------------------------------------

	if req.Frequency != "MONTHLY" &&
		req.Frequency != "SEMI-MONTHLY" {

		return errors.New(
			"invalid loan frequency",
		)
	}

	// ----------------------------------------
	// PN VALUE
	// ----------------------------------------

	expectedPNValue :=
		req.PrincipalAmount +
			req.LoanInterest

	if !moneyEqual(
		req.PNValue,
		expectedPNValue,
	) {
		return errors.New(
			"PN value must equal principal amount plus loan interest",
		)
	}

	// ----------------------------------------
	// SAVE
	// ----------------------------------------

	return s.loanRepository.CreateLoan(req)
}

// ========================================
// UPDATE LOAN
// ========================================

func (s *LoanService) UpdateLoan(
	id int,
	req dto.UpdateLoanRequest,
) error {

	if id <= 0 {
		return errors.New(
			"invalid loan ID",
		)
	}

	// ----------------------------------------
	// NORMALIZE
	// ----------------------------------------

	req.PNNumber =
		strings.TrimSpace(
			req.PNNumber,
		)

	req.Frequency =
		strings.ToUpper(
			strings.TrimSpace(
				req.Frequency,
			),
		)

	// ----------------------------------------
	// BASIC VALIDATION
	// ----------------------------------------

	if req.ClientID <= 0 {
		return errors.New(
			"invalid client ID",
		)
	}

	if req.PNNumber == "" {
		return errors.New(
			"PN number is required",
		)
	}

	if req.PrincipalAmount <= 0 {
		return errors.New(
			"principal amount must be greater than zero",
		)
	}

	if req.InterestRate < 0 {
		return errors.New(
			"interest rate cannot be negative",
		)
	}

	if req.LoanInterest < 0 {
		return errors.New(
			"loan interest cannot be negative",
		)
	}

	if req.LoanTerm <= 0 {
		return errors.New(
			"loan term must be greater than zero",
		)
	}

	if req.AmortizationAmount <= 0 {
		return errors.New(
			"amortization amount must be greater than zero",
		)
	}

	// ----------------------------------------
	// FREQUENCY
	// ----------------------------------------

	if req.Frequency != "MONTHLY" &&
		req.Frequency != "SEMI-MONTHLY" {

		return errors.New(
			"invalid loan frequency",
		)
	}

	// ----------------------------------------
	// PN VALUE
	// ----------------------------------------

	expectedPNValue :=
		req.PrincipalAmount +
			req.LoanInterest

	if !moneyEqual(
		req.PNValue,
		expectedPNValue,
	) {
		return errors.New(
			"PN value must equal principal amount plus loan interest",
		)
	}

	// ----------------------------------------
	// SAVE
	// ----------------------------------------

	return s.loanRepository.UpdateLoan(
		id,
		req,
	)
}

// ========================================
// MONEY COMPARISON
// ========================================

func moneyEqual(
	a float64,
	b float64,
) bool {

	const tolerance = 0.01

	if a > b {
		return a-b <= tolerance
	}

	return b-a <= tolerance
}
