package entity

//OAuthRedirectValues holds values required for Oauth
type OAuthRedirectValues struct {
	GithubOAuthURL string
	ClientID       string
	ClientSecret   string
	RedirectURL    string
	Login          string
	Scope          string
	State          string
	AllowSignUp    string
}

//AccessToken represents access token
type AccessToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}