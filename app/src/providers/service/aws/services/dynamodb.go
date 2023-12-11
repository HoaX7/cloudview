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
	swg := &sync.WaitGroup{}
	ch := make(chan *dynamodb.DescribeTableOutput, len(result.TableNames))
	for _, tableName := range result.TableNames {
		swg.Add(1)
		go s.describeTable(ch, swg, &tableName)
	}
	go func() {
		swg.Wait()
		close(ch)
	}()

	res := []*dynamodb.DescribeTableOutput{}
	for data := range ch {
		if data != nil {
			res = append(res, data)
		}
	}
	respch <- types.GetDataOutput{
		Name:   constants.DYNAMODB,
		Result: res,
	}
	return
}

func (s Services) describeTable(respch chan<- *dynamodb.DescribeTableOutput, wg *sync.WaitGroup, tableName *string) {
	defer wg.Done()
	logger.Logger.Log("fetching table details for tableName:", tableName)
	instance := dynamodb.NewFromConfig(s.Config)
	result, err := instance.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: tableName,
	})
	if err != nil {
		logger.Logger.Error("Failed to describe table...", err)
		respch <- nil
		return
	}
	respch <- result
	return
}
