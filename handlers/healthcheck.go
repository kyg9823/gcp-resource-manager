package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func Healthcheck(ctx *fiber.Ctx) error {

	return ctx.JSON(fiber.Map{
		"status":  200,
		"message": "OK",
	})
}
