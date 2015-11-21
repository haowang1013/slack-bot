package commands

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/nlopes/slack"
	"strings"
)

const (
	region = "ap-southeast-1"
)

var awsConfig = aws.Config{Region: aws.String(region)}

func init() {
	rootHandler.addHandler("aws", &handleCollection{
		"ec2": &handleCollection{
			"list": HandlerFunc(ec2ListHandler),
		},
	})
}

func getInstanceName(inst *ec2.Instance) string {
	for _, tag := range inst.Tags {
		if *tag.Key == "Name" {
			return *tag.Value
		}
	}

	return "???"
}

func ec2ListHandler(fields []string, m *slack.MessageEvent) error {
	svc := ec2.New(session.New(), &awsConfig)
	resp, err := svc.DescribeInstances(nil)
	if err != nil {
		return err
	}

	lines := []string{}
	index := 0
	for _, res := range resp.Reservations {
		for _, inst := range res.Instances {
			lines = append(lines, fmt.Sprintf("%d: %s %s %s %s",
				index,
				*inst.InstanceId,
				getInstanceName(inst),
				*inst.InstanceType,
				*inst.State.Name,
			))
			index++
		}
	}
	sendMessage(strings.Join(lines, "\n"), m.Channel)
	return nil
}
