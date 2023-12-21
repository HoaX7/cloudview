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
	projects_model "cloudview/app/src/models/projects"
	models "cloudview/app/src/models/services"
	"cloudview/app/src/providers/service"
	"cloudview/app/src/providers/service/aws"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

/**
This controller handles all process related to AWS, GCP ...etc services
that returns data from sdks.

Responsible for storing access_keys with encryption and rotating keys.
*/

func (c *ServiceController) StoreAccessKey(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := middleware.RegisterResponses(w)
		authenticatedUser, err := rw.User(db, r)
		if err != nil {
			logger.Logger.Error("ServiceController.StoreAccessKey: ERROR", err)
			if errors.Is(err, custom_errors.NoDataFound) {
				rw.Unauthorized()
				return
			}
			rw.Forbidden()
			return
		}
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			logger.Logger.Error("ServiceController.StoreAccessKey: Error reading request body:", err)
			rw.Error("Bad request", http.StatusBadRequest)
			return
		}
		var request models.Services
		if err := json.Unmarshal(body, &request); err != nil {
			logger.Logger.Error("ServiceController.StoreAccessKey: Error parsing request body:", err)
			rw.Error(custom_errors.UnknownError.Error(), http.StatusUnprocessableEntity)
			return
		}
		if request.Name == "" || request.AccessKeySecret == "" || request.AccessKeyID == "" || request.Provider == "" {
			rw.Error("Missing fields in body. Fields 'name', 'accessKeyId', 'accessKeySecret', 'provider' are required.", http.StatusBadRequest)
			return
		}
		isValidUUID := helpers.IsValidUUID(request.ProjectID.String())
		if !isValidUUID {
			logger.Logger.Error("ServiceController.StoreAccessKey: Invalid project ID provided", err)
			rw.Error("Invalid `projectId` of value uuid provided", http.StatusUnprocessableEntity)
			return
		}
		projectData, err := projects_model.GetById(db, request.ProjectID)
		if err != nil {
			logger.Logger.Error("ServiceController.StoreAccessKey: ERROR projectId", request.ProjectID, err)
			rw.Error("Please select a valid Project.", http.StatusNotFound)
			return
		}
		/**
		Only the project owner can create new services
		with access keys.
		Use the projectId passed in the body to verify
		if the authenticated user is the owner of the project
		before allowing them to create service.
		*/
		if projectData.OwnerID != authenticatedUser.ID {
			logger.Logger.Error("ServiceController.StoreAccessKey: ERROR project owner mismatch, Owner:", projectData.OwnerID, "AuthUser:", authenticatedUser.ID)
			rw.Error("Please contact your project owner to add Access Keys.", http.StatusForbidden)
			return
		}

		/*
			Generate 16 byte random key to encrypt `accessKeySecret`
		*/
		key, err := encryption.GenerateRandomSecretKey(16)
		if err != nil {
			logger.Logger.Error("ServiceController.StoreAccessKey: ERROR unable to generate secret key", err)
			rw.Error("Something went wrong, Please try again later", http.StatusInternalServerError)
			return
		}
		cipherText, err := encryption.Encrypt(request.AccessKeySecret, key)
		if err != nil {
			logger.Logger.Error("ServiceController.StoreAccessKey: ERROR unable to encrypt access key", err)
			rw.Error("Something went wrong, Please try again later", http.StatusInternalServerError)
			return
		}
		request.AccessKeySecret = cipherText
		request.RotationSecretKey = key
		result, err := models.Create(db, request)
		if err != nil {
			logger.Logger.Error("ServiceController.StoreAccessKey: ERROR unable to create data", err)
			rw.Error("Something went wrong, Please try again later", http.StatusInternalServerError)
			return
		}
		rw.Success(result, http.StatusOK)
		return
	}
}

