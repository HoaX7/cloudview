package services

import (
	"cloudview/app/src/api/middleware/logger"
	"cloudview/app/src/helpers/constants"
	"cloudview/app/src/types"
	"context"
	"errors"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2"
)

/*
*
This method only returns the number of apigateways created
with a `AppId` prop which then can be used to fetch
the routes integrated.

To map to ec2/lambda machine use the `IntegrationUri` to match
public ip or lambda ip
*/
func (s Services) GetApiGatewayV2Data(respch chan<- types.GetDataOutput, wg *sync.WaitGroup) {
	defer wg.Done()
	instance := apigatewayv2.NewFromConfig(s.Config)
	result, err := instance.GetApis(context.TODO(), &apigatewayv2.GetApisInput{})
	if err != nil {
		logger.Logger.Error("apigatwwayV2.GetApiGatewayV2Data: ERROR", err)
		respch <- types.GetDataOutput{
			Name:   constants.APIGATEWAYV2,
			Result: nil,
		}
		return
	}
	logger.Logger.Log("apigatewayV2.GetApiGatewayV2Data: success")
	if len(result.Items) <= 0 {
		respch <- types.GetDataOutput{
			Name:   constants.APIGATEWAYV2,
			Result: nil,
		}
		return
	}
	respch <- types.GetDataOutput{
		Name:   constants.APIGATEWAYV2,
		Result: result,
	}
	return
}

func (s Services) GetApiGatewayV2Integrations(apiId string) (*apigatewayv2.GetIntegrationsOutput, error) {
	instance := apigatewayv2.NewFromConfig(s.Config)
	result, err := instance.GetIntegrations(context.TODO(), &apigatewayv2.GetIntegrationsInput{
		ApiId: &apiId,
	})
	if err != nil {
		logger.Logger.Error("apigatewayV2.GetIntagrations: ERROR", err)
		return nil, errors.New("Unable to fetch integrations")
	}

	logger.Logger.Log("apigatewayV2.GetIntegrations: success")
	return result, nil
}
