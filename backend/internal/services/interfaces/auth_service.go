package interfaces

import "github.com/jj.jobo/FGC/internal/dto"

type AuthService interface {
	Login(req dto.LoginRequest) (*dto.LoginResponse, error)
}
