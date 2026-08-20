package repositories

import "errors"

var (
	ErrInvalidReportType = errors.New(
		"invalid report type",
	)

	ErrLoanIDRequired = errors.New(
		"loan id is required for amortization report",
	)
)
