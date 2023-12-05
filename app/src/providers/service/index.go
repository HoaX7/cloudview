package service

import (
	"cloudview/app/src/api/middleware/logger"
	"cloudview/app/src/types"
	"errors"
)

/*
*
Responsible to fetch metric data from AWS,GCP,AZURE ...etc services
*/
func GetData(service types.DataService, accessKeyId string, accessKeySecret string, region string) (any, error) {
	if err := service.Init(accessKeyId, accessKeySecret, region); err != nil {
		logger.Logger.Error("providers.services.index.GetData: Init failed", err)
		return nil, errors.New("Unable to initialize service. Are the credentials correct?")
	}
	logger.Logger.Log("Fetch metric data from service:", service.Name())
	return service.GetData()
}
