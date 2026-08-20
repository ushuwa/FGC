package services

import (
	"github.com/jj.jobo/FGC/internal/dto"
	"github.com/jj.jobo/FGC/internal/repositories"
)

type DashboardService struct {
	dashboardRepository *repositories.DashboardRepository
}

func NewDashboardService(
	dashboardRepository *repositories.DashboardRepository,
) *DashboardService {
	return &DashboardService{
		dashboardRepository: dashboardRepository,
	}
}

func (s *DashboardService) GetSummary() (
	*dto.DashboardSummary,
	error,
) {
	return s.dashboardRepository.GetSummary()
}
