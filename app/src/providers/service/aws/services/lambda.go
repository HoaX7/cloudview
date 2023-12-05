package services

import (
	"cloudview/app/src/api/middleware/logger"
	"cloudview/app/src/helpers/constants"
	"cloudview/app/src/types"
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

func (s Services) GetLambdaFuntions(respch chan<- types.GetDataOutput, wg *sync.WaitGroup) {
	defer wg.Done()

	instance := lambda.NewFromConfig(s.Config)
	result, err := instance.ListFunctions(context.TODO(), &lambda.ListFunctionsInput{})
	if err != nil {
		logger.Logger.Error("lambda.GetLambdaFunctions: ERROR", err)
		respch <- types.GetDataOutput{
			Name:   constants.LAMBDA,
			Result: nil,
		}
		return
	}
	logger.Logger.Log("lambda.GetLambdaFunctions: success")
	if len(result.Functions) <= 0 {
		respch <- types.GetDataOutput{
			Name:   constants.LAMBDA,
			Result: nil,
		}
		return
	}
	respch <- types.GetDataOutput{
		Name:   constants.LAMBDA,
		Result: result,
	}
	return
}
