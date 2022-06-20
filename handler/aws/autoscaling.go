package aws

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/gofiber/fiber/v2"
)

type AutoscalingManagerParam struct {
	Action string   `query:"action"`
	Tags   []string `query:"tags"`
}

func AutoScalingManager(ctx *fiber.Ctx) error {

	param := new(AutoscalingManagerParam)
	if err := ctx.QueryParser(param); err != nil {
		log.Printf("Fail to parse Querystring")
		return err
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("Fail to load default config")
		return err
	}

	client := autoscaling.NewFromConfig(cfg)

	input := &autoscaling.DescribeAutoScalingGroupsInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("tag:AutoStop"),
				Values: []string{"True"},
			},
		},
	}

	result, err := client.DescribeAutoScalingGroups(context.TODO(), input)
	if err != nil {
		log.Printf("Fail to retrieve autoscaling")
		return err
	}

	for _, asg := range result.AutoScalingGroups {
		log.Printf("ASG ID: %s", *asg.AutoScalingGroupName)
	}

	// updateAutoScalingGroupInput := &autoscaling.UpdateAutoScalingGroupInput{
	// 	AutoScalingGroupName: aws.String("WN-HYB006-JOB22000199-SCH202206101127001"),
	// 	MinSize:              aws.Int32(0),
	// 	DesiredCapacity:      aws.Int32(0),
	// 	MaxSize:              aws.Int32(0),
	// }

	// asgResult, err := client.UpdateAutoScalingGroup(context.TODO(), updateAutoScalingGroupInput)
	// if err != nil {
	// 	log.Printf("Fail to resize asg: %v", err)
	// 	return err
	// }
	// log.Printf("%v", asgResult)

	message := "OK"
	return ctx.JSON(fiber.Map{
		"status":  200,
		"message": message,
	})
}
