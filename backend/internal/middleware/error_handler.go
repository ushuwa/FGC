package middleware

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v3"
)

func ErrorHandler(c fiber.Ctx, err error) error {

	statusCode := fiber.StatusInternalServerError

	if fiberErr := new(fiber.Error); errors.As(err, &fiberErr) {
		statusCode = fiberErr.Code
	}

	if statusCode >= 500 {
		log.Printf(
			"Internal server error: %v",
			err,
		)
	}

	message := "Internal server error"

	if statusCode < 500 {
		message = err.Error()
	}

	return c.Status(statusCode).JSON(fiber.Map{
		"success": false,
		"message": message,
	})
}
