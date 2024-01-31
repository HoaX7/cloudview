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
	"cloudview/app/src/types"
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var log_pa = logger.NewLogger()

/**
This controller handles all process related to AWS, GCP ...etc provider accounts
that returns data from sdks.

Responsible for storing access_keys with encryption and rotating keys.
*/

// @deprecated - in favor using 'cross-account-access' to authenticate aws-sdk
func (c *ProviderAccountsController) StoreAccessKey(db *database.DB) http.HandlerFunc {
	log_pa.SetName(c.Name() + ".StoreAccessKey")
	return func(w http.ResponseWriter, r *http.Request) {
		rw := middleware.CustomResponseWriter(w)
		authenticatedUser := rw.SessionUser
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			log_pa.Error("Error reading request body:", err)
			rw.Error("Bad request", http.StatusBadRequest)
			return
		}
		var request models.ProviderAccounts
		if err := json.Unmarshal(body, &request); err != nil {
			log_pa.Error("Error parsing request body:", err)
			rw.Error(custom_errors.UnknownError.Error(), http.StatusUnprocessableEntity)
			return
		}
		if request.Name == "" || request.AccessKeySecret == "" || request.AccessKeyID == "" || request.Provider == "" {
			rw.Error("Missing fields in body. Fields 'name', 'accessKeyId', 'accessKeySecret', 'provider' are required.", http.StatusBadRequest)
			return
		}
		isValidUUID := helpers.IsValidUUID(request.ProjectID.String())
		if !isValidUUID {
			log_pa.Error("Invalid project ID provided", err)
			rw.Error("Invalid `projectId` of value uuid provided", http.StatusUnprocessableEntity)
			return
		}
		projectData, err := projects_model.GetById(db, *request.ProjectID)
		if err != nil {
			log_pa.Error("ERROR projectId", request.ProjectID, err)
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
			log_pa.Error("ERROR project owner mismatch, Owner:", projectData.OwnerID, "AuthUser:", authenticatedUser.ID)
			rw.Error("Please contact your project owner to add Access Keys.", http.StatusForbidden)
			return
		}

		/*
			Generate 16 byte random key to encrypt `accessKeySecret`
		*/
		key, err := encryption.GenerateRandomSecretKey(16)
		if err != nil {
			log_pa.Error("ERROR unable to generate secret key", err)
			rw.Error("Something went wrong, Please try again later", http.StatusInternalServerError)
			return
		}
		cipherText, err := encryption.Encrypt(request.AccessKeySecret, key)
		if err != nil {
			log_pa.Error("ERROR unable to encrypt access key", err)
			rw.Error("Something went wrong, Please try again later", http.StatusInternalServerError)
			return
		}
		request.AccessKeySecret = cipherText
		request.RotationSecretKey = key
		result, err := provider_models.Create(db, request)
		if err != nil {
			log_pa.Error("ERROR unable to create data", err)
			rw.Error("Something went wrong, Please try again later", http.StatusInternalServerError)
			return
		}
		rw.Success(result, http.StatusOK)
		return
	}
}

func (c *ProviderAccountsController) GetById(db *database.DB) http.HandlerFunc {
	log_pa.SetName(c.Name() + ".GetById")
	return func(w http.ResponseWriter, r *http.Request) {
		rw := middleware.CustomResponseWriter(w)
		authenticatedUser := rw.SessionUser
		id := mux.Vars(r)["id"]
		isIDValidUUID := helpers.IsValidUUID(id)
		if !isIDValidUUID {
			log_pa.Error("Invalid provider account ID")
			rw.Error("Invalid provider account ID provided", http.StatusUnprocessableEntity)
			return
		}
		providerAccountId, err := uuid.Parse(id)
		if err != nil {
			rw.Error("Invalid provider account ID provided", http.StatusUnprocessableEntity)
			return
		}
		_, verificationErr := authentication.VerifyProjectAccess(db, authenticatedUser.ID, types.VerifyProjectAccessInput{
			ProviderAccountID: providerAccountId,
		})
		if verificationErr != nil {
			log_pa.Error("Project verification failed", verificationErr)
			rw.Error(verificationErr.Error(), http.StatusForbidden)
			return
		}
		result, err := provider_models.GetById(db, providerAccountId)
		if err != nil {
			log_pa.Error("ERROR", err)
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
	log_pa.SetName(c.Name() + ".GetByProject")
	return func(w http.ResponseWriter, r *http.Request) {
		rw := middleware.CustomResponseWriter(w)
		authenticatedUser := rw.SessionUser
		projectId := r.URL.Query().Get("projectId")
		if projectId == "" {
			rw.Error("Invalid project ID", http.StatusBadRequest)
			return
		}
		/**
		verification to make sure the project
		can be accessed by user.
		*/
		project, err := authentication.VerifyProjectAccess(db, authenticatedUser.ID, types.VerifyProjectAccessInput{
			ProjectID: projectId,
		})
		if err != nil {
			log_pa.Error("invalid project uuid provided", err)
			rw.Error(err.Error(), http.StatusBadRequest)
			return
		}
		log_pa.Log("fetching data for projectId:", project.ProjectAccessDetails.Projects.ID)
		result, err := provider_models.GetByProjectId(db, project.ProjectAccessDetails.Projects.ID)
		if err != nil {
			log_pa.Error("ERROR", err)
			rw.Error("Unable to fetch provider accounts", http.StatusInternalServerError)
			return
		}
		rw.Success(result, http.StatusOK)
		return
	}
}

func (c *ProviderAccountsController) UpdateProviderAccount(db *database.DB) http.HandlerFunc {
	log_pa.SetName(c.Name() + ".UpdateProviderAccount")
	return func(w http.ResponseWriter, r *http.Request) {
		rw := middleware.CustomResponseWriter(w)
		authenticatedUser := rw.SessionUser
		id := mux.Vars(r)["id"]
		isIDValidUUID := helpers.IsValidUUID(id)
		if !isIDValidUUID {
			log_pa.Error("Invalid provider account ID")
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
			log_pa.Error("Error reading request body:", err)
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
			log_pa.Error("Error parsing request body:", err)
			rw.Error(custom_errors.UnknownError.Error(), http.StatusUnprocessableEntity)
			return
		}
		verifiedData, err := authentication.VerifyProjectAccess(db, authenticatedUser.ID, types.VerifyProjectAccessInput{
			ProviderAccountID: providerAccountId,
		})
		if err != nil {
			log_pa.Error("Project verification failed", err)
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
		if *verifiedData.ProjectAccessDetails.OwnerID != authenticatedUser.ID {
			log_pa.Error("ERROR project owner mismatch, Owner:", verifiedData.ProjectAccessDetails.OwnerID, "AuthUser:", authenticatedUser.ID)
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
