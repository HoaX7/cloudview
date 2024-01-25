package models

import (
	"cloudview/app/dbschema/cloudview/public/table"
	"cloudview/app/src/api/encryption"
	custom_errors "cloudview/app/src/api/errors"
	"cloudview/app/src/api/middleware/logger"
	"cloudview/app/src/database"
	"cloudview/app/src/models"
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
)

type Visiblity string // Enum ("PUBLIC", "PRIVATE")

/*
Fetch projects where owner_id matches or if the user id is present
in the member list.
*/
func _getByOwnerId(db *database.DB, userId uuid.UUID) ([]models.Projects, error) {
	stmt := table.Projects.SELECT(table.Projects.AllColumns).
		WHERE(postgres.AND(
			table.Projects.OwnerID.EQ(postgres.UUID(userId)),
			table.Projects.IsDeleted.EQ(postgres.Bool(false)),
		)).ORDER_BY(table.Projects.UpdatedAt.DESC())

	var result []models.Projects
	if err := stmt.Query(db.Postgres, &result); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return result, custom_errors.NoDataFound
		}
		return result, err
	}

	return result, nil
}

func _getById(db *database.DB, id uuid.UUID) (models.Projects, error) {
	// stmt := table.Projects.SELECT(table.Projects.AllColumns).
	// 	WHERE(postgres.AND(
	// 		table.Projects.ID.EQ(postgres.UUID(id)),
	// 		table.Projects.IsDeleted.EQ(postgres.Bool(false)),
	// 	))

	// var result models.Projects
	// if err := stmt.Query(db.Postgres, &result); err != nil {
	// 	if errors.Is(err, qrm.ErrNoRows) {
	// 		return result, custom_errors.NoDataFound
	// 	}
	// 	return result, err
	// }

	// return result, nil
	result, err := GetByIds(db, []uuid.UUID{id})
	return result[0], err
}

func _getByIdAndUserId(db *database.DB, id uuid.UUID, userId uuid.UUID) (models.Projects, error) {
	stmt := table.Projects.SELECT(table.Projects.AllColumns).
		WHERE(postgres.AND(
			table.Projects.ID.EQ(postgres.UUID(id)),
			table.Projects.IsDeleted.EQ(postgres.Bool(false)),
			table.Projects.OwnerID.EQ(postgres.UUID(userId)),
		))

	logger.Logger.Log("models.projects.GetByIdAndUserId", stmt.DebugSql())
	var result models.Projects
	if err := stmt.Query(db.Postgres, &result); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return result, custom_errors.NoDataFound
		}
		return result, err
	}

	return result, nil
}

type CreateWithServiceProps struct {
	Name            string    `json:"name"`
	Description     string    `json:"description,omitempty" optional:"true"`
	AccessKeyID     string    `json:"accessKeyId"`
	AccessKeySecret string    `json:"accessKeySecret"`
	Provider        string    `json:"provider"`
	Type            string    `json:"type"`
	OwnerID         uuid.UUID `json:"ownerId"`
	Email           string    `json:"email"`
}
type ServiceResult struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	ProjectID uuid.UUID `json:"projectId"`
	Provider  string    `json:"provider"`
}
type CreateWithServiceReturnType struct {
	Project models.Projects `json:"project"`
	Service ServiceResult   `json:"service"`
}

