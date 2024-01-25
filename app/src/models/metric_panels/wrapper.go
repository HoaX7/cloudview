package models

import (
	"cloudview/app/src/database"
	"cloudview/app/src/models"

	"github.com/google/uuid"
)

func Create(db *database.DB, data models.MetricPanels) (models.MetricPanels, error) {
	return _create(db, data)
}

func Update(db *database.DB, id uuid.UUID, data models.MetricPanels) error {
	return _update(db, id, data)
}

func GetById(db *database.DB, id uuid.UUID) (models.MetricPanels, error) {
	return _getById(db, id)
}

func GetByProviderAccount(db *database.DB, providerAccountId uuid.UUID) ([]models.MetricPanels, error) {
	return _getByProviderAccount(db, providerAccountId)
}
