package controllers

import (
	"cloudview/app/src/api/authentication"
	"cloudview/app/src/api/encryption"
	custom_errors "cloudview/app/src/api/errors"
	"cloudview/app/src/api/middleware"
	"cloudview/app/src/api/middleware/logger"
	"cloudview/app/src/cache"
	"cloudview/app/src/database"
	"cloudview/app/src/helpers"
	models "cloudview/app/src/models/provider_accounts"
	"cloudview/app/src/providers/service"
	"cloudview/app/src/providers/service/aws"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

/*
*

	Fetch data from AWS,GCP,Azure services to
	show metrics and other data on `client-side`
*/
func (c *ServicesController) GetServiceData(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := middleware.RegisterResponses(w)
		/*
			Need to verify if user can access project.
		*/
		authenticatedUser, err := rw.User(db, r)
		if err != nil {
			logger.Logger.Error("ServicesController.GetServiceData: ERROR", err)
			if errors.Is(err, custom_errors.NoDataFound) {
				rw.Unauthorized()
				return
			}
			rw.Forbidden()
			return
		}
		projectId := r.URL.Query().Get("projectId")
		providerAccountId := r.URL.Query().Get("providerAccountId")
		region := r.URL.Query().Get("region")
		if projectId == "" || providerAccountId == "" {
			logger.Logger.Error("ServicesController.GetServiceData: ERROR", err)
			rw.Error("Invalid `projectId` or `seriviceId` provided", http.StatusBadRequest)
			return
		}
		if region == "" {
			rw.Error("Invalid region provided.", http.StatusBadRequest)
			return
		}
		isValidUUID := helpers.IsValidUUID(providerAccountId)
		if !isValidUUID {
			logger.Logger.Error("ServicesController.GetServiceData: Invalid service ID", err)
			rw.Error("Invalid service ID provided", http.StatusUnprocessableEntity)
			return
		}
		providerAccountUid, err := uuid.Parse(providerAccountId)
		if err != nil {
			rw.Error("Invalid service ID provided", http.StatusUnprocessableEntity)
			return
		}

		project, err := authentication.VerifyProjectAccess(db, projectId, authenticatedUser.ID)
		if err != nil {
			logger.Logger.Error("ServicesController.GetServiceData: ERROR, project verification failed", err)
			rw.Error(err.Error(), http.StatusForbidden)
			return
		}
		logger.Logger.Log("Project access verified", project.ID)
		providerAccount, err := models.GetByIdForSDK(db, providerAccountUid)
		if err != nil {
			logger.Logger.Error("ServicesController.GetServiceData: ERROR", err)
			rw.Error("Unable to fetch service data", http.StatusInternalServerError)
			return
		}

		cacheKey := fmt.Sprintf("service:%s:%s:%s", providerAccountId, providerAccount.Provider, region)
		// Add switch case to switch services like 'aws', 'gcp'
		// Caching data for 15 minutes
		result, err := cache.Fetch(cacheKey, 0, func() (interface{}, error) {
			accessKeySecret, err := encryption.Decrypt(providerAccount.AccessKeySecret, providerAccount.RotationSecretKey)
			if err != nil {
				logger.Logger.Error("Invalid provider access-key-secret", err)
				return nil, errors.New("Invalid provider secret")
			}
			return service.GetData(&aws.AWS{
				Region:            region,
				ProviderAccountID: providerAccountUid,
			}, providerAccount.AccessKeyID, accessKeySecret, region)
		})
		if err != nil {
			logger.Logger.Error("ServicesController.GetServiceData: ERROR fetching metrics", err)
			rw.Error("Unknown error occured", http.StatusInternalServerError)
			return
		}
		rw.Success(result, http.StatusOK)
		return
	}
}