/*
@Deprecated - In favor of using cross account access
*/
func CreateWithService(db *database.DB, data CreateWithServiceProps) (CreateWithServiceReturnType, error) {
	var result CreateWithServiceReturnType
	// Use DB `transactions` to insert projects, project_members
	// and access keys.
	fail := func(err error) error {
		return fmt.Errorf("[models.projects.createWithService] Error running transaction: %v", err)
	}
	ctx := context.TODO()
	tx, err := db.RawDB.BeginTx(ctx, nil)
	if err != nil {
		logger.Logger.Error(fail(err).Error())
		return result, errors.New("Unable to create Project, Please try again later")
	}
	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	projectData := models.Projects{
		OwnerID:     &data.OwnerID,
		Name:        data.Name,
		Description: &data.Description,
		Email:       data.Email,
	}
	// Prepare insert statements and execute queries
	projectStmt := table.Projects.INSERT(table.Projects.Name, table.Projects.OwnerID,
		table.Projects.Description).MODEL(projectData).RETURNING(table.Projects.AllColumns)
	projectSql, args := projectStmt.Sql()
	projectResult := models.Projects{}
	if err := tx.QueryRowContext(ctx, projectSql, args...).Scan(&projectResult.ID, &projectResult.Name,
		&projectResult.Description,
		&projectResult.Email, &projectResult.OwnerID,
		&projectResult.Members, &projectResult.MemberLimit,
		&projectResult.Type, &projectResult.Metadata,
		&projectResult.IsDeleted,
		&projectResult.CreatedAt,
		&projectResult.UpdatedAt); err != nil {
		logger.Logger.Error("models.projects.CreateWithService: ERROR creating project", err.Error())
		return result, errors.New("Unable to insert data")
	}

	// prepare project_member insert statements
	projectMemberStmt := table.ProjectMembers.INSERT(table.ProjectMembers.ProjectID, table.ProjectMembers.UserID,
		table.ProjectMembers.IsOwner).MODEL(models.ProjectMembers{
		ProjectID: projectResult.ID,
		UserID:    *projectResult.OwnerID,
		IsOwner:   true,
	})
	stmtSql, stmtArgs := projectMemberStmt.Sql()
	if _, err := tx.ExecContext(ctx, stmtSql, stmtArgs...); err != nil {
		logger.Logger.Error("models.projects.CreateWithService: ERROR inserting into project_members", err.Error())
		return result, errors.New("Unable to insert data")
	}

	key, err := encryption.GenerateRandomSecretKey(16)
	if err != nil {
		logger.Logger.Error("models.projects.CreateWithService: ERROR unable to generate secret key", err)
		return result, errors.New("Unable to insert data")
	}
	cipherText, err := encryption.Encrypt(data.AccessKeySecret, key)
	if err != nil {
		logger.Logger.Error("models.projects.CreateWithService: ERROR unable to encrypt access key", err)
		return result, errors.New("Unable to insert data")
	}
	serviceData := models.ProviderAccounts{
		AccessKeySecret:   cipherText,
		AccessKeyID:       data.AccessKeyID,
		Provider:          data.Provider,
		ProjectID:         &projectResult.ID,
		RotationSecretKey: key,
		Name:              data.Provider,
	}
	serviceStmt := table.ProviderAccounts.INSERT(table.ProviderAccounts.Name, table.ProviderAccounts.AccessKeyID,
		table.ProviderAccounts.AccessKeySecret, table.ProviderAccounts.RotationSecretKey,
		table.ProviderAccounts.ProjectID, table.ProviderAccounts.Provider).MODEL(serviceData).
		RETURNING(table.ProviderAccounts.ID, table.ProviderAccounts.Name, table.ProviderAccounts.ProjectID, table.ProviderAccounts.Provider)

	serviceSql, sargs := serviceStmt.Sql()
	serviceResult := ServiceResult{}
	if err := tx.QueryRowContext(ctx, serviceSql, sargs...).Scan(&serviceResult.ID,
		&serviceResult.Name,
		&serviceResult.ProjectID,
		&serviceResult.Provider); err != nil {
		logger.Logger.Error("models.projects.CreateWithService: ERROR creating service", err.Error())
		return result, errors.New("Unable to insert data")
	}
	if err := tx.Commit(); err != nil {
		logger.Logger.Error("models.projects.CreateWithService: Transaction Commit Failed", fail(err).Error())
		return result, errors.New("Unable to insert data")
	}
	logger.Logger.Log("models.projects.CreateWithService: Transaction successfully committed")
	return CreateWithServiceReturnType{
		Project: projectResult,
		Service: serviceResult,
	}, nil
}

func _update(db *database.DB, id uuid.UUID, ownerId uuid.UUID, data models.Projects) error {
	columnsList := postgres.ColumnList{}
	if data.Name != "" {
		columnsList = append(columnsList, table.Projects.Name)
	}
	if data.Description != nil && reflect.ValueOf(*data.Description).Kind() == reflect.String {
		columnsList = append(columnsList, table.Projects.Description)
	}
	if data.Email != "" {
		columnsList = append(columnsList, table.Projects.Email)
	}
	if data.IsDeleted != nil && reflect.ValueOf(*data.IsDeleted).Kind() == reflect.Bool {
		columnsList = append(columnsList, table.Projects.IsDeleted)
	}
	if len(columnsList) <= 0 {
		logger.Logger.Log("models.projects.Update: nothing to update")
		return nil
	}
	stmt := table.Projects.UPDATE(columnsList).
		MODEL(data).WHERE(postgres.AND(
		table.Projects.ID.EQ(postgres.UUID(id)),
		table.Projects.OwnerID.EQ(postgres.UUID(ownerId)),
	))

	logger.Logger.Log("models.projects.Update: updating", stmt.DebugSql())
	_, err := stmt.Exec(db.Postgres)
	if err != nil {
		logger.Logger.Error("models.projects.Update: ERROR", err)
		return custom_errors.DBErrors(err)
	}
	return nil
}

func Create(db *database.DB, data models.Projects) (models.Projects, error) {
	stmt := table.Projects.INSERT(table.Projects.Name,
		table.Projects.Description,
		table.Projects.Email,
		table.Projects.OwnerID,
		table.Projects.Type,
		table.Projects.Members,
		table.Projects.MemberLimit,
		table.Projects.Metadata).
		MODEL(data).
		RETURNING(table.Projects.AllColumns)

	queryString, args := stmt.Sql()
	logger.Logger.Log("Inserting into Projects table with data: ", queryString, args)

	result := models.Projects{}
	rows, err := db.Postgres.Query(queryString, args...)
	if err != nil {
		logger.Logger.Error("ProjectsModel.Create: ERROR", err)
		return result, custom_errors.DBErrors(err)
	}
	if rows.Next() {
		if err := rows.Scan(&result.ID, &result.Name,
			&result.Description,
			&result.Email, &result.OwnerID,
			&result.Members, &result.MemberLimit,
			&result.Type, &result.Metadata,
			&result.IsDeleted,
			&result.CreatedAt,
			&result.UpdatedAt); err != nil {

			logger.Logger.Error("models.Projects.Create: ERROR", err)
			return result, err
		}
	}
	return result, nil
}

func GetByIds(db *database.DB, ids []uuid.UUID) ([]models.Projects, error) {
	params := []postgres.Expression{}
	for _, uid := range ids {
		params = append(params, postgres.UUID(uid))
	}
	stmt := table.Projects.SELECT(table.Projects.AllColumns).
		WHERE(postgres.AND(
			table.Projects.ID.IN(params...),
			table.Projects.IsDeleted.EQ(postgres.Bool(false)),
		))

	var result []models.Projects
	if err := stmt.Query(db.Postgres, &result); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return result, custom_errors.NoDataFound
		}
		return result, err
	}

	return result, nil
}
