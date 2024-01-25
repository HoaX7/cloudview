package types

import "cloudview/app/src/models"

type VerifyProjectAccessInput struct {
	ProjectID         interface{}
	ProviderAccountID interface{}
}

type VerifyProjectAccessOutput struct {
	ProjectAccessDetails models.ProjectAccessDetails
	ProviderAccount      *models.ProviderAccounts `json:"providerAccount"`
}
