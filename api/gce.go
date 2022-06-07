package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	compute "cloud.google.com/go/compute/apiv1"
	"github.com/kyg9823/gcp-resource-manager/config"
	"github.com/kyg9823/gcp-resource-manager/types"
	"google.golang.org/api/iterator"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
	"google.golang.org/protobuf/proto"
)

type InstanceInfo struct {
	Project  string
	Zone     string
	Instance string
}

func GceManager(w http.ResponseWriter, r *http.Request) {
	projectId, _ := config.GetProjectId()

	action := r.URL.Query().Get("action")
	if action == "" {
		action = "stop"
	}

	ctx := context.Background()
	instanceClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		log.Printf("Fail to get Instance Client")
		return
	}
	defer instanceClient.Close()

	filter := "labels.auto-" + action + " = true"

	req := &computepb.AggregatedListInstancesRequest{
		Project:    projectId,
		MaxResults: proto.Uint32(10),
		Filter:     &filter,
	}

	it := instanceClient.AggregatedList(ctx, req)
	log.Printf("Instances found:\n")

	for {
		pair, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return
		}
		instances := pair.Value.Instances
		if len(instances) > 0 {
			for _, instance := range instances {
				if action == "start" {
					startInstance(ctx, instanceClient, &InstanceInfo{
						Project:  projectId,
						Zone:     instance.GetZone()[strings.LastIndex(instance.GetZone(), "/")+1:],
						Instance: instance.GetName(),
					})
				} else {
					stopInstance(ctx, instanceClient, &InstanceInfo{
						Project:  projectId,
						Zone:     instance.GetZone()[strings.LastIndex(instance.GetZone(), "/")+1:],
						Instance: instance.GetName(),
					})
				}
			}
		}
	}

	result := &types.Result{
		StatusCode: 200,
		Message:    "OK",
	}
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "Could not marshal JSON output", 500)
		return
	}
	fmt.Fprint(w)
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
