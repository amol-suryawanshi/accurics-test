package handlers

import (
	"accurics-test/db"
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

//GithubAction defines interface for github action
type GithubAction interface {
	Perform(ctx context.Context, user, repo, branchName, token string) error
}

//CreateGithubHandler returns instance of GithubHandler
func CreateGithubHandler(uc GithubAction) *GithubHandler {
	return &GithubHandler{
		uc: uc,
	}
}

//GithubHandler handles github request
type GithubHandler struct {
	uc GithubAction
}

func (gh *GithubHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user := params["user"]
	if user == "" {
		w.Write([]byte("Kindly provide user name"))
		return
	}
	db.UserMapLock.RLock()
	token, ok := db.UserToToken[user]
	db.UserMapLock.RUnlock()
	if !ok {
		w.Write([]byte("Kindly authenticate first!!"))
		return
	}

	repo := params["repo"]
	if repo == "" {
		w.Write([]byte("Kindly provide repository name"))
		return
	}

	branchName := params["branch"]
	if branchName == "" {
		w.Write([]byte("Kindly provide branch name"))
		return
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "fileName", "testfile.txt")

	err := gh.uc.Perform(ctx, user, repo, branchName, token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte("Pull Request created"))
	return

}
