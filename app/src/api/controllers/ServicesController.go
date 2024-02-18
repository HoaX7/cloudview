package controllers

import (
	"cloudview/app/src/api/authentication"
	"cloudview/app/src/api/encryption"
	custom_errors "cloudview/app/src/api/errors"
	"cloudview/app/src/api/middleware"
	"cloudview/app/src/cache"
	"cloudview/app/src/database"
	models "cloudview/app/src/models/provider_accounts"
	"cloudview/app/src/providers/service"
	"cloudview/app/src/providers/service/aws"
	"cloudview/app/src/types"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

/*
*

	Fetch data from AWS,GCP,Azure services to
	show metrics and other data on `client-side`
*/
func (c *ServicesController) GetServiceData(db *database.DB) http.HandlerFunc {
	c.Logger.SetName(c.Name + ".GetServiceData")
	return func(w http.ResponseWriter, r *http.Request) {
		rw := middleware.CustomResponseWriter(w)
		/*
			Need to verify if user can access project.
		*/
		authenticatedUser, err := rw.User(db, r)
		if err != nil {
			c.Logger.Error("ERROR", err)
			if errors.Is(err, custom_errors.NoDataFound) {
				rw.Unauthorized()
				return
			}
			rw.Forbidden()
			return
		}
		providerAccountId := r.URL.Query().Get("providerAccountId")
		region := r.URL.Query().Get("region")
		if providerAccountId == "" {
			c.Logger.Error("ERROR", err)
			rw.Error("Invalid `providerAccountId` provided", http.StatusBadRequest)
			return
		}
		if region == "" {
			rw.Error("Invalid region provided.", http.StatusBadRequest)
			return
		}

		verifiedData, err := authentication.VerifyProjectAccess(db, authenticatedUser.ID, types.VerifyProjectAccessInput{
			ProviderAccountID: providerAccountId,
		})
		if err != nil {
			c.Logger.Error("ERROR, project verification failed", err)
			rw.Error(err.Error(), http.StatusForbidden)
			return
		}
		providerAccount, err := models.GetByIdForSDK(db, verifiedData.ProviderAccount.ID)
		if err != nil {
			c.Logger.Error("ERROR", err)
			rw.Error("Unable to fetch service data", http.StatusInternalServerError)
			return
		}

		cacheKey := fmt.Sprintf("service:%s:%s:%s", providerAccount.ID, providerAccount.Provider, region)
		// Add switch case to switch services like 'aws', 'gcp'
		// Caching data for 15 minutes
		var result interface{}
		if cache.Fetch(cacheKey, 0, &result, func() (interface{}, error) {
			accessKeySecret, err := encryption.Decrypt(providerAccount.AccessKeySecret, providerAccount.RotationSecretKey)
			if err != nil {
				c.Logger.Error("Invalid provider access-key-secret", err)
				return nil, errors.New("Invalid provider secret")
			}
			return service.GetData(&aws.AWS{
				Region:            region,
				ProviderAccountID: providerAccount.ID,
			}, providerAccount.AccessKeyID, accessKeySecret, region)
		}); err != nil {
			c.Logger.Error("ERROR fetching metrics", err)
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
	c.Logger.SetName(c.Name + ".GetApiGatewayV2Integrations")
	return func(w http.ResponseWriter, r *http.Request) {
		rw := middleware.CustomResponseWriter(w)
		authenticatedUser := rw.SessionUser
		providerAccountId := r.URL.Query().Get("providerAccountId")
		region := r.URL.Query().Get("region")
		apiId := r.URL.Query().Get("apiId")
		if providerAccountId == "" {
			rw.Error("Invalid `providerAccountId` provided", http.StatusBadRequest)
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

		verifiedData, err := authentication.VerifyProjectAccess(db, authenticatedUser.ID, types.VerifyProjectAccessInput{
			ProviderAccountID: providerAccountId,
		})
		if err != nil {
			c.Logger.Error("ERROR, project verification failed", err)
			rw.Error(err.Error(), http.StatusForbidden)
			return
		}
		c.Logger.Log("Project access verified", verifiedData.ProjectAccessDetails.Projects.ID)
		providerAccount, err := models.GetByIdForSDK(db, verifiedData.ProviderAccount.ID)
		if err != nil {
			c.Logger.Error("ERROR", err)
			rw.Error("Unable to fetch service data", http.StatusInternalServerError)
			return
		}
		// cache for 15mins
		cacheKey := fmt.Sprintf("integrations:%s", apiId)
		var result interface{}
		if cache.Fetch(cacheKey, 0, &result, func() (interface{}, error) {
			awsClient := &aws.AWS{
				ProviderAccountID: providerAccount.ID,
				Region:            region,
			}
			accessKeySecret, err := encryption.Decrypt(providerAccount.AccessKeySecret, providerAccount.RotationSecretKey)
			if err != nil {
				c.Logger.Error("Invalid provider access-key-secret", err)
				return nil, errors.New("invalid provider secret")
			}
			if err := awsClient.Init(providerAccount.AccessKeyID, accessKeySecret, region); err != nil {
				c.Logger.Error("ERROR unable to initialize aws client", err)
				return nil, custom_errors.UnknownError
			}
			return awsClient.GetApiGatewayV2Integrations(apiId)
		}); err != nil {
			rw.Error(err.Error(), http.StatusInternalServerError)
			return
		}
		rw.Success(result, http.StatusOK)
		return
	}
}

func (c *ServicesController) GetUsage(db *database.DB) http.HandlerFunc {
	c.Logger.SetName(c.Name + ".GetUsage")
	return func(w http.ResponseWriter, r *http.Request) {
		rw := middleware.CustomResponseWriter(w)
		authenticatedUser := rw.SessionUser
		params := mux.Vars(r)
		provider := strings.ToLower(params["provider"])
		if provider == "" {
			rw.Error("Invalid route provided. Malformed arguments", http.StatusBadRequest)
			return
		}
		providerAccountId := r.URL.Query().Get("providerAccountId")
		instance := r.URL.Query().Get("instance")
		instanceId := r.URL.Query().Get("instanceId")
		region := r.URL.Query().Get("region")
		if providerAccountId == "" || instance == "" || instanceId == "" {
			rw.Error("Invalid `providerAccountId` or `instance` or `instanceId` provided", http.StatusBadRequest)
			return
		}
		if region == "" {
			rw.Error("Invalid region provided.", http.StatusBadRequest)
			return
		}

		verifiedData, err := authentication.VerifyProjectAccess(db, authenticatedUser.ID, types.VerifyProjectAccessInput{
			ProviderAccountID: providerAccountId,
		})
		if err != nil {
			c.Logger.Error("ERROR, project verification failed", err)
			rw.Error(err.Error(), http.StatusForbidden)
			return
		}
		providerAccount, err := models.GetByIdForSDK(db, verifiedData.ProviderAccount.ID)
		if err != nil {
			c.Logger.Error("ERROR", err)
			rw.Error("Unable to fetch service data", http.StatusInternalServerError)
			return
		}
		switch provider {
		case "aws":
			// Caching data for 15 minutes
			var result interface{}
			cacheKey := fmt.Sprintf("%s:usage:%s:%s", provider, instance, instanceId)
			if cache.Fetch(cacheKey, 0, &result, func() (interface{}, error) {
				return GetAwsUsageData(db)(r, providerAccount, region, strings.ToLower(instance), instanceId)
			}); err != nil {
				c.Logger.Error("ERROR", err)
				rw.Error(err.Error(), http.StatusBadRequest)
				return
			}
			rw.Success(result, http.StatusOK)
			return
		default:
			rw.Error("404 Route not found", http.StatusNotFound)
			return
		}
	}
}
