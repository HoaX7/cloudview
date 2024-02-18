package oauth

import (
	custom_errors "cloudview/app/src/api/errors"
	"cloudview/app/src/api/middleware/logger"
	"cloudview/app/src/database"
	users_model "cloudview/app/src/models/users"
	"cloudview/app/src/permissions"
	"cloudview/app/src/types"
	"errors"
)

func Login(provider types.OauthLogin, db *database.DB) (types.SessionUser, error) {
	data, err := provider.Login()
	var result types.SessionUser
	if err != nil {
		return result, err
	}

	// Fetch or create user from DB
	var user users_model.Users
	logger.Logger.Log("oauth.Login: fetching user data:", data.Email)
	existingUser, err := users_model.GetByEmail(db, data.Email)
	if err != nil {
		if errors.Is(err, custom_errors.NoDataFound) {
			perms := permissions.SetPermissions([]string{
				permissions.USER_MODIFY_PROJECT,
				permissions.USER_MODIFY_PROVIDER_ACCOUNT,
			})
			logger.Logger.Log("oauth.Login: user not found, creating new user...")
			data, err := users_model.Create(db, users_model.Users{
				Email:       data.Email,
				Username:    data.Username,
				AvatarURL:   &data.AvatarUrl,
				Permissions: &perms,
			})
			if err != nil {
				return result, err
			}
			user = data
		} else {
			return result, err
		}
	} else {
		user = existingUser
	}
	result.Users = &user
	result.Provider = provider.Name()
	result.AccessToken = data.AccessToken

	return result, nil
}
