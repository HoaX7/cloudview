//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import (
	"github.com/google/uuid"
	"time"
)

type ProviderAccounts struct {
	ID                      uuid.UUID `sql:"primary_key"`
	Name                    string
	Description             *string
	AccessKeyID             *string
	AccessKeySecret         *string
	RotationSecretKey       *string
	Provider                *string
	ProjectID               uuid.UUID
	Metadata                *string
	IsDeleted               *bool
	CreatedAt               *time.Time
	UpdatedAt               *time.Time
	DeletedAt               *time.Time
	Type                    *AccessType
	AccountID               string
	AccessRole              *string
	FeatureAccessPermission string
}
