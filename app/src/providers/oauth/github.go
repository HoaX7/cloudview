package oauth

import (
	// "cloudview/app/src/helpers"
	// "cloudview/app/src/helpers/requester"
	"cloudview/app/src/api/middleware/logger"
	"cloudview/app/src/helpers"
	"cloudview/app/src/helpers/requester"
	"cloudview/app/src/types"
	"errors"

	"golang.org/x/exp/slices"
)

type Github struct {
	Code string
}

const accessTokenUrl = "https://github.com/login/oauth/access_token"
const fetchUserDataUrl = "https://api.github.com/user"
const fetchUserEmailsUrl = "https://api.github.com/user/emails"

func (g Github) Login() (*types.ProviderUser, error) {
	logger.Logger.Log("providers.oauth.github.Login: logging in via github")
	clientId := helpers.GoDotEnvVariable("GITHUB_OAUTH_CLIENT_ID")
	clientSecret := helpers.GoDotEnvVariable("GITHUB_OAUTH_CLIENT_SECRET")
	if clientId == "" || clientSecret == "" {
		return nil, errors.New("Missing github environment variables & credentials")
	}
	accessToken, err := g.getAccessToken(clientId, clientSecret)
	if err != nil {
		return nil, err
	}
	user, err := g.GetUserData(accessToken)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (g Github) GetUserData(accessToken string) (*types.ProviderUser, error) {
	var user types.GithubUser
	headers := map[string]string{
		"Authorization": "Bearer " + accessToken,
	}
	logger.Logger.Log("github.getUserData: fetch user_data starting")
	if err := requester.MakeRequest(fetchUserDataUrl, "", "GET", &user, requester.WithHeaders(headers)); err != nil {
		return nil, err
	}
	/*
		Need to check why the func below returns empty object
		instead of converted json
	*/
	// if err := utility.ConvertMapToStruct(target, &user); err != nil {
	// 	logger.Logger.Log("github.getUserData: error parsing user data", err)
	// 	return "", errors.New("Unable to retrieve user data")
	// }
	logger.Logger.Log("github.getUserData: fetch user_data success", user)
	logger.Logger.Log("access_token", accessToken)

	/*
		If your Github profile email is set to private it will not be available
		in response object. The email will be available if it is set to public.
		Will need to make another api call using the same access_token to fetch
		email.
	*/
	if user.Email == "" {
		logger.Logger.Log("github.getUserData: public email not found, fetching private email")
		email, err := g.getUserEmails(accessToken)
		if err != nil {
			return nil, err
		}
		user.Email = email
	}

	var result types.ProviderUser
	result.Email = user.Email
	result.Username = user.Login
	result.AvatarUrl = user.AvatarUrl
	result.AccessToken = accessToken
	logger.Logger.Log("github.getUserData: user authenticated:", result)
	return &result, nil
}

func (g Github) getAccessToken(clientId string, clientSecret string) (string, error) {
	queryParams := "?client_id=" + clientId + "&client_secret=" + clientSecret + "&code=" + g.Code
	headers := map[string]string{
		"Accept": "application/json",
	}

	logger.Logger.Log("github.getAccessToken: fetch access_token starting")
	// Data returned by github api is map[string]
	var target map[string]interface{}
	if err := requester.MakeRequest(accessTokenUrl+queryParams, "", "POST", &target, requester.WithHeaders(headers)); err != nil {
		return "", err
	}
	/*
		Manually handling error from github api.
		This is required as the api does not return errors but 200 status code
		even if there are errors with the code.

		The issue is currently active on github
		https://github.com/orgs/community/discussions/57068
	*/
	if _, err := target["error"]; err {
		description := target["error_description"].(string)
		return "", errors.New(description)
	}

	// Bearer token, use in Authorization headers to retrieve user data
	accessToken := target["access_token"].(string)
	logger.Logger.Log("github.getAccessToken: success", target)
	return accessToken, nil
}

/*
Github returns a slice of structs.
We only need the primary email (github always has atleast 1 primary email)
*/
func (g Github) getUserEmails(accessToken string) (string, error) {
	var target []types.GithubUserEmails
	headers := map[string]string{
		"Authorization": "Bearer " + accessToken,
	}
	logger.Logger.Log("github.getUserEmails: fetch email starting")
	if err := requester.MakeRequest(fetchUserEmailsUrl, "", "GET", &target, requester.WithHeaders(headers)); err != nil {
		logger.Logger.Log("github.getUserEmails: error fetching email", err)
		return "", errors.New("Unable to fetch email")
	}
	idx := slices.IndexFunc(target, func(e types.GithubUserEmails) bool {
		return e.Primary
	})
	if idx == -1 {
		return "", errors.New("Unable to fetch primary email. Make sure you have given proper access.")
	}
	return target[idx].Email, nil
}

func (g Github) Name() string {
	return "github"
}
