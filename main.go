package main

import (
	"dcard-assignment/cmd/connect"
)

func main() {
	connect.DBconnect()

	connect.DBclose()

}
