package ad

import (
	"database/sql"
	"dcard-assignment/cmd/connect"
	. "dcard-assignment/cmd/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Function: Handle post method of ad
func Post(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	body, err := io.ReadAll(request.Body)
	CheckError(err)

	var newAd Ad
	err = json.Unmarshal(body, &newAd)
	CheckError(err)

	success := dbInsert(&newAd)
	if !success {
		postFormatError(writer)
		return
	}

	response, err := json.Marshal(newAd)
	CheckError(err)

	// Return success
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(response)
	CheckError(err)

}

// Function: Insert new ad to database, use transcetion to ensure two table are consistant.
// return: false, if request format is wrong.
func dbInsert(ad *Ad) bool {
	// Initialize query string
	const adColumns = "Title, StartAt, EndAt, AgeStart, AgeEnd, Male, Female, PlatformAndroid, PlatformIos, PlatformWeb"
	const queryValue = "$1, $2, $3, $4, $5, $6, $7, $8, $9, $10"
	query := fmt.Sprintf("INSERT INTO Ad (%s) VALUES (%s) RETURNING ID;", adColumns, queryValue)
	extendCondition := getExtendCondition(ad)

	// Start transection
	db := connect.GetDBconnection()
	tx, err := db.Begin()
	if !TransectionCheckError(err, tx) {
		return false
	}

	// Insert new ad
	var newAdID int
	err = tx.QueryRow(query, ad.Title, ad.StartAt, ad.EndAt, ad.Conditions.AgeStart, ad.Conditions.AgeEnd,
		extendCondition.Male, extendCondition.Female,
		extendCondition.PlatformAndroid, extendCondition.PlatformIos, extendCondition.PlatformWeb).Scan(&newAdID)
	if !TransectionCheckError(err, tx) {
		fmt.Println("Ad insert error.")
		return false
	}

	// Insert countries in another table
	if !handleCountryInsert(ad.Conditions.Country, newAdID, tx) {
		fmt.Println("Country insert error.")
		return false
	}

	// Transection commit
	err = tx.Commit()
	if !TransectionCheckError(err, tx) {
		fmt.Println("Commit error.")
		return false
	}

	fmt.Println("Transection commit success.")
	fmt.Println(query)

	return true
}

// Function: Handle wrong format json
func postFormatError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["message"] = "Bad Request"
	jsonResp, err := json.Marshal(resp)
	CheckError(err)
	w.Write(jsonResp)
}

// Function: Parse gender and platform to boolean variable.
func getExtendCondition(ad *Ad) *extendCondition {
	extend := extendCondition{false, false, false, false, false}
	for _, gender := range ad.Conditions.Gender {
		switch gender {
		case "M":
			extend.Male = true
		case "F":
			extend.Female = true
		}
	}
	for _, platform := range ad.Conditions.Platform {
		switch platform {
		case "android":
			extend.PlatformAndroid = true
		case "ios":
			extend.PlatformIos = true
		case "web":
			extend.PlatformWeb = true
		}
	}
	return &extend
}

func handleCountryInsert(countrys []string, id int, tx *sql.Tx) bool {

	query := "INSERT INTO Country (ID, Country) VALUES ($1, $2);"
	stmt, err := tx.Prepare(query)
	if !TransectionCheckError(err, tx) {
		return false
	}
	for _, country := range countrys {
		_, err := stmt.Exec(id, country)
		if !TransectionCheckError(err, tx) {
			return false
		}
	}
	return true
}
