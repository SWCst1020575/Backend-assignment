package connect

import (
	"database/sql"
	"dcard-assignment/cmd/config"
	"fmt"

	_ "github.com/lib/pq"
)

// Save database connection in global variable.
var dbConnection *sql.DB

func DBconnect() {
	dbConfig := config.GetDbConfig()
	connectionStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbConfig.Addr, dbConfig.Port, dbConfig.Username, dbConfig.Password, dbConfig.Name)
	db, err := sql.Open("postgres", connectionStr)
	checkError(err)

	err = db.Ping()
	checkError(err)

	fmt.Println("Successfully created connection to database.")
	dbConnection = db
	checkTableExist()
}
func DBclose() {
	err := dbConnection.Close()
	checkError(err)
	fmt.Println("Connection is closed.")
}

// Execute SQL query (without return value).
func execSQL(sqlQuery string) {
	_, err := dbConnection.Exec(sqlQuery)
	checkError(err)
}

// Execute SQL query (with return).
// return: sql.Rows
func selectFromSQL(sqlQuery string) *sql.Rows {
	rows, err := dbConnection.Query(sqlQuery)
	checkError(err)
	return rows
}

// The query string to create table in database.
const initCreateTableQuery = `CREATE TABLE Ad (
    							ID SERIAL PRIMARY KEY,
    							Title text NOT NULL,
								StartAt timestamp NOT NULL,
								EndAt timestamp NOT NULL,
								Age int,
								Gender boolean,
							);
							CREATE TABLE Country (
    							ID NOT NULL references Ad(ID),
								Country char(2)
							);
							CREATE TABLE Platform (
    							ID NOT NULL references Ad(ID),
								Platform char(7)
							);`

// Check if ad table exists.
func checkTableExist() {
	_, check := dbConnection.Query("SELECT * FROM ad;")
	if check != nil {
		dbConnection.Exec(initCreateTableQuery)
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
