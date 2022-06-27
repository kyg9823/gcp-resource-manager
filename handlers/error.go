package handlers

import "github.com/gofiber/fiber/v2"

func DefaultErrorHandler(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	return ctx.Status(code).JSON(fiber.Map{
		"status": code,
		"result": err.Error(),
	})
}
