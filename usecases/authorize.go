package usecases

import (
	"context"
	"errors"

	"accurics-test/entity"
)

//CreateGithubToken creates instance of AuthAccessToken
func CreateGithubToken(reader Reader, githubToken GithubToken) *AuthAccessToken {
	return &AuthAccessToken{
		reader: reader,
		gt:     githubToken,
	}
}

//GithubToken defines interface to get AccessToken
type GithubToken interface {
	AccessToken(context.Context, entity.OAuthRedirectValues) error
}

//AuthAccessToken implements GithubAction interface in handler package
type AuthAccessToken struct {
	reader Reader
	gt     GithubToken
}

//Authorize reads env vars and places call to get access token
func (aat *AuthAccessToken) Authorize(ctx context.Context) error {
	oAuthValues, err := aat.reader.Get(ctx)
	if err != nil {
		return err
	}
	code := ctx.Value("code").(string)
	if code == "" {
		// config.AppLogger.ErrorLogger.Println("auth code not found in context or blank auth code set")
		return errors.New("Code not found in context")
	}
	return aat.gt.AccessToken(ctx, oAuthValues)

}
