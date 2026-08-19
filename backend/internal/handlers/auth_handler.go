package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"

	"github.com/jj.jobo/FGC/internal/dto"
	serviceiface "github.com/jj.jobo/FGC/internal/services/interfaces"
	"github.com/jj.jobo/FGC/internal/utils"
)

var validate = validator.New()

type AuthHandler struct {
	service serviceiface.AuthService
}

func NewAuthHandler(
	service serviceiface.AuthService,
) *AuthHandler {

	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) Login(c fiber.Ctx) error {

	var req dto.LoginRequest

	if err := c.Bind().Body(&req); err != nil {
		return utils.Error(
			c,
			fiber.StatusBadRequest,
			"Invalid request body",
		)
	}

	if err := validate.Struct(req); err != nil {

		return utils.ValidationError(
			c,
			err.Error(),
		)
	}

	res, err := h.service.Login(req)

	if err != nil {
		return utils.Error(
			c,
			fiber.StatusUnauthorized,
			err.Error(),
		)
	}

	return utils.Success(
		c,
		"Login successful",
		res,
	)
}
