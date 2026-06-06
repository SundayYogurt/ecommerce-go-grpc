package utils

import "github.com/gofiber/fiber/v2"

func RespondError(ctx *fiber.Ctx, status int, msg string) error {
	return ctx.Status(status).JSON(fiber.Map{"error": msg})
}

// create a generic response function for success
func RespondSuccess(ctx *fiber.Ctx, status int, data interface{}) error {
	return ctx.Status(status).JSON(fiber.Map{"data": data})
}
