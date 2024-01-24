package controllers

import (
	"cloudview/app/src/api/authentication"
	"cloudview/app/src/api/encryption"
	custom_errors "cloudview/app/src/api/errors"
	"cloudview/app/src/api/middleware"
	"cloudview/app/src/api/middleware/logger"
	"cloudview/app/src/database"
	"cloudview/app/src/helpers"
	"cloudview/app/src/models"
	projects_model "cloudview/app/src/models/projects"
	provider_models "cloudview/app/src/models/provider_accounts"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

/**
This controller handles all process related to AWS, GCP ...etc provider accounts
that returns data from sdks.

Responsible for storing access_keys with encryption and rotating keys.
*/

// @deprecated - in favor using 'cross-account-access' to authenticate aws-sdk
func (c *ProviderAccountsController) StoreAccessKey(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := middleware.RegisterResponses(w)
		authenticatedUser, err := rw.User(db, r)
		if err != nil {
			logger.Logger.Error("ProviderAccountsController.StoreAccessKey: ERROR", err)
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
			logger.Logger.Error("ProviderAccountsController.StoreAccessKey: Error reading request body:", err)
			rw.Error("Bad request", http.StatusBadRequest)
			return
		}
		var request models.ProviderAccounts
		if err := json.Unmarshal(body, &request); err != nil {
			logger.Logger.Error("ProviderAccountsController.StoreAccessKey: Error parsing request body:", err)
			rw.Error(custom_errors.UnknownError.Error(), http.StatusUnprocessableEntity)
			return
		}
		if request.Name == "" || request.AccessKeySecret == "" || request.AccessKeyID == "" || request.Provider == "" {
			rw.Error("Missing fields in body. Fields 'name', 'accessKeyId', 'accessKeySecret', 'provider' are required.", http.StatusBadRequest)
			return
		}
		isValidUUID := helpers.IsValidUUID(request.ProjectID.String())
		if !isValidUUID {
			logger.Logger.Error("ProviderAccountsController.StoreAccessKey: Invalid project ID provided", err)
			rw.Error("Invalid `projectId` of value uuid provided", http.StatusUnprocessableEntity)
			return
		}
		projectData, err := projects_model.GetById(db, *request.ProjectID)
		if err != nil {
			logger.Logger.Error("ProviderAccountsController.StoreAccessKey: ERROR projectId", request.ProjectID, err)
			rw.Error("Please select a valid Project.", http.StatusNotFound)
			return
		}
		/**
		Only the project owner can create new provider accounts
		with access keys.
		Use the projectId passed in the body to verify
		if the authenticated user is the owner of the project
		before allowing them to create provider account.
		*/
		if *projectData.OwnerID != authenticatedUser.ID {
			logger.Logger.Error("ProviderAccountsController.StoreAccessKey: ERROR project owner mismatch, Owner:", projectData.OwnerID, "AuthUser:", authenticatedUser.ID)
			rw.Error("Please contact your project owner to add Access Keys.", http.StatusForbidden)
			return
		}

		/*
			Generate 16 byte random key to encrypt `accessKeySecret`
		*/
		key, err := encryption.GenerateRandomSecretKey(16)
		if err != nil {
			logger.Logger.Error("ProviderAccountsController.StoreAccessKey: ERROR unable to generate secret key", err)
			rw.Error("Something went wrong, Please try again later", http.StatusInternalServerError)
			return
		}
		cipherText, err := encryption.Encrypt(request.AccessKeySecret, key)
		if err != nil {
			logger.Logger.Error("ProviderAccountsController.StoreAccessKey: ERROR unable to encrypt access key", err)
			rw.Error("Something went wrong, Please try again later", http.StatusInternalServerError)
			return
		}
		request.AccessKeySecret = cipherText
		request.RotationSecretKey = key
		result, err := provider_models.Create(db, request)
		if err != nil {
			logger.Logger.Error("ProviderAccountsController.StoreAccessKey: ERROR unable to create data", err)
			rw.Error("Something went wrong, Please try again later", http.StatusInternalServerError)
			return
		}
		rw.Success(result, http.StatusOK)
		return
	}
}

