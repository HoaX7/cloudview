package models

import (
	"cloudview/app/dbschema/cloudview/public/table"
	custom_errors "cloudview/app/src/api/errors"
	"cloudview/app/src/api/middleware/logger"
	"cloudview/app/src/database"
	"encoding/json"
	"errors"
	"reflect"
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
)

type Services struct {
	ID                uuid.UUID        `sql:"primary_key" json:"id,omitempty"`
	Name              string           `json:"name,omitempty"`
	Description       string           `json:"description,omitempty"`
	AccessKeyID       string           `json:"accessKeyId,omitempty"`
	AccessKeySecret   string           `json:"accessKeySecret,omitempty"`
	RotationSecretKey string           `json:"rotationSecretKey,omitempty"`
	Provider          string           `json:"provider,omitempty"`
	ProjectID         uuid.UUID        `json:"projectId,omitempty" jet:"nullable"`
	Metadata          *json.RawMessage `json:"metadata,omitempty" jet:"type=jsonb,nullable"`
	IsDeleted         *bool            `json:"isDeleted,omitempty"`
	CreatedAt         time.Time        `json:"createdAt,omitempty"`
	UpdatedAt         time.Time        `json:"updatedAt,omitempty"`
}

/*
*
TODO - Implement rotation secret key logic
to hash `accessSecretKey`
*/
func Create(db *database.DB, data Services) (Services, error) {
	stmt := table.Services.INSERT(
		table.Services.Name,
		table.Services.Description,
		table.Services.AccessKeyID,
		table.Services.AccessKeySecret,
		table.Services.RotationSecretKey,
		table.Services.Provider,
		table.Services.ProjectID,
	).MODEL(data).RETURNING(table.Services.AllColumns)

	queryString, args := stmt.Sql()
	logger.Logger.Log("Inserting into services table with data: ", queryString, args)

	result := Services{}
	rows, err := db.Postgres.Query(queryString, args...)
	if err != nil {
		logger.Logger.Error("models.Services.Create: ERROR", err)
		return result, custom_errors.DBErrors(err)
	}
	if rows.Next() {
		if err := rows.Scan(&result.ID, &result.Name,
			&result.Description,
			&result.AccessKeyID, &result.AccessKeySecret,
			&result.RotationSecretKey,
			&result.Provider, &result.ProjectID,
			&result.Metadata,
			&result.IsDeleted,
			&result.CreatedAt,
			&result.UpdatedAt); err != nil {

			logger.Logger.Error("models.Services.Create: ERROR", err)
			return result, err
		}
	}
	return result, nil
}

/*
*
TODO - Add pagination
*/
func GetByProjectId(db *database.DB, projectId uuid.UUID) ([]Services, error) {
	stmt := table.Services.SELECT(table.Services.ID, table.Services.Name,
		table.Services.Description, table.Services.Provider,
		table.Services.ProjectID,
		table.Services.CreatedAt, table.Services.UpdatedAt).
		WHERE(postgres.AND(
			table.Services.ProjectID.EQ(postgres.UUID(projectId)),
			table.Services.IsDeleted.EQ(postgres.Bool(false)),
		))

	logger.Logger.Log(stmt.DebugSql())
	var result []Services
	if err := stmt.Query(db.Postgres, &result); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return result, nil
		}
		return result, err
	}

	return result, nil
}

func GetById(db *database.DB, id uuid.UUID) (Services, error) {
	stmt := table.Services.SELECT(table.Services.AllColumns).
		WHERE(postgres.AND(
			table.Services.ID.EQ(postgres.UUID(id)),
			table.Services.IsDeleted.EQ(postgres.Bool(false)),
		))

	var result Services
	if err := stmt.Query(db.Postgres, &result); err != nil {
		logger.Logger.Error("models.services.GetById: ERROR", err)
		if errors.Is(err, qrm.ErrNoRows) {
			return result, nil
		}
		return result, err
	}
	return result, nil
}

func Update(db *database.DB, id uuid.UUID, data Services) error {
	/**
	When you use `unmarshal` to set struct data from json body
	default values for non-specified fields are autoamtically added
	by go. If the field is empty string we should not be updating the value.
	therefore, build dynamic query.
	*/

	columnsList := postgres.ColumnList{}
	if data.Name != "" {
		columnsList = append(columnsList, table.Services.Name)
	}
	if data.Description != "" {
		columnsList = append(columnsList, table.Services.Description)
	}
	if data.AccessKeyID != "" {
		columnsList = append(columnsList, table.Services.AccessKeyID)
	}
	if data.AccessKeySecret != "" {
		columnsList = append(columnsList, table.Services.AccessKeySecret)
	}
	if data.Provider != "" {
		columnsList = append(columnsList, table.Services.Provider)
	}
	if data.IsDeleted != nil && reflect.ValueOf(*data.IsDeleted).Kind() == reflect.Bool {
		columnsList = append(columnsList, table.Projects.IsDeleted)
	}
	if len(columnsList) <= 0 {
		logger.Logger.Log("models.services.Update: nothing to update")
		return nil
	}
	stmt := table.Services.UPDATE(columnsList).
		MODEL(data).WHERE(table.Services.ID.EQ(postgres.UUID(id)))

	logger.Logger.Log("models.services.Update: updating", stmt.DebugSql())
	_, err := stmt.Exec(db.Postgres)
	if err != nil {
		logger.Logger.Error("models.services.Update: ERROR", err)
		return custom_errors.DBErrors(err)
	}
	return nil
}
