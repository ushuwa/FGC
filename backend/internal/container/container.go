package container

import (
	"github.com/jj.jobo/FGC/internal/handlers"
	"github.com/jj.jobo/FGC/internal/repositories"
	"github.com/jj.jobo/FGC/internal/services"
)

type Container struct {
	AuthHandler   *handlers.AuthHandler
	ClientHandler *handlers.ClientHandler
	LoanHandler   *handlers.LoanHandler
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

	// =========================
	// CONTAINER
	// =========================

	return &Container{
		AuthHandler:   authHandler,
		ClientHandler: clientHandler,
		LoanHandler:   loanHandler,
	}
}
