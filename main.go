package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kyg9823/gcp-resource-manager/config"
	"github.com/kyg9823/gcp-resource-manager/handler"
	"github.com/kyg9823/gcp-resource-manager/router"
	"github.com/kyg9823/gcp-resource-manager/utils"
)

func main() {

	app := fiber.New()

	app.Get("/healthcheck", handler.Healthcheck)

	router.SetupRouters(app)

	profile := config.GetConfig("PROFILE")
	if profile == "dev" {
		utils.StartServer(app)
	} else {
		utils.StartServerWithGracefulShutdown(app)
	}

	// http.HandleFunc("/api/v1/healthcheck", api.Healthcheck)
	// http.HandleFunc("/api/v1/gcp/:ProjectID/gce", gcp.GceManager)
	// port := config.GetPort()
	// if port == "" {
	// 	port = "8080"
	// }

	// log.Printf("Listening on port %s", port)
	// if err := http.ListenAndServe(":"+port, nil); err != nil {
	// 	log.Fatal(err)
	// }
}
