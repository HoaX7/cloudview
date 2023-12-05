package services

import (
	"cloudview/app/src/api/middleware/logger"
	"cloudview/app/src/helpers/constants"
	"cloudview/app/src/types"
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func (s Services) GetDynamodbTables(respch chan<- types.GetDataOutput, wg *sync.WaitGroup) {
	defer wg.Done()
	instance := dynamodb.NewFromConfig(s.Config)
	result, err := instance.ListTables(context.TODO(), &dynamodb.ListTablesInput{})
	if err != nil {
		logger.Logger.Error("dynamodb.GetDynamodbTables: ERROR", err)
		respch <- types.GetDataOutput{
			Name: constants.DYNAMODB,
		}
		return
	}
	logger.Logger.Log("dynamodb.GetDynamodbTables: success")
	if len(result.TableNames) <= 0 {
		respch <- types.GetDataOutput{
			Name:   constants.DYNAMODB,
			Result: nil,
		}
		return
	}
	respch <- types.GetDataOutput{
		Name:   constants.DYNAMODB,
		Result: result,
	}
	return
}
