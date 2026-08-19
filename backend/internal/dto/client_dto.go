package dto

type CreateClientRequest struct {
	FirstName      string  `json:"first_name" validate:"required"`
	LastName       string  `json:"last_name" validate:"required"`
	ContactNumber  *string `json:"contact_number"`
	Email          *string `json:"email"`
	CurrentAddress *string `json:"current_address"`
}

type UpdateClientRequest struct {
	FirstName      string  `json:"first_name" validate:"required"`
	LastName       string  `json:"last_name" validate:"required"`
	ContactNumber  *string `json:"contact_number"`
	Email          *string `json:"email"`
	CurrentAddress *string `json:"current_address"`
}
