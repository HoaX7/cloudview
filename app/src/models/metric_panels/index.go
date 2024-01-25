package models

import (
	"cloudview/app/dbschema/cloudview/public/table"
	custom_errors "cloudview/app/src/api/errors"
	"cloudview/app/src/api/middleware/logger"
	"cloudview/app/src/database"
	"cloudview/app/src/models"
	"errors"
	"reflect"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
)

func _create(db *database.DB, data models.MetricPanels) (models.MetricPanels, error) {
	columnList := postgres.ColumnList{
		table.MetricPanels.Name,
		table.MetricPanels.Panels,
		table.MetricPanels.ProviderAccountID,
	}
	if data.Description != nil {
		columnList = append(columnList, table.MetricPanels.Description)
	}

	stmt := table.ProjectMembers.INSERT(columnList).MODEL(data).
		RETURNING(table.MetricPanels.AllColumns)

	queryString, args := stmt.Sql()
	logger.Logger.Log("Inserting into Metric Panels table with data: ", queryString, args)

	result := models.MetricPanels{}
	rows, err := db.Postgres.Query(queryString, args...)
	if err != nil {
		logger.Logger.Error("models.metric_panels.Create: ERROR", err)
		return result, custom_errors.DBErrors(err)
	}
	if rows.Next() {
		if err := rows.Scan(&result.ID,
			&result.Name, &result.Description,
			&result.Panels, &result.ProviderAccountID,
			&result.Metadata,
			&result.IsDeleted,
			&result.CreatedAt,
			&result.UpdatedAt); err != nil {

			logger.Logger.Error("models.metric_panels.Create: ERROR", err)
			return result, err
		}
	}
	return result, nil
}

func _update(db *database.DB, id uuid.UUID, data models.MetricPanels) error {
	columnsList := postgres.ColumnList{}
	if data.Name != "" {
		columnsList = append(columnsList, table.MetricPanels.Name)
	}
	if data.Description != nil {
		columnsList = append(columnsList, table.MetricPanels.Description)
	}
	if data.IsDeleted != nil && reflect.ValueOf(*data.IsDeleted).Kind() == reflect.Bool {
		columnsList = append(columnsList, table.MetricPanels.IsDeleted)
	}
	if data.Metadata != nil {
		columnsList = append(columnsList, table.MetricPanels.Metadata)
	}
	if data.Panels != nil {
		columnsList = append(columnsList, table.MetricPanels.Panels)
	}
	if len(columnsList) <= 0 {
		logger.Logger.Log("models.metric_panels.Update: nothing to update")
		return nil
	}
	stmt := table.MetricPanels.UPDATE(columnsList).
		MODEL(data).WHERE(table.MetricPanels.ID.EQ(postgres.UUID(id)))

	logger.Logger.Log("models.metric_panels.Update: updating", stmt.DebugSql())
	_, err := stmt.Exec(db.Postgres)
	if err != nil {
		logger.Logger.Error("models.metric_panels.Update: ERROR", err)
		return custom_errors.DBErrors(err)
	}
	return nil
}

func _getById(db *database.DB, id uuid.UUID) (models.MetricPanels, error) {
	stmt := table.MetricPanels.SELECT(table.MetricPanels.AllColumns).
		WHERE(postgres.AND(
			table.MetricPanels.ID.EQ(postgres.UUID(id)),
			table.MetricPanels.IsDeleted.EQ(postgres.Bool(false)),
		))

	var result models.MetricPanels
	if err := stmt.Query(db.Postgres, &result); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return result, custom_errors.NoDataFound
		}
		return result, err
	}

	return result, nil
}

func _getByProviderAccount(db *database.DB, providerAccount uuid.UUID) ([]models.MetricPanels, error) {
	stmt := table.MetricPanels.SELECT(table.MetricPanels.AllColumns).
		WHERE(postgres.AND(
			table.MetricPanels.ProviderAccountID.EQ(postgres.UUID(providerAccount)),
			table.MetricPanels.IsDeleted.EQ(postgres.Bool(false)),
		))

	var result []models.MetricPanels
	if err := stmt.Query(db.Postgres, &result); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return result, custom_errors.NoDataFound
		}
		return result, err
	}

	return result, nil
}
