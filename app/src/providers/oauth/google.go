package oauth

import (
	"cloudview/app/src/types"
	"errors"
)

type Google struct{ Code string }

func (g Google) Login() (*types.ProviderUser, error) {
	return nil, nil
}

func (g Google) GetUserData(access_token string) (*types.ProviderUser, error) {
	return nil, errors.New("Unimplemented")
}

func (g Google) Name() string {
	return "google"
}
