package utils

import (
	"database/sql"
	"fmt"
)

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

// Transection rollback if error occur.
func TransectionCheckError(err error, tx *sql.Tx) bool {
	if err != nil {
		fmt.Println(err.Error())
		err = tx.Rollback()
		CheckError(err)
		fmt.Println("Transection failed.")
		return false
	}
	return true
}
