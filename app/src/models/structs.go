package models

import (
	user_model "cloudview/app/src/models/users"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ProjectMembers struct {
	ID          uuid.UUID        `sql:"primary_key" json:"id"`
	ProjectID   uuid.UUID        `json:"projectId"`
	UserID      uuid.UUID        `json:"userId"`
	IsOwner     bool             `json:"isOwner"`
	IsActive    *bool            `json:"isActive"`
	Permissions *json.RawMessage `json:"permissions,omitempty"`
	Metadata    *json.RawMessage `json:"metadata,omitempty"`
	IsDeleted   *bool            `json:"isDeleted"`
	CreatedAt   *time.Time       `json:"createdAt"`
	UpdatedAt   *time.Time       `json:"updatedAt"`
}
type ProjectMembersWithUserInfo struct {
	ProjectMembers
	User user_model.Users `json:"user"`
}

type Projects struct {
	ID          uuid.UUID        `sql:"primary_key" json:"id,omitempty"`
	Name        string           `json:"name,omitempty"`
	Description *string          `json:"description,omitempty"`
	Email       string           `db:"email,omitempty" json:"email,omitempty"`
	OwnerID     uuid.UUID        `json:"ownerId,omitempty"`
	Members     *json.RawMessage `jet:"type=jsonb,nullable" json:"members,omitempty"`
	MemberLimit int32            `json:"memberLimit,omitempty"`
	Type        string           `json:"type,omitempty"`
	Metadata    *json.RawMessage `jet:"type=jsonb,nullable" json:"metadata,omitempty"`
	IsDeleted   *bool            `json:"isDeleted,omitempty"`
	CreatedAt   time.Time        `json:"createdAt,omitempty"`
	UpdatedAt   time.Time        `json:"updatedAt,omitempty"`
}
