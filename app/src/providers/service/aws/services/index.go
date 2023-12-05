package services

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
)

type Services struct {
	Config     aws.Config
	Region     string
	CloudWatch *cloudwatch.Client
}
