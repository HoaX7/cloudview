package controllers

import (
	"cloudview/app/src/api/encryption"
	"cloudview/app/src/api/middleware/logger"
	"cloudview/app/src/database"
	models "cloudview/app/src/models"
	"cloudview/app/src/providers/service/aws"
	"cloudview/app/src/types"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
)

func GetAwsUsageData(db *database.DB) func(*http.Request, models.ProviderAccounts, string, string, string) (interface{}, error) {
	return func(r *http.Request, providerAccount models.ProviderAccounts, region string, instance string, instanceId string) (interface{}, error) {
		client := aws.AWS{
			Region:            region,
			ProviderAccountID: providerAccount.ID,
		}
		accessKeySecret, err := encryption.Decrypt(providerAccount.AccessKeySecret, providerAccount.RotationSecretKey)
		if err != nil {
			logger.Logger.Error("Invalid provider access-key-secret", err)
			return nil, errors.New("Invalid provider secret")
		}
		client.Init(providerAccount.AccessKeyID, accessKeySecret, region)
		caller := client.GetServiceCaller()
		caller.CloudWatchInit()
		switch instance {
		case "ec2":
			var bandwidthResult *cloudwatch.GetMetricDataOutput
			var cpuResult *cloudwatch.GetMetricDataOutput
			wg := &sync.WaitGroup{}
			respch := make(chan *types.Usage, 2)

			// fetch ec2 bandwidth and cpu usage metrics
			wg.Add(1)
			go caller.GetEC2Bandwidth(respch, wg, instanceId)

			wg.Add(1)
			go caller.GetEC2CpuUsage(respch, wg, instanceId)
			wg.Wait()
			close(respch)

			for {
				select {
				case val, ok := <-respch:
					if !ok {
						return &types.Usage{
							Bandwidth: bandwidthResult,
							Cpu:       cpuResult,
						}, nil
					}
					if val.Bandwidth != nil {
						bandwidthResult = val.Bandwidth
					}
					if val.Cpu != nil {
						cpuResult = val.Cpu
					}
				}
			}
		default:
			return "", errors.New(fmt.Errorf("Service %s does not exist", instance).Error())
		}
	}
}
