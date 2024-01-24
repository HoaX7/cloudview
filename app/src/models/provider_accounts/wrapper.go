package models

import (
	"cloudview/app/src/cache"
	cache_keys "cloudview/app/src/cache/keys"
	"cloudview/app/src/database"
	"cloudview/app/src/models"
	"fmt"

	"github.com/google/uuid"
)

func GetByIdForSDK(db *database.DB, id uuid.UUID) (models.ProviderAccounts, error) {
	return _getByIdForSDK(db, id)
}

func GetById(db *database.DB, id uuid.UUID) (interface{}, error) {
	key := fmt.Sprintf("%s_%s", cache_keys.PROVIDER_ACCOUNT, id)
	result, err := cache.Fetch(key, 0, func() (interface{}, error) {
		return _getById(db, id)
	})
	return result, err
}

func Update(db *database.DB, id uuid.UUID, data models.ProviderAccounts) error {
	return _update(db, id, data)
}
