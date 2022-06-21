package gcp

import (
	"context"
	"log"
	"strings"

	compute "cloud.google.com/go/compute/apiv1"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/iterator"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
	"google.golang.org/protobuf/proto"
)

type InstanceInfo struct {
	Project           string
	Zone              string
	Instance          string
	InstanceGroupName string
	Size              int32
}

type GceManagerParam struct {
	Action string `query:"action"`
}

func GceStateManager(ctx *fiber.Ctx) error {
	projectId := ctx.Params("ProjectId")

	param := new(GceManagerParam)
	if err := ctx.QueryParser(param); err != nil {
		return err
	}

	if param.Action == "" {
		param.Action = "stop"
	}

	c := context.Background()
	instanceClient, err := compute.NewInstancesRESTClient(c)
	if err != nil {
		log.Printf("Fail to get Instance Client")
		return err
	}
	defer instanceClient.Close()

	instanceGroupManager, err := compute.NewInstanceGroupManagersRESTClient(c)
	if err != nil {
		log.Printf("Fail to get Instance Group Manager Client")
		return err
	}
	defer instanceGroupManager.Close()

	filter := "labels.auto-" + param.Action + " = true"

	req := &computepb.AggregatedListInstancesRequest{
		Project:    projectId,
		MaxResults: proto.Uint32(10),
		Filter:     &filter,
	}

	it := instanceClient.AggregatedList(c, req)
	log.Printf("Instances found:\n")

	for {
		pair, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		instances := pair.Value.Instances
		if len(instances) > 0 {
			for _, instance := range instances {
				instanceGroup := getMetadata(instance.Metadata.Items, "created-by")
				if instanceGroup != "" {
					instanceInfo := &InstanceInfo{
						Project:           projectId,
						Zone:              instance.GetZone()[strings.LastIndex(instance.GetZone(), "/")+1:],
						Instance:          instance.GetName(),
						InstanceGroupName: instanceGroup[strings.LastIndex(instanceGroup, "/")+1:],
						Size:              0,
					}
					if param.Action == "start" {
						instanceInfo.Size = 1
					}
					resizeInstanceGroup(c, instanceGroupManager, instanceInfo)
				} else {
					if param.Action == "start" {
						startInstance(c, instanceClient, &InstanceInfo{
							Project:  projectId,
							Zone:     instance.GetZone()[strings.LastIndex(instance.GetZone(), "/")+1:],
							Instance: instance.GetName(),
						})
					} else {
						stopInstance(c, instanceClient, &InstanceInfo{
							Project:  projectId,
							Zone:     instance.GetZone()[strings.LastIndex(instance.GetZone(), "/")+1:],
							Instance: instance.GetName(),
						})
					}
				}
			}
		}
	}

	message := "OK"

	return ctx.JSON(fiber.Map{
		"status":  200,
		"message": message,
	})
}

func getMetadata(items []*computepb.Items, key string) string {
	result := ""

	for _, item := range items {
		if *item.Key == key {
			result = *item.Value
			break
		}
	}

	return result
}

func startInstance(ctx context.Context, instanceClient *compute.InstancesClient, instanceInfo *InstanceInfo) {
	req := &computepb.StartInstanceRequest{
		Project:  instanceInfo.Project,
		Zone:     instanceInfo.Zone,
		Instance: instanceInfo.Instance,
	}

	op, err := instanceClient.Start(ctx, req)
	if err != nil {
		log.Printf("Failed to start instance\n")
	}

	if err = op.Wait(ctx); err != nil {
		log.Printf("unable to wait for the operation: %v\n", err)
	}

	log.Printf("Instance started\n")
}

func stopInstance(ctx context.Context, instanceClient *compute.InstancesClient, instanceInfo *InstanceInfo) {
	req := &computepb.StopInstanceRequest{
		Project:  instanceInfo.Project,
		Zone:     instanceInfo.Zone,
		Instance: instanceInfo.Instance,
	}

	op, err := instanceClient.Stop(ctx, req)
	if err != nil {
		log.Printf("Failed to start instance\n")
	}

	if err = op.Wait(ctx); err != nil {
		log.Printf("unable to wait for the operation: %v\n", err)
	}

	log.Printf("Instance stopped\n")
}

func resizeInstanceGroup(ctx context.Context, instanceGroupManager *compute.InstanceGroupManagersClient, instanceInfo *InstanceInfo) {
	req := &computepb.ResizeInstanceGroupManagerRequest{
		InstanceGroupManager: instanceInfo.InstanceGroupName,
		Project:              instanceInfo.Project,
		Zone:                 instanceInfo.Zone,
		Size:                 0,
	}

	op, err := instanceGroupManager.Resize(ctx, req)
	if err != nil {
		log.Printf("Failed to resize Instance Group \n")
	}

	if err = op.Wait(ctx); err != nil {
		log.Printf("unable to wait for the operation: %v\n", err)
	}

	log.Printf("Resize Complete\n")
	log.Printf("%v", instanceInfo)

}
