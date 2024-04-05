package router

import (
	"dcard-assignment/api/v1/ad"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/ad", ad.Post).Methods("POST")
	r.HandleFunc("/api/v1/ad", ad.Get).Methods("GET")
	return r
}
