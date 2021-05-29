package router

import (
	"accurics-test/router/muxrouter"
	"net/http"
)

func NewRouter() http.Handler {
	return muxrouter.GetMuxRouter()
}
