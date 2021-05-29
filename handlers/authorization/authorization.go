package authorization

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"

	"accurics-test/db"
)

//Redirector defines Redirect signature
type Redirector interface {
	Redirect(ctx context.Context) (string, error)
}

//AuthHandler handles authorization requests
type AuthHandler struct {
	uc Redirector
}

//CreateHandler creates authorization handler
func CreateHandler(uc Redirector) *AuthHandler {
	return &AuthHandler{
		uc: uc,
	}
}

func (ah *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user := params["user"]
	db.UserMapLock.RLock()
	_, ok := db.UserToToken[user]
	db.UserMapLock.RUnlock()
	if ok {
		w.Write([]byte("User Authenticated"))
		return
	}
	ctx := context.WithValue(context.Background(), "username", user)
	reDirectURL, err := ah.uc.Redirect(ctx)
	if err != nil {
		w.Write([]byte("Error in getting redirect URL" + err.Error()))
		return
	}
	http.Redirect(w, r, reDirectURL, http.StatusMovedPermanently)
}
