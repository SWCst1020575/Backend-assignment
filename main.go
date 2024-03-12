package main

import (
	"dcard-assignment/cmd/connect"
	"dcard-assignment/cmd/router"
	"dcard-assignment/cmd/utils"
	"net/http"
)

func main() {
	connect.DBconnect()

	r := router.NewRouter()
	err := http.ListenAndServe(":3000", r)
	utils.CheckError(err)

	connect.DBclose()

}
