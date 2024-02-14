package oauth

/*

This provider is being used to allow users to see a
sample / demo of the product.

username: demo@gmail.com
password: demo123

*/

import (
	"cloudview/app/src/types"
	"errors"
)

type Local struct {
	Code     string
	Username string
	Password string
}

func (l Local) Login() (*types.ProviderUser, error) {
	if l.Username != "demo@gmail.com" || l.Password != "demo123" {
		return nil, errors.New("Invalid `username` or `password` provided.")
	}
	return l.GetUserData("dummy-token")
}

func (l Local) GetUserData(access_token string) (*types.ProviderUser, error) {
	return &types.ProviderUser{
		Email:       "demo@gmail.com",
		AccessToken: access_token,
		Username:    "demo",
	}, nil
}

func (l Local) Name() string {
	return "Local"
}