/*
*

	Could considering add more restrictions to access services
	when more collaborators are added.

	For example: Permissions to view only 1 type of service (AWS, GCP)
*/
func (c *ServiceController) GetByProject(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := middleware.RegisterResponses(w)
		authenticatedUser, err := rw.User(db, r)
		if err != nil {
			logger.Logger.Error("ServiceController.GetServiceData: ERROR", err)
			if errors.Is(err, custom_errors.NoDataFound) {
				rw.Unauthorized()
				return
			}
			rw.Forbidden()
			return
		}
		projectId := r.URL.Query().Get("projectId")
		if projectId == "" {
			rw.Error("Invalid project ID", http.StatusBadRequest)
			return
		}
		/**
		verification to make sure the project
		can be accessed by user.
		*/
		project, err := authentication.VerifyProjectAccess(db, projectId, authenticatedUser.ID)
		if err != nil {
			logger.Logger.Error("ServiceController.GetByProject: invalid project uuid provided", err)
			rw.Error(err.Error(), http.StatusBadRequest)
			return
		}
		logger.Logger.Log("ServiceController.GetByProject: fetching data for projectId:", project.ID)
		result, err := models.GetByProjectId(db, project.ID)
		if err != nil {
			logger.Logger.Error("ServiceController.GetByProject: ERROR", err)
			rw.Error("Unable to fetch services", http.StatusInternalServerError)
			return
		}
		rw.Success(result, http.StatusOK)
		return
	}
}

/*
*

	Fetch data from AWS,GCP,Azure services to
	show metrics and other data on `client-side`
*/
func (c *ServiceController) GetServiceData(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := middleware.RegisterResponses(w)
		/*
			Need to verify if user can access project.
		*/
		authenticatedUser, err := rw.User(db, r)
		if err != nil {
			logger.Logger.Error("ServiceController.GetServiceData: ERROR", err)
			if errors.Is(err, custom_errors.NoDataFound) {
				rw.Unauthorized()
				return
			}
			rw.Forbidden()
			return
		}
		projectId := r.URL.Query().Get("projectId")
		serviceId := r.URL.Query().Get("serviceId")
		region := r.URL.Query().Get("region")
		if projectId == "" || serviceId == "" {
			logger.Logger.Error("ServiceController.GetServiceData: ERROR", err)
			rw.Error("Invalid `projectId` or `seriviceId` provided", http.StatusBadRequest)
			return
		}
		if region == "" {
			rw.Error("Invalid region provided.", http.StatusBadRequest)
			return
		}
		isValidUUID := helpers.IsValidUUID(serviceId)
		if !isValidUUID {
			logger.Logger.Error("ServiceController.GetServiceData: Invalid service ID", err)
			rw.Error("Invalid service ID provided", http.StatusUnprocessableEntity)
			return
		}
		serviceUid, err := uuid.Parse(serviceId)
		if err != nil {
			rw.Error("Invalid service ID provided", http.StatusUnprocessableEntity)
			return
		}

		project, err := authentication.VerifyProjectAccess(db, projectId, authenticatedUser.ID)
		if err != nil {
			logger.Logger.Error("ServiceController.GetServiceData: ERROR, project verification failed", err)
			rw.Error(err.Error(), http.StatusForbidden)
			return
		}
		logger.Logger.Log("Project access verified", project.ID)
		serviceData, err := models.GetById(db, serviceUid)
		if err != nil {
			logger.Logger.Error("ServiceController.GetServiceData: ERROR", err)
			rw.Error("Unable to fetch service data", http.StatusInternalServerError)
			return
		}

		cacheKey := fmt.Sprintf("service:%s:%s:%s", serviceId, serviceData.Provider, region)
		// Add switch case to switch services like 'aws', 'gcp'
		// Caching data for 15 minutes
		result, err := cache.Fetch(cacheKey, 0, func() (interface{}, error) {
			accessKeySecret, err := encryption.Decrypt(serviceData.AccessKeySecret, serviceData.RotationSecretKey)
			if err != nil {
				logger.Logger.Error("Invalid provider access-key-secret", err)
				return nil, errors.New("Invalid provider secret")
			}
			return service.GetData(&aws.AWS{
				Region:    region,
				ServiceId: serviceUid,
			}, serviceData.AccessKeyID, accessKeySecret, region)
		})
		if err != nil {
			logger.Logger.Error("ServiceController.GetServiceData: ERROR fetching metrics", err)
			rw.Error("Unknown error occured", http.StatusInternalServerError)
			return
		}
		rw.Success(result, http.StatusOK)
		return
	}
}

