package usecases

import (
	"accurics-test/entity"
	"context"
	"fmt"
)

const (
	redirecQueryParam = "%s/authorize?client_id=%s&redirect_uri=%s&login=%s&scope=%s&state=%s&allow_signup=%s"
)

//CreateAuthRedirect creates AuthRedirect instance
func CreateAuthRedirect(reader Reader) *AuthRedirect {
	return &AuthRedirect{
		reader: reader,
	}
}

//Reader to read oauth related values from environment and config
type Reader interface {
	Get(ctx context.Context) (entity.OAuthRedirectValues, error)
}

//AuthRedirect generates redirection url
type AuthRedirect struct {
	reader Reader
}

//Redirect creates a redirection url
func (ar *AuthRedirect) Redirect(ctx context.Context) (string, error) {
	oAuthValues, err := ar.reader.Get(ctx)
	if err != nil {
		return "", err
	}
	return createRedirectURL(oAuthValues), nil
}

func createRedirectURL(redirectValues entity.OAuthRedirectValues) string {

	return fmt.Sprintf(redirecQueryParam, redirectValues.GithubOAuthURL, redirectValues.ClientID,
		redirectValues.RedirectURL, redirectValues.Login, redirectValues.Scope, redirectValues.State,
		redirectValues.AllowSignUp)
}
