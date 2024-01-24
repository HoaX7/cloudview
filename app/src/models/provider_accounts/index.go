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

type AccessType string // Access Keys, Cross Account Role

func _getById(db *database.DB, id uuid.UUID) (models.ProviderAccountWithProject, error) {
	stmt := table.ProviderAccounts.SELECT(table.ProviderAccounts.ID, table.ProviderAccounts.Name,
		table.ProviderAccounts.Description, table.ProviderAccounts.Provider,
		table.ProviderAccounts.AccountID,
		table.ProviderAccounts.AccessRole, table.ProviderAccounts.Metadata,
		table.Projects.Name, table.Projects.ID).
		FROM(table.ProviderAccounts.LEFT_JOIN(table.Projects,
			table.ProviderAccounts.ProjectID.EQ(table.Projects.ID))).
		WHERE(postgres.AND(
			table.ProviderAccounts.ID.EQ(postgres.UUID(id)),
			table.ProviderAccounts.IsDeleted.EQ(postgres.Bool(false)),
		))

	logger.Logger.Log("models.provider_accounts.GetById", stmt.DebugSql())
	var result models.ProviderAccountWithProject
	if err := stmt.Query(db.Postgres, &result); err != nil {
		logger.Logger.Error("models.provideraccounts.GetById: ERROR", err)
		if errors.Is(err, qrm.ErrNoRows) {
			return result, nil
		}
		return result, err
	}
	return result, nil
}

func _getByIdForSDK(db *database.DB, id uuid.UUID) (models.ProviderAccounts, error) {
	stmt := table.ProviderAccounts.SELECT(table.ProviderAccounts.AccessKeyID, table.ProviderAccounts.AccessKeySecret,
		table.ProviderAccounts.AccessRole, table.ProviderAccounts.AccountID, table.ProviderAccounts.FeatureAccessPermission,
		table.ProviderAccounts.RotationSecretKey, table.ProviderAccounts.Provider).
		WHERE(postgres.AND(
			table.ProviderAccounts.ID.EQ(postgres.UUID(id)),
			table.ProviderAccounts.IsDeleted.EQ(postgres.Bool(false)),
		))

	logger.Logger.Log("models.provider_accounts._getByIdForSDK", stmt.DebugSql())
	var result models.ProviderAccounts
	if err := stmt.Query(db.Postgres, &result); err != nil {
		logger.Logger.Error("models.provideraccounts.GetById: ERROR", err)
		if errors.Is(err, qrm.ErrNoRows) {
			return result, nil
		}
		return result, err
	}
	return result, nil
}

func _update(db *database.DB, id uuid.UUID, data models.ProviderAccounts) error {
	/**
	When you use `unmarshal` to set struct data from json body
	default values for non-specified fields are autoamtically added
	by go. If the field is empty string we should not be updating the value.
	therefore, build dynamic query.
	*/

	columnsList := postgres.ColumnList{}
	if data.Name != "" {
		columnsList = append(columnsList, table.ProviderAccounts.Name)
	}
	if data.Description != "" {
		columnsList = append(columnsList, table.ProviderAccounts.Description)
	}
	if data.AccessKeyID != "" {
		columnsList = append(columnsList, table.ProviderAccounts.AccessKeyID)
	}
	if data.AccessKeySecret != "" {
		columnsList = append(columnsList, table.ProviderAccounts.AccessKeySecret)
	}
	if data.Provider != "" {
		columnsList = append(columnsList, table.ProviderAccounts.Provider)
	}
	if data.IsDeleted != nil && reflect.ValueOf(*data.IsDeleted).Kind() == reflect.Bool {
		columnsList = append(columnsList, table.ProviderAccounts.IsDeleted)
	}
	if data.Metadata != nil {
		columnsList = append(columnsList, table.ProviderAccounts.Metadata)
	}
	if len(columnsList) <= 0 {
		logger.Logger.Log("models.provideraccounts.Update: nothing to update")
		return nil
	}
	stmt := table.ProviderAccounts.UPDATE(columnsList).
		MODEL(data).WHERE(table.ProviderAccounts.ID.EQ(postgres.UUID(id)))

	logger.Logger.Log("models.provideraccounts.Update: updating", stmt.DebugSql())
	_, err := stmt.Exec(db.Postgres)
	if err != nil {
		logger.Logger.Error("models.provideraccounts.Update: ERROR", err)
		return custom_errors.DBErrors(err)
	}
	return nil
}

/*
*
TODO - Implement rotation secret key logic
to hash `accessSecretKey`
*/
func Create(db *database.DB, data models.ProviderAccounts) (models.ProviderAccounts, error) {
	stmt := table.ProviderAccounts.INSERT(
		table.ProviderAccounts.Name,
		table.ProviderAccounts.Description,
		table.ProviderAccounts.AccessKeyID,
		table.ProviderAccounts.AccessKeySecret,
		table.ProviderAccounts.RotationSecretKey,
		table.ProviderAccounts.Provider,
		table.ProviderAccounts.ProjectID,
	).MODEL(data).RETURNING(table.ProviderAccounts.AllColumns)

	queryString, args := stmt.Sql()
	logger.Logger.Log("Inserting into provideraccounts table with data: ", queryString, args)

	result := models.ProviderAccounts{}
	rows, err := db.Postgres.Query(queryString, args...)
	if err != nil {
		logger.Logger.Error("models.ProviderAccounts.Create: ERROR", err)
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

			logger.Logger.Error("models.ProviderAccounts.Create: ERROR", err)
			return result, err
		}
	}
	return result, nil
}

/*
*
TODO - Add pagination
*/
func GetByProjectId(db *database.DB, projectId uuid.UUID) ([]models.ProviderAccounts, error) {
	stmt := table.ProviderAccounts.SELECT(table.ProviderAccounts.ID, table.ProviderAccounts.Name,
		table.ProviderAccounts.Description, table.ProviderAccounts.Provider,
		table.ProviderAccounts.ProjectID, table.ProviderAccounts.AccountID,
		table.ProviderAccounts.AccessRole,
		table.ProviderAccounts.CreatedAt, table.ProviderAccounts.UpdatedAt).
		WHERE(postgres.AND(
			table.ProviderAccounts.ProjectID.EQ(postgres.UUID(projectId)),
			table.ProviderAccounts.IsDeleted.EQ(postgres.Bool(false)),
		))

	logger.Logger.Log(stmt.DebugSql())
	var result []models.ProviderAccounts
	if err := stmt.Query(db.Postgres, &result); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return result, nil
		}
		return result, err
	}

	return result, nil
}
