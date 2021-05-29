package repositories

import (
	"accurics-test/entity"
	"context"
	"errors"

	"accurics-test/config"
	"accurics-test/db"

	"github.com/google/uuid"
)

//CreateOAuthReader creates OAuthReader
func CreateOAuthReader() *OAuthReader {
	return &OAuthReader{}
}

//OAuthReader implementation for Get()
type OAuthReader struct{}

//Get value from env and config
func (oar *OAuthReader) Get(ctx context.Context) (entity.OAuthRedirectValues, error) {
	OAuthValues := entity.OAuthRedirectValues{}
	state, err := getState()
	if err != nil {
		return OAuthValues, err
	}
	OAuthValues.ClientID = config.ServerConfig.ClientID
	OAuthValues.GithubOAuthURL = config.ServerConfig.GithubOAuthURL
	OAuthValues.RedirectURL = config.ServerConfig.OAuthRedirectURL
	userID, ok := ctx.Value("username").(string)
	if !ok {
		return OAuthValues, errors.New("context doesn't contain userName")
	}
	OAuthValues.Login = userID
	OAuthValues.Scope = config.ServerConfig.Scope
	OAuthValues.State = state
	OAuthValues.AllowSignUp = config.ServerConfig.AllowSignUp
	db.StateMapLock.Lock()
	db.StateToUser[state] = userID
	db.StateMapLock.Unlock()

	return OAuthValues, nil
}

func getState() (string, error) {
	uid, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	return uid.String(), nil
}

//CreateOAuthClientSecretReader creates instance for OAuthClientSecretReader
func CreateOAuthClientSecretReader() *OAuthClientSecretReader {
	return &OAuthClientSecretReader{}
}

//OAuthClientSecretReader will read client secret and does not need state
type OAuthClientSecretReader struct{}

//Get value from config
func (oacsr *OAuthClientSecretReader) Get(ctx context.Context) (entity.OAuthRedirectValues, error) {
	OAuthValues := entity.OAuthRedirectValues{}
	OAuthValues.GithubOAuthURL = config.ServerConfig.GithubOAuthURL
	OAuthValues.ClientID = config.ServerConfig.ClientID
	OAuthValues.ClientSecret = config.ServerConfig.ClientSecret
	OAuthValues.RedirectURL = config.ServerConfig.OAuthRedirectURL
	var ok bool
	OAuthValues.State, ok = ctx.Value("state").(string)
	if !ok {
		return OAuthValues, errors.New("state not present in context")
	}

	return OAuthValues, nil
}
