package models

import (
	custom_errors "cloudview/app/src/api/errors"
	"cloudview/app/src/api/middleware/logger"
	"cloudview/app/src/database"
	models "cloudview/app/src/models/structs"
	"errors"
)

func UpsertTimeSeries(db *database.DB, data models.TimeSeries) error {
	res, err := _getExistingTs(db, GetExistingTsInput{
		TruncatedTimestamp: data.TruncatedTimestamp,
		MetricPanelID:      data.MetricPanelID,
		MetricType:         string(data.Type),
	})
	if err != nil {
		if errors.Is(err, custom_errors.NoDataFound) {
			logger.Logger.Log("Creating new time series row..")
			return _create(db, data)
		}
		logger.Logger.Error("unable to fetch existing rows..", err)
		return err
	}
	data.ID = res.ID
	logger.Logger.Log("appending time_series data for id:", res.ID)
	if err := _appendTsData(db, data); err != nil {
		logger.Logger.Error("unable to append ts data: ERROR", err)
		return err
	}
	return nil
}

func GetByMetricId(db *database.DB, opts GetByMetricIdInput) ([]models.TimeSeries, error) {
	res, err := _getByMetricId(db, opts)
	if err != nil {
		if errors.Is(err, custom_errors.NoDataFound) {
			return nil, nil
		}
		return nil, err
	}
	return res, nil
}
