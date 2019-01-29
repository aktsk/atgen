package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
)

// DummyRouter is router for test
func DummyRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/v1/", func(w http.ResponseWriter, r *http.Request) {})
	return router
}
