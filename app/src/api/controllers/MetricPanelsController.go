package controllers

import (
	"cloudview/app/src/api/authentication"
	custom_errors "cloudview/app/src/api/errors"
	"cloudview/app/src/api/middleware"
	"cloudview/app/src/api/middleware/logger"
	"cloudview/app/src/database"
	"cloudview/app/src/models"
	metric_panels_model "cloudview/app/src/models/metric_panels"
	"cloudview/app/src/permissions"
	"cloudview/app/src/types"
	"encoding/json"
	"io"
	"net/http"
)

func (c *MetricPanelsController) CreateMetricPanel(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := middleware.RegisterResponses(w)
		authUser, err := rw.User(db, r)
		if err != nil {
			logger.Logger.Error("MetricPanelsController.Create: ERROR", err)
			rw.Forbidden()
			return
		}
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			logger.Logger.Error("MetricPanelsController.Create: Error reading request body:", err)
			rw.Error("Bad request", http.StatusUnprocessableEntity)
			return
		}
		var request models.MetricPanels
		if err := json.Unmarshal(body, &request); err != nil {
			logger.Logger.Error("MetricPanelsController.Create: Error parsing request body:", err)
			rw.Error(custom_errors.UnknownError.Error(), http.StatusUnprocessableEntity)
			return
		}
		if request.Panels == nil {
			rw.Error("Required `panels` field", http.StatusBadRequest)
			return
		}
		if request.Name == "" {
			rw.Error("Required `name` field", http.StatusBadRequest)
			return
		}
		verifiedData, err := authentication.VerifyProjectAccess(db, authUser.ID, types.VerifyProjectAccessInput{
			ProviderAccountID: request.ProviderAccountID,
		})
		if err != nil {
			logger.Logger.Error("MetricPanelsController.CreateMetricPanel: ERROR", err)
			rw.Error(err.Error(), http.StatusBadRequest)
			return
		}
		perms := verifiedData.ProjectAccessDetails.ProjectMembers.Permissions
		if perms == "" {
			perms = "0"
		}
		canContinue := permissions.VerifyPermissions([]string{
			permissions.MANAGE_METRICS_PANEL,
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
	return func(w http.ResponseWriter, r *http.Request) {
		rw := middleware.RegisterResponses(w)
		authUser, err := rw.User(db, r)
		if err != nil {
			logger.Logger.Error("MetricPanelsController.Create: ERROR", err)
			rw.Forbidden()
			return
		}
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			logger.Logger.Error("MetricPanelsController.Create: Error reading request body:", err)
			rw.Error("Bad request", http.StatusUnprocessableEntity)
			return
		}
		var request models.MetricPanels
		if err := json.Unmarshal(body, &request); err != nil {
			logger.Logger.Error("MetricPanelsController.Create: Error parsing request body:", err)
			rw.Error(custom_errors.UnknownError.Error(), http.StatusUnprocessableEntity)
			return
		}
	}
}
