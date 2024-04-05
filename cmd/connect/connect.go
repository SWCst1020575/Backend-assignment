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

// Function: Get database connection
func GetDBconnection() *sql.DB {
	return dbConnection
}

// Constant: The query string to create table in database.
var initCreateTableQuery = [7]string{
	`CREATE TABLE Ad (
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
								PlatformWeb boolean
	);`,
	`CREATE TABLE Country (
    	ID int NOT NULL references Ad(ID),
		Country char(2)
	);`,
	"CREATE INDEX country_index ON Country USING HASH (country);",
	"CREATE INDEX ad_agestart_index ON Ad USING BTREE (agestart);",
	"CREATE INDEX ad_agestart_index ON Ad USING BTREE (ageend);",
	"CREATE INDEX ad_endat_index ON Ad USING BTREE (endat);",
	"CREATE INDEX ad_startat_index ON Ad USING BTREE (startat);"}

// Function: Check if ad table exists.
func checkTableExist() {

	for i := 0; i < 7; i++ {
		var check error
		if i == 0 {
			_, check = dbConnection.Query("SELECT * FROM Ad;")
		}
		if i == 1 {
			_, check = dbConnection.Query("SELECT * FROM Country;")
		}
		if check != nil {
			_, err := dbConnection.Exec(initCreateTableQuery[i])
			CheckError(err)
			fmt.Println("Create.")
		}
	}

}
