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

func _create(db *database.DB, data models.ProjectMembers) (models.ProjectMembers, error) {
	stmt := table.ProjectMembers.INSERT(
		table.ProjectMembers.ProjectID,
		table.ProjectMembers.UserID,
		table.ProjectMembers.IsOwner,
	).MODEL(data).RETURNING(table.ProjectMembers.AllColumns)

	queryString, args := stmt.Sql()
	logger.Logger.Log("Inserting into Project Members table with data: ", queryString, args)

	result := models.ProjectMembers{}
	rows, err := db.Postgres.Query(queryString, args...)
	if err != nil {
		logger.Logger.Error("models.project_members.Create: ERROR", err)
		return result, custom_errors.DBErrors(err)
	}
	if rows.Next() {
		if err := rows.Scan(&result.ID,
			&result.ProjectID, &result.UserID,
			&result.IsOwner, &result.IsActive,
			&result.Permissions, &result.Metadata,
			&result.IsDeleted,
			&result.CreatedAt,
			&result.UpdatedAt); err != nil {

			logger.Logger.Error("models.project_members.Create: ERROR", err)
			return result, err
		}
	}
	return result, nil
}

func _getByIdAndUserId(db *database.DB, id uuid.UUID, userId uuid.UUID) (models.ProjectMembers, error) {
	stmt := table.ProjectMembers.SELECT(table.ProjectMembers.AllColumns).
		WHERE(postgres.AND(
			table.ProjectMembers.ID.EQ(postgres.UUID(id)),
			table.ProjectMembers.IsDeleted.EQ(postgres.Bool(false)),
			table.ProjectMembers.UserID.EQ(postgres.UUID(userId)),
		))

	logger.Logger.Log("models.project_members.GetByIdAndUserId", stmt.DebugSql())
	var result models.ProjectMembers
	if err := stmt.Query(db.Postgres, &result); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return result, custom_errors.NoDataFound
		}
		return result, err
	}

	return result, nil
}

func _getProjectsByUserId(db *database.DB, userId uuid.UUID) ([]models.Projects, error) {
	stmt := table.Projects.SELECT(table.Projects.AllColumns).
		FROM(table.Projects.LEFT_JOIN(table.ProjectMembers, table.ProjectMembers.ProjectID.EQ(table.Projects.ID))).
		WHERE(postgres.AND(
			table.ProjectMembers.UserID.EQ(postgres.UUID(userId)),
			table.Projects.IsDeleted.EQ(postgres.Bool(false)),
		))

	logger.Logger.Log("models.project_members.GetProjectsByUserId", stmt.DebugSql())
	var result []models.Projects
	if err := stmt.Query(db.Postgres, &result); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return result, custom_errors.NoDataFound
		}
		return result, err
	}

	return result, nil
}

// This function is used to verify project members.
func _getProjectByIdAndUserId(db *database.DB, id uuid.UUID, userId uuid.UUID) (models.ProjectAccessDetails, error) {
	stmt := table.Projects.SELECT(table.Projects.ID, table.Projects.Name,
		table.Projects.OwnerID, table.ProjectMembers.Permissions).
		FROM(table.Projects.LEFT_JOIN(table.ProjectMembers, table.ProjectMembers.ProjectID.EQ(table.Projects.ID))).
		WHERE(postgres.AND(
			table.ProjectMembers.UserID.EQ(postgres.UUID(userId)),
			table.ProjectMembers.IsActive.EQ(postgres.Bool(true)),
			table.Projects.IsDeleted.EQ(postgres.Bool(false)),
			table.ProjectMembers.ProjectID.EQ(postgres.UUID(id)),
		))

	logger.Logger.Log("models.project_members.GetProjectByIdAndUserId", stmt.DebugSql())
	var result models.ProjectAccessDetails
	if err := stmt.Query(db.Postgres, &result); err != nil {
		logger.Logger.Error("models.project_members.GetProjectByIdAndUserId: ERROR", err)
		if errors.Is(err, qrm.ErrNoRows) {
			return result, custom_errors.NoDataFound
		}
		return result, err
	}
	return result, nil
}

/*
TODO - Add pagination
*/
func _getMembersByProjectId(db *database.DB, projectId uuid.UUID) ([]models.ProjectMembersWithUserInfo, error) {
	stmt := table.ProjectMembers.SELECT(table.ProjectMembers.AllColumns, table.Users.ID,
		table.Users.Email,
		table.Users.Username, table.Users.LastLoginAt).
		FROM(table.ProjectMembers.LEFT_JOIN(table.Users,
			table.ProjectMembers.UserID.EQ(table.Users.ID))).
		WHERE(postgres.AND(
			table.ProjectMembers.IsDeleted.EQ(postgres.Bool(false)),
			table.ProjectMembers.ProjectID.EQ(postgres.UUID(projectId)),
		)).ORDER_BY(table.ProjectMembers.CreatedAt.DESC())

	logger.Logger.Log("models.project_members.GetMembersByProjectId", stmt.DebugSql())
	var result []models.ProjectMembersWithUserInfo
	if err := stmt.Query(db.Postgres, &result); err != nil {
		logger.Logger.Error("models.project_members.GetMembersByProjectId: ERROR", err)
		if errors.Is(err, qrm.ErrNoRows) {
			return result, custom_errors.NoDataFound
		}
		return result, err
	}

	return result, nil
}

func _update(db *database.DB, id uuid.UUID, data models.ProjectMembers) error {
	columnsList := postgres.ColumnList{}

	if data.IsActive != nil && reflect.ValueOf(*data.IsActive).Kind() == reflect.Bool {
		columnsList = append(columnsList, table.ProjectMembers.IsActive)
	}
	if data.IsDeleted != nil && reflect.ValueOf(*data.IsDeleted).Kind() == reflect.Bool {
		columnsList = append(columnsList, table.ProjectMembers.IsDeleted)
	}

	if len(columnsList) <= 0 {
		logger.Logger.Log("models.project_members.Update: nothing to update")
		return nil
	}
	stmt := table.ProjectMembers.UPDATE(columnsList).
		MODEL(data).WHERE(table.ProjectMembers.ID.EQ(postgres.UUID(id)))

	logger.Logger.Log("models.project_members.Update: updating", stmt.DebugSql())
	_, err := stmt.Exec(db.Postgres)
	if err != nil {
		logger.Logger.Error("models.project_members.Update: ERROR", err)
		return custom_errors.DBErrors(err)
	}

	return nil
}
