package types

import "github.com/aws/aws-sdk-go-v2/service/cloudwatch"

type GetDataOutput struct {
	Name   string      `json:"name"`
	Result interface{} `json:"result"`
}

type DataService interface {
	Name() string // returns service name (aws, gcp, azure)
	GetData() ([]GetDataOutput, error)
	Init(accessKeyId string, accessKeySecret string, region string) error
}

type DataServiceReturnType struct {
	Services interface{} `json:"services"`
}

type Usage struct {
	Bandwidth *cloudwatch.GetMetricDataOutput
	Cpu       *cloudwatch.GetMetricDataOutput
}
