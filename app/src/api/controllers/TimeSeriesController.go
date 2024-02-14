package controllers

import (
	"cloudview/app/src/api/authentication"
	custom_errors "cloudview/app/src/api/errors"
	"cloudview/app/src/api/middleware"
	"cloudview/app/src/database"
	"cloudview/app/src/helpers"
	models "cloudview/app/src/models/structs"
	time_series_model "cloudview/app/src/models/time_series"
	"cloudview/app/src/types"
	"cloudview/app/src/utility"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type series struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}
type tsBody struct {
	Series        *[]series `json:"series"`
	Timestamp     int64     `json:"timestamp"`
	MetricPanelID *string   `json:"metricPanelId"` // MetricPanelID
	Type          string    `json:"type"`
}

type getByMetricIdqueryParams struct {
	MetricID          string `json:"metricId" query:"metricId"`
	ProviderAccountID string `json:"providerAccountId" query:"providerAccountId"`
	MetricType        string `json:"type" query:"type"`
	Limit             int    `json:"limit" query:"limit"`
	Page              int    `json:"page" query:"page"`
	StartTime         int64  `json:"startTime,omitempty" query:"startTime"`
	EndTime           int64  `json:"endTime,omitempty" query:"endTime"`
	Period            int    `json:"period,omitempty" query:"period"`
	Aggregation       string `json:"aggregation,omitempty" query:"aggregation"`
}

/*
These values are enums from `time_series` table. If more enums
are added make sure to update this array.
*/
var metricsToLog = []string{"CPU_USAGE", "RAM_USAGE", "DISK_OPERATIONS"}

/*
Record usage metrics from binaries. You can only report
metrics from a 'system' user who has 'report_metrics' permissions.

TODO:
1. Route is called every 5 seconds to save data. Need to add rate
limits.
2. Authenticate these requests by store a auth_key in config.yaml and DB compare the values.

NOTICE: The timestamp in request body must be in ms.
*/
func (c *TimeSeriesController) SaveSeries(db *database.DB) http.HandlerFunc {
	c.Logger.SetName(c.Name + ".SaveSeries")
	return func(w http.ResponseWriter, r *http.Request) {
		rw := middleware.RegisterResponses(w)
		/*
			Auth headers must contain a system user with required permissions to
			save metric series.
		*/
		// system := rw.SessionUser
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			c.Logger.Error("Error reading request body:", err)
			rw.Error("Bad request", http.StatusUnprocessableEntity)
			return
		}
		var request tsBody
		if err := json.Unmarshal(body, &request); err != nil {
			c.Logger.Error("Error parsing request body:", err)
			rw.Error(custom_errors.UnknownError.Error(), http.StatusUnprocessableEntity)
			return
		}
		if request.MetricPanelID == nil {
			rw.Error("Invalid `metricPanelId` provided.", http.StatusUnprocessableEntity)
			return
		}
		if request.Series == nil {
			rw.Error("Missing time series data, field 'series' is required", http.StatusBadRequest)
			return
		}
		validType, enumStr := utility.ContainsString(metricsToLog, request.Type)
		if !validType {
			c.Logger.Error("Invalid enum type, got: ", request.Type)
			msg := fmt.Sprintf("Field 'type' must be one of (%s), got: %s", enumStr, request.Type)
			rw.Error(msg, http.StatusUnprocessableEntity)
			return
		}

		valid, uid := helpers.IsValidUUID(*request.MetricPanelID)
		if !valid {
			rw.Error("Invalid 'metricPanelId' provided. Expected type uuid.", http.StatusBadRequest)
			return
		}

		/*	Optimization:
			We round off the timestampto the lowest whole hour (hh:00:00).
			i.e. Each row will contain metric points of 1 hour (3600 points).
			This will helps us determine if we need to create a new row or append to
			existing row.
			There will be 24 rows created in 1 day.
		*/
		ts := c.roundEpochToHour(request.Timestamp)

		tsData := c.prepareSeriesData(*request.Series)
		c.Logger.Log(fmt.Sprintf("upserting time_series data: %+v\n", *tsData))
		err = time_series_model.UpsertTimeSeries(db, models.TimeSeries{
			TruncatedTimestamp: ts,
			MetricPanelID:      *uid,
			Type:               models.MetricType(request.Type),
			Series:             tsData,
		})
		if err != nil {
			c.Logger.Error("unable to insert data:", err)
			rw.Error("unknown error occured while inserting time series", http.StatusInternalServerError)
			return
		}
		rw.Success("ok", http.StatusOK)
		return
	}
}

