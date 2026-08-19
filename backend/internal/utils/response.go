package utils

import "github.com/gofiber/fiber/v3"

type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

func Success(
	c fiber.Ctx,
	message string,
	data any,
) error {

	return c.JSON(APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(
	c fiber.Ctx,
	status int,
	message string,
) error {

	return c.Status(status).JSON(APIResponse{
		Success: false,
		Message: message,
	})
}

func ValidationError(
	c fiber.Ctx,
	errors any,
) error {

	return c.Status(fiber.StatusBadRequest).JSON(APIResponse{
		Success: false,
		Message: "Validation failed",
		Errors:  errors,
	})
}
