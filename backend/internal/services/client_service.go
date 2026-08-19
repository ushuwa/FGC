package services

import (
	"errors"
	"strings"

	"github.com/jj.jobo/FGC/internal/dto"
	"github.com/jj.jobo/FGC/internal/models"
	"github.com/jj.jobo/FGC/internal/repositories"
	"gorm.io/gorm"
)

type ClientService struct {
	clientRepository *repositories.ClientRepository
}

func NewClientService(
	clientRepository *repositories.ClientRepository,
) *ClientService {

	return &ClientService{
		clientRepository: clientRepository,
	}
}

func (s *ClientService) GetClients(
	search string,
) ([]models.Client, error) {

	return s.clientRepository.FindAll(
		search,
	)
}

func (s *ClientService) GetClient(
	id uint,
) (*models.Client, error) {

	client, err :=
		s.clientRepository.FindByID(id)

	if err != nil {

		if errors.Is(
			err,
			gorm.ErrRecordNotFound,
		) {
			return nil, errors.New(
				"client not found",
			)
		}

		return nil, err
	}

	return client, nil
}

func (s *ClientService) CreateClient(
	req dto.CreateClientRequest,
) (*models.Client, error) {

	firstName :=
		strings.TrimSpace(
			req.FirstName,
		)

	lastName :=
		strings.TrimSpace(
			req.LastName,
		)

	if firstName == "" {
		return nil, errors.New(
			"first name is required",
		)
	}

	if lastName == "" {
		return nil, errors.New(
			"last name is required",
		)
	}

	client := &models.Client{
		FirstName: firstName,
		LastName:  lastName,
		ContactNumber: cleanOptional(
			req.ContactNumber,
		),
		Email: cleanOptional(
			req.Email,
		),
		CurrentAddress: cleanOptional(
			req.CurrentAddress,
		),
	}

	err := s.clientRepository.CreateClient(
		client,
	)

	if err != nil {
		return nil, err
	}

	return client, nil
}

func (s *ClientService) UpdateClient(
	id uint,
	req dto.UpdateClientRequest,
) (*models.Client, error) {

	firstName :=
		strings.TrimSpace(
			req.FirstName,
		)

	lastName :=
		strings.TrimSpace(
			req.LastName,
		)

	if firstName == "" {
		return nil, errors.New(
			"first name is required",
		)
	}

	if lastName == "" {
		return nil, errors.New(
			"last name is required",
		)
	}

	client := &models.Client{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		ContactNumber: cleanOptional(
			req.ContactNumber,
		),
		Email: cleanOptional(
			req.Email,
		),
		CurrentAddress: cleanOptional(
			req.CurrentAddress,
		),
	}

	err := s.clientRepository.UpdateClient(
		client,
	)

	if err != nil {
		return nil, err
	}

	return s.GetClient(id)
}

func (s *ClientService) DeleteClient(
	id uint,
) error {

	_, err := s.GetClient(id)

	if err != nil {
		return err
	}

	return s.clientRepository.DeleteClient(
		id,
	)
}

func cleanOptional(
	value *string,
) *string {

	if value == nil {
		return nil
	}

	cleaned :=
		strings.TrimSpace(
			*value,
		)

	if cleaned == "" {
		return nil
	}

	return &cleaned
}

func (s *ClientService) GetClientProfile(
	id int,
) (*dto.ClientProfileResponse, error) {

	if id <= 0 {
		return nil, errors.New(
			"invalid client id",
		)
	}

	return s.clientRepository.GetProfile(
		id,
	)
}
