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
	Permissions string           `json:"permissions,omitempty"`
	Metadata    *json.RawMessage `json:"metadata,omitempty"`
	IsDeleted   *bool            `json:"isDeleted"`
	IsSystem    *bool            `json:"isSystem"`
	CreatedAt   *time.Time       `json:"createdAt"`
	UpdatedAt   *time.Time       `json:"updatedAt"`
}
type ProjectMembersWithUserInfo struct {
	ProjectMembers
	User user_model.Users `json:"user"`
}

type ProviderAccountWithProject struct {
	ProviderAccounts
	Project Projects `json:"project"`
}

type Projects struct {
	ID          uuid.UUID        `sql:"primary_key" json:"id,omitempty"`
	Name        string           `json:"name,omitempty"`
	Description *string          `json:"description,omitempty"`
	Email       string           `db:"email,omitempty" json:"email,omitempty"`
	OwnerID     *uuid.UUID       `json:"ownerId,omitempty"`
	Members     *json.RawMessage `jet:"type=jsonb,nullable" json:"members,omitempty"`
	MemberLimit int32            `json:"memberLimit,omitempty"`
	Type        string           `json:"type,omitempty"`
	Metadata    *json.RawMessage `jet:"type=jsonb,nullable" json:"metadata,omitempty"`
	IsDeleted   *bool            `json:"isDeleted,omitempty"`
	CreatedAt   *time.Time       `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time       `json:"updatedAt,omitempty"`
}

type ProviderAccounts struct {
	ID                      uuid.UUID        `sql:"primary_key" json:"id,omitempty"`
	Name                    string           `json:"name,omitempty"`
	Description             string           `json:"description,omitempty"`
	AccessKeyID             string           `json:"accessKeyId,omitempty"`
	AccessKeySecret         string           `json:"accessKeySecret,omitempty"`
	RotationSecretKey       string           `json:"rotationSecretKey,omitempty"`
	Provider                string           `json:"provider,omitempty"`
	Type                    *string          `json:"type,omitempty"`
	AccountID               string           `json:"accountId,omitempty"`
	AccessRole              *string          `json:"accessRole,omitempty"`
	FeatureAccessPermission string           `json:"featureAccessPermission,omitempty"`
	ProjectID               *uuid.UUID       `json:"projectId,omitempty" jet:"nullable"`
	Metadata                *json.RawMessage `json:"metadata,omitempty" jet:"type=jsonb,nullable"`
	IsDeleted               *bool            `json:"isDeleted,omitempty"`
	CreatedAt               *time.Time       `json:"createdAt,omitempty"`
	UpdatedAt               *time.Time       `json:"updatedAt,omitempty"`
}

type MetricPanels struct {
	ID                uuid.UUID        `sql:"primary_key" json:"id,omitempty"`
	Name              string           `json:"name,omitempty"`
	Description       *string          `json:"description,omiempty"`
	Panels            *json.RawMessage `json:"panels,omitempty" jet:"type=jsonb,nullable"`
	ProviderAccountID uuid.UUID        `json:"providerAccountId,omitempty"`
	Metadata          *json.RawMessage `json:"metadata,omitempty" jet:"type=jsonb,nullable"`
	IsDeleted         *bool            `json:"isDeleted,omitempty"`
	CreatedAt         *time.Time       `json:"createdAt,omitempty"`
	UpdatedAt         *time.Time       `json:"updatedAt,omitempty"`
	InstanceID        string           `json:"instanceId,omitempty"`
	HealthStatus      string           `json:"healthStatus,omitempty"`
}

type ProjectAccessDetails struct {
	Projects
	ProjectMembers ProjectMembers `json:"projectMembers"`
}
