package handlers

import "github.com/gorilla/mux"

func DummyRouter() *mux.Router {
	router := mux.NewRouter()
	return router
}
