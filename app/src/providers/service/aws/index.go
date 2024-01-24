package aws

import (
	custom_errors "cloudview/app/src/api/errors"
	"cloudview/app/src/api/middleware/logger"
	"cloudview/app/src/providers/service/aws/services"
	"cloudview/app/src/types"
	"context"
	"sync"
	"time"

	/**
	AWS services must be installed separetly
	using `go get <service name>`

	read: https://pkg.go.dev/github.com/aws/aws-sdk-go-v2#section-readme
	to see all available services
	*/
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/google/uuid"
)

type AWS struct {
	Config            aws.Config
	Region            string
	ProviderAccountID uuid.UUID
}

func (s *AWS) Name() string {
	return "AWS"
}

func (s *AWS) Init(accessKeyId string, accessKeySecret string, region string) error {
	creds := aws.Credentials{
		AccessKeyID:     accessKeyId,
		SecretAccessKey: accessKeySecret,
		// SessionToken:    "init-session", // Not sure what to send. But works even if this field isnt passed
	}
	start := time.Now()
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithCredentialsProvider(
		credentials.StaticCredentialsProvider{
			Value: creds,
		},
	))
	if err != nil {
		logger.Logger.Error("aws.Init(): Unable to initialize sdk", err)
		return err
	}

	cfg.Region = region
	logger.Logger.Log(cfg.Region, accessKeyId, accessKeySecret)
	/*
		Required to make sure the credentials provided are correct.
	*/
	stsClient := sts.NewFromConfig(cfg)
	res, err := stsClient.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
	if err != nil {
		logger.Logger.Error("aws.Init(): Error validating credentials--", err)
		return custom_errors.InvalidCredentials
	}
	logger.Logger.Log("aws.Init(): initialized successfully for account:", *res.Account, "in:", time.Since(start))

	// Initialize config to AWS struct
	// to be used in other services to fetch data
	s.Config = cfg

	return nil
}

func (s *AWS) GetData() ([]types.GetDataOutput, error) {
	caller := services.Services{
		Config: s.Config,
		Region: s.Region,
	}
	/*
		Goroutines are used to run functions
		concurrently.
	*/
	start := time.Now()
	tasksToComplete := 11

	// Buffered channels, Make sure the stack size is equal
	// to number of functions being called.
	//
	// Note: Always use buffered channels to prevent deadlocks

	respch := make(chan types.GetDataOutput, tasksToComplete)
	wg := &sync.WaitGroup{}

	wg.Add(tasksToComplete)
	go caller.GetEC2Instances(respch, wg)
	go caller.GetApiGatewayV2Data(respch, wg)
	go caller.GetLambdaFuntions(respch, wg)
	go caller.GetS3Buckets(respch, wg)
	go caller.GetRoute53Data(respch, wg)
	go caller.GetRDSInstances(respch, wg)
	go caller.GetEFSData(respch, wg)
	go caller.GetEKSData(respch, wg)
	go caller.GetELBData(respch, wg)
	go caller.GetCloudfrontDistributions(respch, wg)
	go caller.GetDynamodbTables(respch, wg)

	wg.Wait()
	close(respch)

	logger.Logger.Log("aws.GetData: finished in ", time.Since(start))
	// Grab the responses from channel
	// in the order of execution.
	// ex: 1st data will always be ec2. 2nd is apigateway etc...
	var result []types.GetDataOutput
	for {
		select {
		case val, ok := <-respch:
			// Since the channel is already closed
			// return the data here...
			if !ok {
				return result, nil
			}
			if val.Result != nil {
				result = append(result, val)
			}
			// switch v := val.(type) {
			// case services.Ec2Response:
			// 	result = append(result, types.GetDataOutput{
			// 		Name:   "ec2",
			// 		Result: v.Ec2,
			// 	})
			// case services.ApiGatewayV2Response:
			// 	result = append(result, types.GetDataOutput{
			// 		Name:   "apigatewayV2",
			// 		Result: v.AGV2,
			// 	})
			// }
		}
	}
}

func (s *AWS) GetApiGatewayV2Integrations(apiId string) (*apigatewayv2.GetIntegrationsOutput, error) {
	caller := services.Services{
		Config: s.Config,
		Region: s.Region,
	}
	return caller.GetApiGatewayV2Integrations(apiId)
}

func (s *AWS) GetServiceCaller() *services.Services {
	caller := services.Services{
		Config: s.Config,
		Region: s.Region,
	}
	return &caller
}