func (c *ServiceController) UpdateService(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := middleware.RegisterResponses(w)
		authenticatedUser, err := rw.User(db, r)
		if err != nil {
			logger.Logger.Error("ServiceController.UpdateService: ERROR", err)
			if errors.Is(err, custom_errors.NoDataFound) {
				rw.Unauthorized()
				return
			}
			rw.Forbidden()
			return
		}
		id := mux.Vars(r)["id"]
		isIDValidUUID := helpers.IsValidUUID(id)
		if !isIDValidUUID {
			logger.Logger.Error("ServiceController.UpdateService: Invalid service ID", err)
			rw.Error("Invalid service ID provided", http.StatusUnprocessableEntity)
			return
		}
		serviceId, err := uuid.Parse(id)
		if err != nil {
			rw.Error("Invalid service ID provided", http.StatusUnprocessableEntity)
			return
		}
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			logger.Logger.Error("ServiceController.UpdateService: Error reading request body:", err)
			rw.Error("Bad request", http.StatusBadRequest)
			return
		}
		/**
		NOTE: `omitempty` has no effect on `json.Unmarshal`, so the
		`request` body will have default values for the fields with no values.

		To update data you must validate each field manually to append to a column list.
		*/
		var request models.Services
		if err := json.Unmarshal(body, &request); err != nil {
			logger.Logger.Error("ServiceController.UpdateService: Error parsing request body:", err)
			rw.Error(custom_errors.UnknownError.Error(), http.StatusUnprocessableEntity)
			return
		}
		project, err := authentication.VerifyProjectAccess(db, request.ProjectID, authenticatedUser.ID)
		if err != nil {
			logger.Logger.Error("ServiceController.UpdateService: Project verification failed", err)
			rw.Error(err.Error(), http.StatusForbidden)
			return
		}
		/**
		Only the project owner can create new services
		with access keys.
		Use the projectId passed in the body to verify
		if the authenticated user is the owner of the project
		before allowing them to create service.
		*/
		if project.OwnerID != authenticatedUser.ID {
			logger.Logger.Error("ServiceController.UpdateService: ERROR project owner mismatch, Owner:", project.OwnerID, "AuthUser:", authenticatedUser.ID)
			rw.Error("Please contact your project owner to edit service details.", http.StatusForbidden)
			return
		}
		if err := models.Update(db, serviceId, request); err != nil {
			rw.Error("Unable to save data", http.StatusInternalServerError)
			return
		}
		rw.Success("data saved", http.StatusOK)
		return
	}
}

