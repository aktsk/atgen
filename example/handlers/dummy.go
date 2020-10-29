package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
)

// DummyRouter is router for test
func GetDummyRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/v1/", func(w http.ResponseWriter, r *http.Request) {})
	return router
}
