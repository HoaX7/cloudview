package types

import (
	users_model "cloudview/app/src/models/users"
)

type SessionUser struct {
	*users_model.Users
	AccessToken string `json:"accessToken"`
	Provider    string `json:"provider"`
}
