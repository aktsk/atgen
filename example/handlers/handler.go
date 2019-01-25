package handlers

import "github.com/gorilla/mux"

// DummyRouter is router for test
func DummyRouter() *mux.Router {
	router := mux.NewRouter()
	return router
}
