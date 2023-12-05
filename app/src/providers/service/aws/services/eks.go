package services

import (
	"cloudview/app/src/api/middleware/logger"
	"cloudview/app/src/helpers/constants"
	"cloudview/app/src/types"
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/eks"
)

/*Elastic Kubernetes Service*/

func (s Services) GetEKSData(respch chan<- types.GetDataOutput, wg *sync.WaitGroup) {
	defer wg.Done()
	instance := eks.NewFromConfig(s.Config)
	result, err := instance.ListClusters(context.TODO(), &eks.ListClustersInput{})
	if err != nil {
		logger.Logger.Error("eks.GetEKSData: ERROR", err)
		respch <- types.GetDataOutput{
			Name:   constants.EKS,
			Result: nil,
		}
		return
	}

	logger.Logger.Log("eks.GetEKSData: success")
	if len(result.Clusters) <= 0 {
		respch <- types.GetDataOutput{
			Name:   constants.EKS,
			Result: nil,
		}
		return
	}
	respch <- types.GetDataOutput{
		Name:   constants.EKS,
		Result: result,
	}
	return
}
