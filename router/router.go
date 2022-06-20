package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/kyg9823/gcp-resource-manager/handler/aws"
	"github.com/kyg9823/gcp-resource-manager/handler/gcp"
)

func SetupRouters(app *fiber.App) {
	api := app.Group("/api/v1", logger.New())

	api.Get("/gcp/:ProjectId/gce/state", gcp.GceStateManager)
	api.Get("/gcp/:ProjectId/gke/:ClusterId/node", gcp.GkeNodeManager)

	api.Get("/aws/ec2", aws.EC2StateManager)
	api.Get("/aws/autoscaling", aws.AutoScalingManager)
}
