//
// Responsible for handle queries related to users table.
// The model uses pg-jet to prepare dynamic statements to be executed.
//

package models

import (
	// Jet is a query statement builder only.

	"cloudview/app/dbschema/cloudview/public/table"
	custom_errors "cloudview/app/src/api/errors"
	"cloudview/app/src/api/middleware/logger"
	"cloudview/app/src/database"
	"errors"

	"encoding/json"
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	_ "github.com/lib/pq"

	"github.com/google/uuid"
)

type Users struct {
	ID                   uuid.UUID        `json:"id,omitempty" sql:"primary_key"`
	Username             string           `json:"username,omitempty"`
	Email                string           `json:"email,omitempty"`
	AvatarURL            *string          `json:"avatarUrl,omitempty"`
	SubscribedSince      *time.Time       `json:"subscribedSince,omitempty" jet:"nullable"`
	SubscriptionDaysLeft *int32           `json:"subscriptionDaysLeft,omitempty"`
	LastLoginAt          *time.Time       `json:"lastLoginAt,omitempty"`
	SubscriptionPlanID   *int32           `json:"subscriptionPlanId,omitempty"`
	Metadata             *json.RawMessage `jet:"type=jsonb,nullable" json:"metadata,omitempty"`
	IsDeleted            *bool            `json:"isDeleted,omitempty"`
	CreatedAt            *time.Time       `json:"createdAt,omitempty"`
	UpdatedAt            *time.Time       `json:"updatedAt,omitempty"`
	Permissions          *string          `json:"permissions,omitempty"`
}

func Create(db *database.DB, data Users) (Users, error) {
	//
	// Jet query statement builder returns the query & args.
	// Need to pass both for query to execute properly.
	//
	stmt := table.Users.INSERT(table.Users.Username,
		table.Users.Email, table.Users.AvatarURL,
		table.Users.Metadata, table.Users.Permissions).
		MODEL(data).
		RETURNING(table.Users.AllColumns)

	queryString, args := stmt.Sql()
	logger.Logger.Log("Inserting into users table with data: ", queryString, args)

	result := Users{}
	// Currently using raw query since go-jet does not return pq.Error
	// if err := stmt.Query(db.Postgres, &result); err != nil {
	// 	logger.Logger.Log("error", err)
	// 	return result, custom_errors.DBErrors(err)
	// }
	rows, err := db.Postgres.Query(queryString, args...)
	if err != nil {
		return result, custom_errors.DBErrors(err)
	}
	/*
		While running raw query, need to manually enter
		all required columns to be added to struct.

		WARNING: THIS METHODOLOGY IS STRICTLY NOT RECOMMENDED TO USE.
		IT IS ONLY BEING USED TO RETURN DATA AFTER INSERT TO CATCH POSTGRES
		ERROR CODES.
	*/
	if rows.Next() {
		if err := rows.Scan(&result.ID, &result.Username,
			&result.Email,
			&result.AvatarURL, &result.SubscribedSince,
			&result.SubscriptionDaysLeft, &result.LastLoginAt,
			&result.Metadata, &result.IsDeleted,
			&result.CreatedAt,
			&result.UpdatedAt,
			&result.SubscriptionPlanID,
			&result.Permissions); err != nil {
			logger.Logger.Error("model.users.Create: ERROR", err)
			return result, custom_errors.UnknownError
		}
	}
	return result, nil
}

func _getByEmail(db *database.DB, email string) (Users, error) {
	stmt := table.Users.SELECT(table.Users.AllColumns).
		WHERE(postgres.AND(
			table.Users.Email.EQ(postgres.String(email)),
			table.Users.IsDeleted.EQ(postgres.Bool(false)),
		))

	result := Users{}
	if err := stmt.Query(db.Postgres, &result); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return result, custom_errors.NoDataFound
		}
		return result, err
	}

	return result, nil
}
