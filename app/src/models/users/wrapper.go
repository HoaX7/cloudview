package models

import "cloudview/app/src/database"

func GetByEmail(db *database.DB, email string) (Users, error) {
	return _getByEmail(db, email)
}
