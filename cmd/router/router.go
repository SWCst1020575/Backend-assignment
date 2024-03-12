package router

import (
	"dcard-assignment/api/v1/ad"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	// TODO: Design api router (/api/v1/ad) post and get.
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/ad", ad.Post).Methods("POST")
	return r
}
