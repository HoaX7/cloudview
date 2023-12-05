package services

import (
	"cloudview/app/src/api/middleware/logger"
	"cloudview/app/src/helpers/constants"
	"cloudview/app/src/types"
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/route53"
)

func (s Services) GetRoute53Data(respch chan<- types.GetDataOutput, wg *sync.WaitGroup) {
	defer wg.Done()

	instance := route53.NewFromConfig(s.Config)
	result, err := instance.ListHostedZones(context.TODO(), &route53.ListHostedZonesInput{})
	if err != nil {
		logger.Logger.Error("route53.GeetRoute53Data: ERROR", err)
		respch <- types.GetDataOutput{
			Name:   constants.ROUTE53,
			Result: nil,
		}
		return
	}
	logger.Logger.Log("route53.GetRoute53Data: success")
	if len(result.HostedZones) <= 0 {
		respch <- types.GetDataOutput{
			Name:   constants.ROUTE53,
			Result: nil,
		}
		return
	}
	respch <- types.GetDataOutput{
		Name:   constants.ROUTE53,
		Result: result,
	}
	return
}
