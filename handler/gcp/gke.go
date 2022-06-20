package gcp

import (
	"context"
	"fmt"
	"log"
	"time"

	container "cloud.google.com/go/container/apiv1"
	"github.com/gofiber/fiber/v2"
	"github.com/kyg9823/gcp-resource-manager/config"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
)

type GkeNodeManagerParam struct {
	Replicas uint8 `query:"replicas"`
}

func GkeNodeManager(ctx *fiber.Ctx) error {
	projectId := ctx.Params("ProjectId")
	clusterId := ctx.Params("ClusterId")
	location := config.GetConfig("LOCATION")

	param := new(GkeNodeManagerParam)
	if err := ctx.QueryParser(param); err != nil {
		return err
	}

	log.Printf("Project ID: %s\n", projectId)
	log.Printf("Cluster ID: %s\n", clusterId)
	log.Printf("Replicas: %d\n", param.Replicas)

	c := context.Background()
	containerClient, err := container.NewClusterManagerClient(c)
	if err != nil {
		log.Printf("Fail to get Container Client")
		return err
	}
	defer containerClient.Close()

	req := &containerpb.ListNodePoolsRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s/clusters/%s", projectId, location, clusterId),
	}

	response, err := containerClient.ListNodePools(c, req)
	if err != nil {
		log.Printf("Fail to get NodePools")
		return err
	}

	log.Printf("Result: %v\n", response)

	for _, node := range response.NodePools {
		log.Printf("%v", node)
		nodeName := node.Name

		req := &containerpb.SetNodePoolSizeRequest{
			Name:      fmt.Sprintf("projects/%s/locations/%s/clusters/%s/nodePools/%s", projectId, location, clusterId, nodeName),
			NodeCount: 1,
		}

		response, err := containerClient.SetNodePoolSize(c, req)
		if err != nil {
			log.Printf("Fail to set NodePool Size")
			return err
		}

		operationId := response.Name
		waitCount := 1

		for {
			req := &containerpb.GetOperationRequest{
				Name: fmt.Sprintf("projects/%s/locations/%s/operations/%s", projectId, location, operationId),
			}

			response, err := containerClient.GetOperation(c, req)
			if err != nil {
				log.Printf("Fail to get Container Operation")
				return err
			}

			if response.Status == 3 {
				log.Printf("Done")
				break
			}
			time.Sleep(5 * time.Second)
			log.Printf("Wait %d seconds for resizing nodepool %s", (waitCount * 5), nodeName)
			waitCount += 1
		}
	}

	message := "OK"

	return ctx.JSON(fiber.Map{
		"status":  200,
		"message": message,
	})
}