func (c *ProviderAccountsController) GetById(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := middleware.RegisterResponses(w)
		authenticatedUser, err := rw.User(db, r)
		if err != nil {
			logger.Logger.Error("ProviderAccountsController.GetById: ERROR", err)
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
			logger.Logger.Error("ProviderAccountsController.GetById: Invalid provider account ID", err)
			rw.Error("Invalid provider account ID provided", http.StatusUnprocessableEntity)
			return
		}
		providerAccountId, err := uuid.Parse(id)
		if err != nil {
			rw.Error("Invalid provider account ID provided", http.StatusUnprocessableEntity)
			return
		}
		projectId := r.URL.Query().Get("projectId")
		isProjectIDValidUUID := helpers.IsValidUUID(projectId)
		if !isProjectIDValidUUID {
			logger.Logger.Error("ProviderAccountsController.GetById: Invalid project ID", err)
			rw.Error("Please select a valid project", http.StatusUnprocessableEntity)
			return
		}
		projectUUID, err := uuid.Parse(projectId)
		if err != nil {
			logger.Logger.Error("ProviderAccountsController.GetById: Unable to parse project UUID", err)
			rw.Error(custom_errors.UnknownError.Error(), http.StatusInternalServerError)
			return
		}
		_, verificationErr := authentication.VerifyProjectAccess(db, projectUUID, authenticatedUser.ID)
		if verificationErr != nil {
			logger.Logger.Error("ProviderAccountsController.GetById: Project verification failed", verificationErr)
			rw.Error(verificationErr.Error(), http.StatusForbidden)
			return
		}
		result, err := provider_models.GetById(db, providerAccountId)
		if err != nil {
			logger.Logger.Error("ProviderAccountsController.GetById: ERROR", err)
			rw.Error("Unable to fetch accounts details", http.StatusInternalServerError)
			return
		}
		rw.Success(result, http.StatusOK)
		return
	}
}

/*
*

	Could considering add more restrictions to access provider accounts
	when more collaborators are added.

	For example: Permissions to view only 1 type of provider account (AWS, GCP)
*/
func (c *ProviderAccountsController) GetByProject(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := middleware.RegisterResponses(w)
		authenticatedUser, err := rw.User(db, r)
		if err != nil {
			logger.Logger.Error("ProviderAccountsController.GetByProject: ERROR", err)
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
			logger.Logger.Error("ProviderAccountsController.GetByProject: invalid project uuid provided", err)
			rw.Error(err.Error(), http.StatusBadRequest)
			return
		}
		logger.Logger.Log("ProviderAccountsController.GetByProject: fetching data for projectId:", project.ID)
		result, err := provider_models.GetByProjectId(db, project.ID)
		if err != nil {
			logger.Logger.Error("ProviderAccountsController.GetByProject: ERROR", err)
			rw.Error("Unable to fetch provider accounts", http.StatusInternalServerError)
			return
		}
		rw.Success(result, http.StatusOK)
		return
	}
}

func (c *ProviderAccountsController) UpdateProviderAccount(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := middleware.RegisterResponses(w)
		authenticatedUser, err := rw.User(db, r)
		if err != nil {
			logger.Logger.Error("ProviderAccountsController.UpdateProviderAccount: ERROR", err)
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
			logger.Logger.Error("ProviderAccountsController.UpdateProviderAccount: Invalid provider account ID", err)
			rw.Error("Invalid provider account ID provided", http.StatusUnprocessableEntity)
			return
		}
		providerAccountId, err := uuid.Parse(id)
		if err != nil {
			rw.Error("Invalid provider account ID provided", http.StatusUnprocessableEntity)
			return
		}
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			logger.Logger.Error("ProviderAccountsController.UpdateProviderAccount: Error reading request body:", err)
			rw.Error("Bad request", http.StatusBadRequest)
			return
		}
		/**
		NOTE: `omitempty` has no effect on `json.Unmarshal`, so the
		`request` body will have default values for the fields with no values.

		To update data you must validate each field manually to append to a column list.
		*/
		var request models.ProviderAccounts
		if err := json.Unmarshal(body, &request); err != nil {
			logger.Logger.Error("ProviderAccountsController.UpdateProviderAccount: Error parsing request body:", err)
			rw.Error(custom_errors.UnknownError.Error(), http.StatusUnprocessableEntity)
			return
		}
		project, err := authentication.VerifyProjectAccess(db, request.ProjectID, authenticatedUser.ID)
		if err != nil {
			logger.Logger.Error("ProviderAccountsController.UpdateProviderAccount: Project verification failed", err)
			rw.Error(err.Error(), http.StatusForbidden)
			return
		}
		/**
		Only the project owner can create new provider accounts
		with access keys.
		Use the projectId passed in the body to verify
		if the authenticated user is the owner of the project
		before allowing them to create provider account.
		*/
		if *project.OwnerID != authenticatedUser.ID {
			logger.Logger.Error("ProviderAccountsController.UpdateProviderAccount: ERROR project owner mismatch, Owner:", project.OwnerID, "AuthUser:", authenticatedUser.ID)
			rw.Error("Please contact your project owner to edit provider account details.", http.StatusForbidden)
			return
		}
		if err := provider_models.Update(db, providerAccountId, request); err != nil {
			rw.Error("Unable to save data", http.StatusInternalServerError)
			return
		}
		rw.Success("data saved", http.StatusOK)
		return
	}
}
