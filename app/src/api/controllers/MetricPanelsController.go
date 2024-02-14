package controllers

import (
	"cloudview/app/src/api/authentication"
	custom_errors "cloudview/app/src/api/errors"
	"cloudview/app/src/api/middleware"
	"cloudview/app/src/database"
	"cloudview/app/src/helpers"
	"cloudview/app/src/models"
	metric_panels_model "cloudview/app/src/models/metric_panels"
	"cloudview/app/src/permissions"
	"cloudview/app/src/types"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

func (c *MetricPanelsController) CreateMetricPanel(db *database.DB) http.HandlerFunc {
	c.Logger.SetName(c.Name + ".CreateMetricPanel")
	return func(w http.ResponseWriter, r *http.Request) {
		rw := middleware.CustomResponseWriter(w)
		authUserId := rw.SessionUser.ID
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			c.Logger.Error("Error reading request body:", err)
			rw.Error("Bad request", http.StatusUnprocessableEntity)
			return
		}
		var request models.MetricPanels
		if err := json.Unmarshal(body, &request); err != nil {
			c.Logger.Error("Error parsing request body:", err)
			rw.Error(custom_errors.UnknownError.Error(), http.StatusUnprocessableEntity)
			return
		}
		if request.Panels == nil {
			rw.Error("Required `panels` field", http.StatusBadRequest)
			return
		}
		if request.Name == "" || request.InstanceID == "" {
			rw.Error("Required `name` and `instanceId` field", http.StatusBadRequest)
			return
		}
		verifiedData, err := authentication.VerifyProjectAccess(db, authUserId, types.VerifyProjectAccessInput{
			ProviderAccountID: request.ProviderAccountID,
		})
		if err != nil {
			c.Logger.Error("MetricPanelsController.CreateMetricPanel: ERROR", err)
			rw.Error(err.Error(), http.StatusBadRequest)
			return
		}
		perms := verifiedData.ProjectAccessDetails.ProjectMembers.Permissions
		if perms == "" {
			perms = "0"
		}
		canContinue := permissions.VerifyPermissions([]string{
			permissions.MEMBER_MANAGE_METRICS_PANEL,
		}, perms)

		if !canContinue {
			rw.ErrorMessage = "You do not have permissions to 'create' or 'modify' metric panels. Please contact your administrator."
			rw.Forbidden()
			return
		}
		result, err := metric_panels_model.Create(db, request)
		if err != nil {
			rw.Error(err.Error(), http.StatusInternalServerError)
			return
		}
		rw.Success(result, http.StatusOK)
		return
	}
}

func (c *MetricPanelsController) UpdateMetricPanel(db *database.DB) http.HandlerFunc {
	c.Logger.SetName(c.Logger.Name + ".UpdateMetricPanel")
	return func(w http.ResponseWriter, r *http.Request) {
		rw := middleware.CustomResponseWriter(w)
		authUser := rw.SessionUser
		id := mux.Vars(r)["id"]
		isUUIDValid, uuid := helpers.IsValidUUID(id)
		if !isUUIDValid {
			rw.Error("Invalid project ID provided", http.StatusBadRequest)
			return
		}
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			c.Logger.Error("Error reading request body:", err)
			rw.Error("Bad request", http.StatusUnprocessableEntity)
			return
		}
		var request models.MetricPanels
		if err := json.Unmarshal(body, &request); err != nil {
			c.Logger.Error("Error parsing request body:", err)
			rw.Error(custom_errors.UnknownError.Error(), http.StatusUnprocessableEntity)
			return
		}
		verifiedData, err := authentication.VerifyProjectAccess(db, authUser.ID, types.VerifyProjectAccessInput{
			ProviderAccountID: request.ProviderAccountID,
		})
		if err != nil {
			c.Logger.Error("error verifying project access", err)
			rw.Error(err.Error(), http.StatusBadRequest)
			return
		}
		perms := verifiedData.ProjectAccessDetails.ProjectMembers.Permissions
		if perms == "" {
			perms = "0"
		}
		canContinue := permissions.VerifyPermissions([]string{
			permissions.MEMBER_MANAGE_METRICS_PANEL,
		}, perms)

		if !canContinue {
			rw.ErrorMessage = "You do not have permissions to 'modify' metric panels. Please contact your administrator."
			rw.Forbidden()
			return
		}
		if err := metric_panels_model.Update(db, *uuid, request); err != nil {
			c.Logger.Error("Unable to edit metric_panels", err)
			rw.Error(err.Error(), http.StatusInternalServerError)
			return
		}
		rw.Success("Data updated successfully", http.StatusOK)
		return
	}
}

/*
TODO - Add pagination
*/
func (c *MetricPanelsController) GetPanels(db *database.DB) http.HandlerFunc {
	c.Logger.SetName(c.Name + ".GetPanels")
	return func(w http.ResponseWriter, r *http.Request) {
		rw := middleware.CustomResponseWriter(w)
		authUser := rw.SessionUser
		providerAccountId := r.URL.Query().Get("providerAccountId")
		verifyData, err := authentication.VerifyProjectAccess(db, authUser.ID, types.VerifyProjectAccessInput{
			ProviderAccountID: providerAccountId,
		})
		if err != nil {
			c.Logger.Error("Unable to verify project access", err)
			rw.Forbidden()
			return
		}
		result, err := metric_panels_model.GetByProviderAccount(db, verifyData.ProviderAccount.ID)
		if err != nil {
			c.Logger.Error("Unable to fetch data", err)
			rw.Error("unable to fetch metric panels", http.StatusInternalServerError)
			return
		}
		rw.Success(result, http.StatusOK)
		return
	}
}
