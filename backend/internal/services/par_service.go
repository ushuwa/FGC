package services

import (
	"github.com/jj.jobo/FGC/internal/dto"
	"github.com/jj.jobo/FGC/internal/repositories"
)

type PARService struct {
	parRepository *repositories.PARRepository
}

func NewPARService(
	parRepository *repositories.PARRepository,
) *PARService {
	return &PARService{
		parRepository: parRepository,
	}
}

func (s *PARService) GetPAR(
	search string,
	status string,
	aging string,
	page int,
	limit int,
) (*dto.PARResponse, error) {

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

	// ========================================
	// GET PAR
	// ========================================

	return s.parRepository.GetPAR(
		search,
		status,
		aging,
		page,
		limit,
	)
}
