package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/kyg9823/gcp-resource-manager/handler"
)

func SetupRouters(app *fiber.App) {
	api := app.Group("/api/v1", logger.New())

	api.Get("/gcp/:ProjectId/gce/state", handler.GceStateManager)

}
