package controllers

import (
	custom_errors "cloudview/app/src/api/errors"
	"cloudview/app/src/api/middleware"
	"cloudview/app/src/api/middleware/logger"
	"cloudview/app/src/database"
	"cloudview/app/src/models"
	"encoding/json"
	"io"
	"net/http"
)

func (c *MetricPanelsController) CreateMetricPanel(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := middleware.RegisterResponses(w)
		_, err := rw.User(db, r)
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
