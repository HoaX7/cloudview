package services

import (
	"cloudview/app/src/api/middleware/logger"
	"cloudview/app/src/helpers/constants"
	"cloudview/app/src/types"
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/efs"
)

/* Elastic File System */

func (s Services) GetEFSData(respch chan<- types.GetDataOutput, wg *sync.WaitGroup) {
	defer wg.Done()

	instance := efs.NewFromConfig(s.Config)
	result, err := instance.DescribeFileSystems(context.TODO(), &efs.DescribeFileSystemsInput{})
	if err != nil {
		logger.Logger.Error("efs.GetEFSData: ERROR", err)
		respch <- types.GetDataOutput{
			Name:   constants.EFS,
			Result: nil,
		}
		return
	}
	logger.Logger.Log("efs.GetEFSData: success")
	if len(result.FileSystems) <= 0 {
		respch <- types.GetDataOutput{
			Name:   constants.EFS,
			Result: nil,
		}
		return
	}
	respch <- types.GetDataOutput{
		Name:   constants.EFS,
		Result: result,
	}
	return
}
