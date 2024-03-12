package connect

import (
	"database/sql"
	"dcard-assignment/cmd/config"
	. "dcard-assignment/cmd/utils"
	"fmt"

	_ "github.com/lib/pq"
)

// Variable: Save database connection in global variable.
var dbConnection *sql.DB

func DBconnect() {
	dbConfig := config.GetDbConfig()
	connectionStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbConfig.Addr, dbConfig.Port, dbConfig.Username, dbConfig.Password, dbConfig.Name)
	db, err := sql.Open("postgres", connectionStr)
	CheckError(err)

	err = db.Ping()
	CheckError(err)

	fmt.Println("Successfully created connection to database.")
	dbConnection = db
	checkTableExist()
}
func DBclose() {
	err := dbConnection.Close()
	CheckError(err)
	fmt.Println("Connection is closed.")
}

// Function: Get database connection
func GetDBconnection() *sql.DB {
	return dbConnection
}

// Constant: The query string to create table in database.
const initCreateTableQuery = `CREATE TABLE Ad (
    							ID SERIAL PRIMARY KEY,
    							Title text NOT NULL,
								StartAt timestamp NOT NULL,
								EndAt timestamp NOT NULL,
								AgeStart int,
								AgeEnd int,
								Male boolean,
								Female boolean,
								PlatformAndroid boolean,
								PlatformIos boolean,
								PlatformWeb boolean,
							);
							CREATE TABLE Country (
    							ID NOT NULL references Ad(ID),
								Country char(2)
							);`

// Function: Check if ad table exists.
func checkTableExist() {
	_, check := dbConnection.Query("SELECT * FROM ad;")
	if check != nil {
		dbConnection.Exec(initCreateTableQuery)
	}
}
