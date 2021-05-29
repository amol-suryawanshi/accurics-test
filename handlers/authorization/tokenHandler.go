package authorization

import (
	"context"
	"net/http"

	"accurics-test/db"
)

//Authorizer defines Authorize signature
type Authorizer interface {
	Authorize(ctx context.Context) error
}

func CreateAccessTokenHandler(uc Authorizer) *AccessTokenHandler {
	return &AccessTokenHandler{
		uc: uc,
	}
}

//AccessTokenHandler handler function for /callback
type AccessTokenHandler struct {
	uc Authorizer
}

func (ath *AccessTokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// config.AppLogger.InfoLogger.Println("recieved /github/callback call")
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	db.StateMapLock.RLock()
	userID, ok := db.StateToUser[state]
	db.StateMapLock.RUnlock()
	if !ok {
		// config.AppLogger.ErrorLogger.Println("recieved call for unknown state:", state)
		w.Write([]byte("Kindly authenticate first!"))
		return
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "code", code)
	ctx = context.WithValue(ctx, "state", state)
	ctx = context.WithValue(ctx, "username", userID)

	err := ath.uc.Authorize(ctx)
	if err != nil {
		w.Write([]byte("Error occured during authorizing :" + err.Error()))
		return
	}
	w.Write([]byte("Authorization Success"))
}
