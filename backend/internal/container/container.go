package container

import (
	"github.com/jj.jobo/FGC/internal/handlers"
	"github.com/jj.jobo/FGC/internal/repositories"
	"github.com/jj.jobo/FGC/internal/services"
)

type Container struct {
	AuthHandler      *handlers.AuthHandler
	ClientHandler    *handlers.ClientHandler
	LoanHandler      *handlers.LoanHandler
	DashboardHandler *handlers.DashboardHandler
	PARHandler       *handlers.PARHandler
	ReportHandler    *handlers.ReportHandler
}

func BuildContainer() *Container {

	// =========================
	// REPOSITORIES
	// =========================

	userRepository :=
		repositories.NewUserRepository()

	clientRepository :=
		repositories.NewClientRepository()

	loanRepository :=
		repositories.NewLoanRepository()

	dashboardRepository :=
		repositories.NewDashboardRepository()

	parRepository :=
		repositories.NewPARRepository()

	reportRepository :=
		repositories.NewReportRepository()

	// =========================
	// SERVICES
	// =========================

	authService :=
		services.NewAuthService(
			userRepository,
		)

	clientService :=
		services.NewClientService(
			clientRepository,
		)

	loanService :=
		services.NewLoanService(
			loanRepository,
		)

	dashboardService :=
		services.NewDashboardService(
			dashboardRepository,
		)
	parService :=
		services.NewPARService(
			parRepository,
		)

	reportService :=
		services.NewReportService(
			reportRepository,
		)
	reportPDFService :=
		services.NewReportPDFService("./storage/reports")

	// =========================
	// HANDLERS
	// =========================

	authHandler :=
		handlers.NewAuthHandler(
			authService,
		)

	clientHandler :=
		handlers.NewClientHandler(
			clientService,
		)

	loanHandler :=
		handlers.NewLoanHandler(
			loanService,
		)

	dashboardHandler :=
		handlers.NewDashboardHandler(
			dashboardService,
		)

	parHandler :=
		handlers.NewPARHandler(
			parService,
		)

	reportHandler :=
		handlers.NewReportHandler(
			reportService,
			reportPDFService,
		)

	// =========================
	// CONTAINER
	// =========================

	return &Container{
		AuthHandler:      authHandler,
		ClientHandler:    clientHandler,
		LoanHandler:      loanHandler,
		DashboardHandler: dashboardHandler,
		PARHandler:       parHandler,
		ReportHandler:    reportHandler,
	}
}
