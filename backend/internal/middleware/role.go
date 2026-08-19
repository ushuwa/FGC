package middleware

import "github.com/gofiber/fiber/v3"

func RoleMiddleware(roles ...string) fiber.Handler {

	return func(c fiber.Ctx) error {

		roleValue := c.Locals("role")

		role, ok := roleValue.(string)

		if !ok || role == "" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Access denied",
			})
		}

		for _, allowedRole := range roles {

			if role == allowedRole {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "Access denied",
		})
	}
}
