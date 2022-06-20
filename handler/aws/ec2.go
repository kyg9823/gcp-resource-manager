package aws

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type EC2ManagerParam struct {
	Action string `query:"action"`
}

func EC2StateManager(ctx *fiber.Ctx) error {

	param := new(EC2ManagerParam)
	if err := ctx.QueryParser(param); err != nil {
		return err
	}

	if param.Action == "" {
		param.Action = "stop"
	}

	caser := cases.Title(language.English)
	tagKey := "tag:Auto" + caser.String(param.Action)

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("Fail to load default config")
		return err
	}

	client := ec2.NewFromConfig(cfg)

	input := &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   &tagKey,
				Values: []string{"True"},
			},
		},
	}

	result, err := client.DescribeInstances(context.TODO(), input)

	if err != nil {
		log.Printf("Fail to retrieve instances")
		return err
	}

	instanceIds := []string{}

	for _, r := range result.Reservations {
		log.Printf("Reservation ID: %s", *r.ReservationId)
		log.Printf("Instance IDs: ")

		for _, instance := range r.Instances {
			log.Printf("  %v", instance.InstanceId)
			instanceIds = append(instanceIds, *instance.InstanceId)
		}
		log.Printf("\n")
	}

	log.Printf("%v", instanceIds)

	if param.Action == "start" {
		instanceParam := &ec2.StartInstancesInput{
			InstanceIds: instanceIds,
		}

		instanceResult, err := client.StartInstances(context.TODO(), instanceParam)
		if err != nil {
			log.Printf("Fail to stop Instances")
			return err
		}
		log.Printf("%v", instanceResult)
	} else {
		instanceParam := &ec2.StopInstancesInput{
			InstanceIds: instanceIds,
		}

		instanceResult, err := client.StopInstances(context.TODO(), instanceParam)
		if err != nil {
			log.Printf("Fail to stop Instances")
			return err
		}
		log.Printf("%v", instanceResult)
	}

	message := "OK"

	return ctx.JSON(fiber.Map{
		"status":  200,
		"message": message,
	})
}
