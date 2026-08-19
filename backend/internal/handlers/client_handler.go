package handlers

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"

	"github.com/jj.jobo/FGC/internal/dto"
	"github.com/jj.jobo/FGC/internal/services"
	"github.com/jj.jobo/FGC/internal/utils"
)

type ClientHandler struct {
	service *services.ClientService
}

func NewClientHandler(
	service *services.ClientService,
) *ClientHandler {

	return &ClientHandler{
		service: service,
	}
}

func (h *ClientHandler) GetClients(
	c fiber.Ctx,
) error {

	search := c.Query("search")

	clients, err :=
		h.service.GetClients(search)

	if err != nil {
		return utils.Error(
			c,
			fiber.StatusInternalServerError,
			"Failed to retrieve clients",
		)
	}

	return utils.Success(
		c,
		"Clients retrieved successfully",
		clients,
	)
}

func (h *ClientHandler) GetClient(
	c fiber.Ctx,
) error {

	id, err := strconv.ParseUint(
		c.Params("id"),
		10,
		64,
	)

	if err != nil || id == 0 {
		return utils.Error(
			c,
			fiber.StatusBadRequest,
			"Invalid client ID",
		)
	}

	client, err :=
		h.service.GetClient(uint(id))

	if err != nil {

		if errors.Is(
			err,
			gorm.ErrRecordNotFound,
		) {
			return utils.Error(
				c,
				fiber.StatusNotFound,
				"Client not found",
			)
		}

		if err.Error() ==
			"client not found" {
			return utils.Error(
				c,
				fiber.StatusNotFound,
				"Client not found",
			)
		}

		return utils.Error(
			c,
			fiber.StatusInternalServerError,
			"Failed to retrieve client",
		)
	}

	return utils.Success(
		c,
		"Client retrieved successfully",
		client,
	)
}

func (h *ClientHandler) CreateClient(
	c fiber.Ctx,
) error {

	var req dto.CreateClientRequest

	if err := c.Bind().Body(&req); err != nil {
		return utils.Error(
			c,
			fiber.StatusBadRequest,
			"Invalid request body",
		)
	}

	if err := validate.Struct(req); err != nil {
		return utils.Error(
			c,
			fiber.StatusBadRequest,
			err.Error(),
		)
	}

	client, err :=
		h.service.CreateClient(req)

	if err != nil {
		return utils.Error(
			c,
			fiber.StatusBadRequest,
			err.Error(),
		)
	}

	return c.Status(
		fiber.StatusCreated,
	).JSON(fiber.Map{
		"success": true,
		"message": "Client created successfully",
		"data":    client,
	})
}

func (h *ClientHandler) UpdateClient(
	c fiber.Ctx,
) error {

	id, err := strconv.ParseUint(
		c.Params("id"),
		10,
		64,
	)

	if err != nil || id == 0 {
		return utils.Error(
			c,
			fiber.StatusBadRequest,
			"Invalid client ID",
		)
	}

	var req dto.UpdateClientRequest

	if err := c.Bind().Body(&req); err != nil {
		return utils.Error(
			c,
			fiber.StatusBadRequest,
			"Invalid request body",
		)
	}

	if err := validate.Struct(req); err != nil {
		return utils.Error(
			c,
			fiber.StatusBadRequest,
			err.Error(),
		)
	}

	client, err :=
		h.service.UpdateClient(
			uint(id),
			req,
		)

	if err != nil {

		if err.Error() ==
			"record not found" {
			return utils.Error(
				c,
				fiber.StatusNotFound,
				"Client not found",
			)
		}

		return utils.Error(
			c,
			fiber.StatusBadRequest,
			err.Error(),
		)
	}

	return utils.Success(
		c,
		"Client updated successfully",
		client,
	)
}

func (h *ClientHandler) DeleteClient(
	c fiber.Ctx,
) error {

	id, err := strconv.ParseUint(
		c.Params("id"),
		10,
		64,
	)

	if err != nil || id == 0 {
		return utils.Error(
			c,
			fiber.StatusBadRequest,
			"Invalid client ID",
		)
	}

	if err := h.service.DeleteClient(
		uint(id),
	); err != nil {

		if err.Error() ==
			"client not found" {
			return utils.Error(
				c,
				fiber.StatusNotFound,
				"Client not found",
			)
		}

		return utils.Error(
			c,
			fiber.StatusInternalServerError,
			"Failed to delete client",
		)
	}

	return utils.Success(
		c,
		"Client deleted successfully",
		nil,
	)
}

func (h *ClientHandler) GetClientProfile(
	c fiber.Ctx,
) error {

	id, err := strconv.Atoi(
		c.Params("id"),
	)

	if err != nil || id <= 0 {
		return utils.Error(
			c,
			fiber.StatusBadRequest,
			"Invalid client ID",
		)
	}

	profile, err :=
		h.service.GetClientProfile(id)

	if err != nil {

		if err.Error() ==
			"client not found" {
			return utils.Error(
				c,
				fiber.StatusNotFound,
				"Client not found",
			)
		}

		return utils.Error(
			c,
			fiber.StatusInternalServerError,
			"Failed to retrieve client profile",
		)
	}

	return utils.Success(
		c,
		"Client profile retrieved successfully",
		profile,
	)
}
