package types

type OauthProviders struct {
	GOOGLE string
	GITHUB string
}

// Common struct to be returned from both github/google oauth login
type ProviderUser struct {
	Email       string
	Username    string
	AvatarUrl   string
	AccessToken string
}

type OauthLogin interface {
	Login() (*ProviderUser, error)
	GetUserData(access_token string) (*ProviderUser, error)
	Name() string
}

// This struct is to capture the data returned from oauth `getUserData` func
type GithubUser struct {
	AvatarUrl string `json:"avatar_url"`
	Email     string `json:"email"`
	Login     string `json:"login"` // has username
	// check documentation for all available props
}

type GithubUserEmails struct {
	Email     string
	Primary   bool
	Verified  bool
	Visiblity string
}