/*
AWS - fetch integrations for apigateway route.
*/
func (c *ServicesController) GetApiGatewayV2Integrations(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := middleware.RegisterResponses(w)
		authenticatedUser, err := rw.User(db, r)
		if err != nil {
			logger.Logger.Error("ServicesController.GetApiGatewayV2Integrations: ERROR", err)
			if errors.Is(err, custom_errors.NoDataFound) {
				rw.Unauthorized()
				return
			}
			rw.Forbidden()
			return
		}
		projectId := r.URL.Query().Get("projectId")
		providerAccountId := r.URL.Query().Get("providerAccountId")
		region := r.URL.Query().Get("region")
		apiId := r.URL.Query().Get("apiId")
		if projectId == "" || providerAccountId == "" {
			logger.Logger.Error("ServicesController.GetApiGatewayV2Integrations: ERROR", err)
			rw.Error("Invalid `projectId` or `seriviceId` provided", http.StatusBadRequest)
			return
		}
		if region == "" {
			rw.Error("Invalid region provided.", http.StatusBadRequest)
			return
		}
		if apiId == "" {
			rw.Error("Invalid appId provided", http.StatusBadRequest)
			return
		}
		isValidUUID := helpers.IsValidUUID(providerAccountId)
		if !isValidUUID {
			logger.Logger.Error("ServicesController.GetApiGatewayV2Integrations: Invalid service ID", err)
			rw.Error("Invalid service ID provided", http.StatusUnprocessableEntity)
			return
		}
		providerAccountUid, err := uuid.Parse(providerAccountId)
		if err != nil {
			rw.Error("Invalid service ID provided", http.StatusUnprocessableEntity)
			return
		}

		project, err := authentication.VerifyProjectAccess(db, projectId, authenticatedUser.ID)
		if err != nil {
			logger.Logger.Error("ServicesController.GetApiGatewayV2Integrations: ERROR, project verification failed", err)
			rw.Error(err.Error(), http.StatusForbidden)
			return
		}
		logger.Logger.Log("Project access verified", project.ID)
		providerAccount, err := models.GetByIdForSDK(db, providerAccountUid)
		if err != nil {
			logger.Logger.Error("ServicesController.GetApiGatewayV2Integrations: ERROR", err)
			rw.Error("Unable to fetch service data", http.StatusInternalServerError)
			return
		}
		// cache for 15mins
		cacheKey := fmt.Sprintf("integrations:%s", apiId)
		result, err := cache.Fetch(cacheKey, 0, func() (interface{}, error) {
			awsClient := &aws.AWS{
				ProviderAccountID: providerAccountUid,
				Region:            region,
			}
			if err := awsClient.Init(providerAccount.AccessKeyID, providerAccount.AccessKeySecret, region); err != nil {
				logger.Logger.Error("ServicesController.GetApiGatewayV2Integrations: ERROR unable to initialize aws client", err)
				return nil, custom_errors.UnknownError
			}
			return awsClient.GetApiGatewayV2Integrations(apiId)
		})
		if err != nil {
			rw.Error(err.Error(), http.StatusInternalServerError)
			return
		}
		rw.Success(result, http.StatusOK)
		return
	}
}

func (c *ServicesController) GetUsage(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := middleware.RegisterResponses(w)
		authenticatedUser, err := rw.User(db, r)
		if err != nil {
			logger.Logger.Error("ServicesController.GetUsage: ERROR", err)
			if errors.Is(err, custom_errors.NoDataFound) {
				rw.Unauthorized()
				return
			}
			rw.Forbidden()
			return
		}
		params := mux.Vars(r)
		provider := strings.ToLower(params["provider"])
		if provider == "" {
			rw.Error("Invalid route provided. Malformed arguments", http.StatusBadRequest)
			return
		}
		projectId := r.URL.Query().Get("projectId")
		providerAccountId := r.URL.Query().Get("providerAccountId")
		instance := r.URL.Query().Get("instance")
		instanceId := r.URL.Query().Get("instanceId")
		region := r.URL.Query().Get("region")
		if projectId == "" || providerAccountId == "" || instance == "" || instanceId == "" {
			logger.Logger.Error("ServicesController.GetUsage: ERROR", err)
			rw.Error("Invalid `projectId`or `seriviceId` or `instance` or `instanceId` provided", http.StatusBadRequest)
			return
		}
		if region == "" {
			rw.Error("Invalid region provided.", http.StatusBadRequest)
			return
		}
		isValidUUID := helpers.IsValidUUID(providerAccountId)
		if !isValidUUID {
			logger.Logger.Error("ServicesController.GetUsage: Invalid service ID", err)
			rw.Error("Invalid service ID provided", http.StatusUnprocessableEntity)
			return
		}
		providerAccountUid, err := uuid.Parse(providerAccountId)
		if err != nil {
			rw.Error("Invalid service ID provided", http.StatusUnprocessableEntity)
			return
		}

		project, err := authentication.VerifyProjectAccess(db, projectId, authenticatedUser.ID)
		if err != nil {
			logger.Logger.Error("ServicesController.GetUsage: ERROR, project verification failed", err)
			rw.Error(err.Error(), http.StatusForbidden)
			return
		}
		logger.Logger.Log("Project access verified", project.ID)
		providerAccount, err := models.GetByIdForSDK(db, providerAccountUid)
		if err != nil {
			logger.Logger.Error("ServicesController.GetUsage: ERROR", err)
			rw.Error("Unable to fetch service data", http.StatusInternalServerError)
			return
		}
		switch provider {
		case "aws":
			// Caching data for 15 minutes
			cacheKey := fmt.Sprintf("%s:usage:%s:%s", provider, instance, instanceId)
			result, err := cache.Fetch(cacheKey, 0, func() (interface{}, error) {
				return GetAwsUsageData(db)(r, providerAccount, region, strings.ToLower(instance), instanceId)
			})
			if err != nil {
				logger.Logger.Error("ServicesController.GetUsage: ERROR", err)
				rw.Error(err.Error(), http.StatusBadRequest)
			}
			rw.Success(result, http.StatusOK)
			return
		default:
			rw.Error("404 Route not found", http.StatusNotFound)
			return
		}
	}
}
