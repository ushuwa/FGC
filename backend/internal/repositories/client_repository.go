package repositories

import (
	"errors"
	"strings"

	"github.com/jj.jobo/FGC/internal/database"
	"github.com/jj.jobo/FGC/internal/dto"
	"github.com/jj.jobo/FGC/internal/models"
)

type ClientRepository struct {
	*BaseRepository[models.Client]
}

func NewClientRepository() *ClientRepository {
	return &ClientRepository{
		BaseRepository: NewBaseRepository[models.Client](
			database.DB,
		),
	}
}

func (r *ClientRepository) FindAll(
	search string,
) ([]models.Client, error) {

	var clients []models.Client

	query := database.DB.
		Model(&models.Client{})

	search = strings.TrimSpace(search)

	if search != "" {
		searchPattern := "%" + search + "%"

		query = query.Where(
			`first_name ILIKE ? 
			OR last_name ILIKE ?
			OR first_name || ' ' || last_name ILIKE ?
			OR contact_number ILIKE ?
			OR email ILIKE ?`,
			searchPattern,
			searchPattern,
			searchPattern,
			searchPattern,
			searchPattern,
		)
	}

	err := query.
		Order("last_name ASC").
		Order("first_name ASC").
		Order("id ASC").
		Find(&clients).Error

	if err != nil {
		return nil, err
	}

	return clients, nil
}

func (r *ClientRepository) FindByID(
	id uint,
) (*models.Client, error) {

	var client models.Client

	err := database.DB.
		First(&client, id).
		Error

	if err != nil {
		return nil, err
	}

	return &client, nil
}

func (r *ClientRepository) CreateClient(
	client *models.Client,
) error {

	return database.DB.
		Create(client).
		Error
}

func (r *ClientRepository) UpdateClient(
	client *models.Client,
) error {

	return database.DB.
		Model(&models.Client{}).
		Where("id = ?", client.ID).
		Updates(map[string]interface{}{
			"first_name":      client.FirstName,
			"last_name":       client.LastName,
			"contact_number":  client.ContactNumber,
			"email":           client.Email,
			"current_address": client.CurrentAddress,
		}).
		Error
}

func (r *ClientRepository) DeleteClient(
	id uint,
) error {

	return database.DB.
		Delete(
			&models.Client{},
			id,
		).
		Error
}

func (r *ClientRepository) GetProfile(
	id int,
) (*dto.ClientProfileResponse, error) {

	var client dto.ClientProfileClient

	err := database.DB.Raw(`
		SELECT
			id,
			first_name,
			last_name,
			contact_number,
			email,
			current_address
		FROM clients
		WHERE id = ?
	`, id).Scan(&client).Error

	if err != nil {
		return nil, err
	}

	if client.ID == 0 {
		return nil, errors.New(
			"client not found",
		)
	}

	var summary dto.ClientSummary

	err = database.DB.Raw(`
		SELECT
			COUNT(l.id) AS total_loans,

			COUNT(
				CASE
					WHEN l.status = 'ACTIVE'
					THEN 1
				END
			) AS active_loans,

			COALESCE(
				SUM(l.principal_amount),
				0
			) AS total_principal,

			COALESCE(
				(
					SELECT SUM(p.amount_paid)
					FROM payments p
					INNER JOIN loans lp
						ON lp.id = p.loan_id
					WHERE lp.client_id = ?
				),
				0
			) AS total_paid,

			COALESCE(
				SUM(l.pn_value),
				0
			)
			-
			COALESCE(
				(
					SELECT SUM(p.amount_paid)
					FROM payments p
					INNER JOIN loans lp
						ON lp.id = p.loan_id
					WHERE lp.client_id = ?
				),
				0
			) AS total_outstanding

		FROM loans l
		WHERE l.client_id = ?
	`,
		id,
		id,
		id,
	).Scan(&summary).Error

	if err != nil {
		return nil, err
	}

	var loans []dto.ClientLoanSummary

	err = database.DB.Raw(`
		SELECT
			l.id,
			l.pn_number,
			l.loan_type,
			l.principal_amount,
			l.interest_rate,
			l.loan_interest,
			l.pn_value,
			l.loan_term,
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
				SUM(p.amount_paid),
				0
			) AS total_paid,

			l.pn_value -
			COALESCE(
				SUM(p.amount_paid),
				0
			) AS outstanding_balance

		FROM loans l

		LEFT JOIN payments p
			ON p.loan_id = l.id

		WHERE l.client_id = ?

		GROUP BY
			l.id,
			l.pn_number,
			l.loan_type,
			l.principal_amount,
			l.interest_rate,
			l.loan_interest,
			l.pn_value,
			l.loan_term,
			l.amortization_amount,
			l.disbursement_date,
			l.maturity_date,
			l.status

		ORDER BY
			l.id DESC
	`,
		id,
	).Scan(&loans).Error

	if err != nil {
		return nil, err
	}

	if loans == nil {
		loans = []dto.ClientLoanSummary{}
	}

	var payments []dto.ClientPayment

	err = database.DB.Raw(`
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

		INNER JOIN loans l
			ON l.id = p.loan_id

		WHERE l.client_id = ?

		ORDER BY
			p.payment_date ASC,
			p.id ASC
	`,
		id,
	).Scan(&payments).Error

	if err != nil {
		return nil, err
	}

	if payments == nil {
		payments = []dto.ClientPayment{}
	}

	return &dto.ClientProfileResponse{
		Client:   client,
		Summary:  summary,
		Loans:    loans,
		Payments: payments,
	}, nil
}
