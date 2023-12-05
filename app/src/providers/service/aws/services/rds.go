package services

import (
	"cloudview/app/src/api/middleware/logger"
	"cloudview/app/src/helpers/constants"
	"cloudview/app/src/types"
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/rds"
)

func (s Services) GetRDSInstances(respch chan<- types.GetDataOutput, wg *sync.WaitGroup) {
	defer wg.Done()
	instance := rds.NewFromConfig(s.Config)
	result, err := instance.DescribeDBInstances(context.TODO(), &rds.DescribeDBInstancesInput{})
	if err != nil {
		logger.Logger.Error("s3.GetS3Buckets: ERROR", err)
		respch <- types.GetDataOutput{
			Name:   constants.RDS,
			Result: nil,
		}
		return
	}

	logger.Logger.Log("rds.GetRDSInstances: success")
	if len(result.DBInstances) <= 0 {
		respch <- types.GetDataOutput{
			Name:   constants.RDS,
			Result: nil,
		}
		return
	}
	respch <- types.GetDataOutput{
		Name:   constants.RDS,
		Result: result,
	}
	return
}