func (c *TimeSeriesController) GetByMetricId(db *database.DB) http.HandlerFunc {
	c.Logger.SetName(c.Name + ".GetByMetricId")
	return func(w http.ResponseWriter, r *http.Request) {
		rw := middleware.CustomResponseWriter(w)
		var queryP getByMetricIdqueryParams
		err := utility.ExtractQueryParams(r, &queryP)
		if err != nil {
			c.Logger.Error("Unable to extract query params", err)
			rw.Error("Unknown error occured", http.StatusInternalServerError)
			return
		}
		c.Logger.Log(fmt.Sprintf("fetching time series metrics of type: %s, for metricId: %s", queryP.MetricType, queryP.MetricID))
		valid, enumStr := utility.ContainsString(metricsToLog, queryP.MetricType)
		if !valid {
			rw.Error(fmt.Sprintf("query 'type' must be one of: %s, got: %s", enumStr, queryP.MetricType), http.StatusBadRequest)
			return
		}
		valid, uid := helpers.IsValidUUID(queryP.MetricID)
		if !valid {
			rw.Error("Invalid `metricId` provided", http.StatusBadRequest)
			return
		}
		p_valid, p_uid := helpers.IsValidUUID(queryP.ProviderAccountID)
		if !p_valid {
			rw.Error("Invalid `providerAccountId` field.", http.StatusBadRequest)
			return
		}
		_, err = authentication.VerifyProjectAccess(db, rw.SessionUser.ID, types.VerifyProjectAccessInput{
			ProviderAccountID: queryP.ProviderAccountID,
		})
		if err != nil {
			c.Logger.Error(fmt.Sprintf("project access verification failed for user: %s, provider account: %s", rw.SessionUser.ID, queryP.ProviderAccountID))
			rw.Forbidden()
			return
		}
		params := time_series_model.GetByMetricIdInput{
			MetricPanelID:     *uid,
			MetricType:        queryP.MetricType,
			ProviderAccountID: *p_uid,
			Aggregation:       &queryP.Aggregation,
			StartTime:         &queryP.StartTime,
			EndTime:           &queryP.EndTime,
			Period:            &queryP.Period,
		}
		if *params.StartTime > 0 {
			if *params.EndTime == 0 {
				c.Logger.Error("Invalid `startTime` or `endTime` specified specified.")
				rw.Error("Invalid 'endTime' field. Expected valid date timestamp.", http.StatusBadRequest)
				return
			}
		} else {
			/*
				By default - we will fetch metric data at an interval of
				2mins. Implied as real-time
			*/
			*params.StartTime = pastDuration(2)
			*params.EndTime = pastDuration(0)
		}
		result, err := time_series_model.GetByMetricId(db, params)
		if err != nil {
			c.Logger.Error("unable to fetch data", err)
			rw.Error("Unknown error occured", http.StatusInternalServerError)
			return
		}
		c.Logger.Log("success")
		rw.Success(result, http.StatusOK)
		return
	}
}

func pastDuration(minutes uint8) int64 {
	now := time.Now()
	then := now.Add(time.Duration(-minutes) * time.Minute)
	return then.Unix()
}

func (c *TimeSeriesController) roundEpochToHour(ms int64) int64 {
	t := time.Unix(ms, 0)
	t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())
	return t.Unix()
}

func (c *TimeSeriesController) prepareSeriesData(data []series) *[][]float64 {
	result := make([][]float64, len(data))
	for i, v := range data {
		result[i] = []float64{float64(v.Timestamp), v.Value}
	}

	return &result
}