/*
AWS - fetch integrations for apigateway route.
*/
func (c *ServiceController) GetApiGatewayV2Integrations(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := middleware.RegisterResponses(w)
		authenticatedUser, err := rw.User(db, r)
		if err != nil {
			logger.Logger.Error("ServiceController.GetApiGatewayV2Integrations: ERROR", err)
			if errors.Is(err, custom_errors.NoDataFound) {
				rw.Unauthorized()
				return
			}
			rw.Forbidden()
			return
		}
		projectId := r.URL.Query().Get("projectId")
		serviceId := r.URL.Query().Get("serviceId")
		region := r.URL.Query().Get("region")
		apiId := r.URL.Query().Get("apiId")
		if projectId == "" || serviceId == "" {
			logger.Logger.Error("ServiceController.GetApiGatewayV2Integrations: ERROR", err)
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
		isValidUUID := helpers.IsValidUUID(serviceId)
		if !isValidUUID {
			logger.Logger.Error("ServiceController.GetApiGatewayV2Integrations: Invalid service ID", err)
			rw.Error("Invalid service ID provided", http.StatusUnprocessableEntity)
			return
		}
		serviceUid, err := uuid.Parse(serviceId)
		if err != nil {
			rw.Error("Invalid service ID provided", http.StatusUnprocessableEntity)
			return
		}

		project, err := authentication.VerifyProjectAccess(db, projectId, authenticatedUser.ID)
		if err != nil {
			logger.Logger.Error("ServiceController.GetApiGatewayV2Integrations: ERROR, project verification failed", err)
			rw.Error(err.Error(), http.StatusForbidden)
			return
		}
		logger.Logger.Log("Project access verified", project.ID)
		serviceData, err := models.GetById(db, serviceUid)
		if err != nil {
			logger.Logger.Error("ServiceController.GetApiGatewayV2Integrations: ERROR", err)
			rw.Error("Unable to fetch service data", http.StatusInternalServerError)
			return
		}
		// cache for 15mins
		cacheKey := fmt.Sprintf("integrations:%s", apiId)
		result, err := cache.Fetch(cacheKey, 0, func() (interface{}, error) {
			awsClient := &aws.AWS{
				ServiceId: serviceUid,
				Region:    region,
			}
			if err := awsClient.Init(serviceData.AccessKeyID, serviceData.AccessKeySecret, region); err != nil {
				logger.Logger.Error("ServiceController.GetApiGatewayV2Integrations: ERROR unable to initialize aws client", err)
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

func (c *ServiceController) GetUsage(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := middleware.RegisterResponses(w)
		authenticatedUser, err := rw.User(db, r)
		if err != nil {
			logger.Logger.Error("ServiceController.GetUsage: ERROR", err)
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
		serviceId := r.URL.Query().Get("serviceId")
		instance := r.URL.Query().Get("instance")
		instanceId := r.URL.Query().Get("instanceId")
		region := r.URL.Query().Get("region")
		if projectId == "" || serviceId == "" || instance == "" || instanceId == "" {
			logger.Logger.Error("ServiceController.GetUsage: ERROR", err)
			rw.Error("Invalid `projectId`or `seriviceId` or `instance` or `instanceId` provided", http.StatusBadRequest)
			return
		}
		if region == "" {
			rw.Error("Invalid region provided.", http.StatusBadRequest)
			return
		}
		isValidUUID := helpers.IsValidUUID(serviceId)
		if !isValidUUID {
			logger.Logger.Error("ServiceController.GetUsage: Invalid service ID", err)
			rw.Error("Invalid service ID provided", http.StatusUnprocessableEntity)
			return
		}
		serviceUid, err := uuid.Parse(serviceId)
		if err != nil {
			rw.Error("Invalid service ID provided", http.StatusUnprocessableEntity)
			return
		}

		project, err := authentication.VerifyProjectAccess(db, projectId, authenticatedUser.ID)
		if err != nil {
			logger.Logger.Error("ServiceController.GetUsage: ERROR, project verification failed", err)
			rw.Error(err.Error(), http.StatusForbidden)
			return
		}
		logger.Logger.Log("Project access verified", project.ID)
		serviceData, err := models.GetById(db, serviceUid)
		if err != nil {
			logger.Logger.Error("ServiceController.GetUsage: ERROR", err)
			rw.Error("Unable to fetch service data", http.StatusInternalServerError)
			return
		}
		switch provider {
		case "aws":
			// Caching data for 15 minutes
			cacheKey := fmt.Sprintf("%s:usage:%s:%s", provider, instance, instanceId)
			result, err := cache.Fetch(cacheKey, 0, func() (interface{}, error) {
				return GetAwsUsageData(db)(r, serviceData, region, strings.ToLower(instance), instanceId)
			})
			if err != nil {
				logger.Logger.Error("ServiceController.GetUsage: ERROR", err)
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
