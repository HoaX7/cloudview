package models

import (
	"cloudview/app/src/database"
	"cloudview/app/src/models"

	"github.com/google/uuid"
)

func Create(db *database.DB, data models.ProjectMembers) (models.ProjectMembers, error) {
	return _create(db, data)
}

func Update(db *database.DB, id uuid.UUID, data models.ProjectMembers) error {
	return _update(db, id, data)
}

func GetMembersByProjectId(db *database.DB, projectId uuid.UUID) ([]models.ProjectMembersWithUserInfo, error) {
	return _getMembersByProjectId(db, projectId)
}

func GetProjectByIdAndUserId(db *database.DB, id uuid.UUID, userId uuid.UUID) (models.Projects, error) {
	return _getProjectByIdAndUserId(db, id, userId)
}

func GetProjectsByUserId(db *database.DB, userId uuid.UUID) ([]models.Projects, error) {
	return _getProjectsByUserId(db, userId)
}

func GetByIdAndUserId(db *database.DB, id uuid.UUID, userId uuid.UUID) (models.ProjectMembers, error) {
	return _getByIdAndUserId(db, id, userId)
}
