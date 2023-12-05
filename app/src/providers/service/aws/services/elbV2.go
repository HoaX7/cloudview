package services

import (
	"cloudview/app/src/api/middleware/logger"
	"cloudview/app/src/helpers/constants"
	"cloudview/app/src/types"
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
)

/*Elastic Loadbalancer*/

func (s Services) GetELBData(respch chan<- types.GetDataOutput, wg *sync.WaitGroup) {
	defer wg.Done()

	instance := elasticloadbalancingv2.NewFromConfig(s.Config)
	result, err := instance.DescribeLoadBalancers(context.TODO(), &elasticloadbalancingv2.DescribeLoadBalancersInput{})
	if err != nil {
		logger.Logger.Error("elbV2.GetELBData: ERROR", err)
		respch <- types.GetDataOutput{
			Name:   constants.ELBV2,
			Result: nil,
		}
		return
	}
	logger.Logger.Log("elbV2.GetELBData: success")
	if len(result.LoadBalancers) <= 0 {
		respch <- types.GetDataOutput{
			Name:   constants.ELBV2,
			Result: nil,
		}
		return
	}
	respch <- types.GetDataOutput{
		Name:   constants.ELBV2,
		Result: result,
	}
	return
}
