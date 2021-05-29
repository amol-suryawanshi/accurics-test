package muxrouter

import (
	"accurics-test/handlers"
	"accurics-test/handlers/authorization"
	"accurics-test/repositories"
	"accurics-test/usecases"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	authHandler     = "auth"
	callbackHandler = "callback"
	createHandler   = "create"
)

//GetMuxRouter returns instance of mux router
func GetMuxRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.Handle("/api/v1/{user}/authorize", getHandler(authHandler)).Methods(http.MethodGet)
	router.Handle("/api/v1/callback", getHandler(callbackHandler)).Methods(http.MethodGet)
	router.Handle("/api/v1/create/{user}/{repo}/{branch}", getHandler(createHandler)).Methods(http.MethodGet)

	return router
}

func getHandler(key string) http.Handler {
	switch key {

	case authHandler:
		repo := repositories.CreateOAuthReader()
		uc := usecases.CreateAuthRedirect(repo)
		return authorization.CreateHandler(uc)

	case callbackHandler:
		accessTokenRepo := repositories.CreateGithubTokenImpl()
		oAuthRepo := repositories.CreateOAuthClientSecretReader()
		uc := usecases.CreateGithubToken(oAuthRepo, accessTokenRepo)
		return authorization.CreateAccessTokenHandler(uc)

	case createHandler:
		repo := repositories.CreateGithubImpl()
		uc := usecases.CreateCombinedUseCase(repo)
		return handlers.CreateGithubHandler(uc)

	}
	return nil
}
