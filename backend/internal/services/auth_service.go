package services

import (
	"errors"

	"github.com/jj.jobo/FGC/internal/dto"
	repoiface "github.com/jj.jobo/FGC/internal/repositories/interfaces"
	"github.com/jj.jobo/FGC/internal/utils"
)

type AuthService struct {
	userRepo repoiface.UserRepository
}

func NewAuthService(userRepo repoiface.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.FindByUsername(req.Username)

	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	if !user.IsActive {
		return nil, errors.New("user account is inactive")
	}

	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid username or password")
	}

	token, err := utils.GenerateJWT(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
		User: dto.UserDTO{
			ID:       user.ID,
			Username: user.Username,
			FullName: user.FullName,
			Role:     user.Role,
		},
	}, nil
}
