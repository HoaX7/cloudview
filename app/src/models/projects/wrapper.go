package models

import (
	"cloudview/app/src/cache"
	cache_keys "cloudview/app/src/cache/keys"
	"cloudview/app/src/database"
	"cloudview/app/src/models"
	"fmt"

	"github.com/google/uuid"
)

func GetByOwnerId(db *database.DB, ownerId uuid.UUID) ([]models.Projects, error) {
	return _getByOwnerId(db, ownerId)
}

func GetById(db *database.DB, id uuid.UUID) (models.Projects, error) {
	key := fmt.Sprintf("%s_%s", cache_keys.PROJECT, id)
	var result models.Projects
	err := cache.Fetch(key, 0, &result, func() (interface{}, error) {
		return _getById(db, id)
	})
	return result, err
}

func GetByIdAndUserId(db *database.DB, id uuid.UUID, userId uuid.UUID) (models.Projects, error) {
	return _getByIdAndUserId(db, id, userId)
}

func Update(db *database.DB, id uuid.UUID, ownerId uuid.UUID, data models.Projects) error {
	key := fmt.Sprintf("%s_%s", cache_keys.PROJECT, id)
	cache.Del(key)
	return _update(db, id, ownerId, data)
}
