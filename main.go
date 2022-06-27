package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kyg9823/gcp-resource-manager/config"
	"github.com/kyg9823/gcp-resource-manager/handlers"
	"github.com/kyg9823/gcp-resource-manager/router"
	"github.com/kyg9823/gcp-resource-manager/utils"
)

func main() {

	app := fiber.New(fiber.Config{
		ErrorHandler: handlers.DefaultErrorHandler,
	})

	app.Get("/healthcheck", handlers.Healthcheck)

	router.SetupRouters(app)

	profile := config.GetConfig("PROFILE")
	if profile == "dev" {
		utils.StartServer(app)
	} else {
		utils.StartServerWithGracefulShutdown(app)
	}
}
