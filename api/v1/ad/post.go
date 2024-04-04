package ad

import (
	"database/sql"
	"dcard-assignment/cmd/connect"
	"dcard-assignment/cmd/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/biter777/countries"
)

// Function: Handle post method of ad
func Post(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	body, err := io.ReadAll(request.Body)
	utils.CheckError(err)

	var newAd Ad
	err = json.Unmarshal(body, &newAd)
	utils.CheckError(err)

	success := dbInsert(&newAd)
	if !success {
		postFormatError(writer)
		return
	}

	resp := make(map[string]string)
	resp["message"] = "Success"
	jsonResp, err := json.Marshal(resp)
	utils.CheckError(err)

	// Return success
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(jsonResp)
	utils.CheckError(err)

}

// Function: Insert new ad to database, use transcetion to ensure two table are consistant.
// return: false, if request format is wrong.
func dbInsert(ad *Ad) bool {
	// Initialize query string

	if !checkNotNull(ad) {
		return false
	}
	const adColumns = "Title, StartAt, EndAt, AgeStart, AgeEnd, Male, Female, PlatformAndroid, PlatformIos, PlatformWeb"
	const queryValue = "$1, $2, $3, $4, $5, $6, $7, $8, $9, $10"
	query := fmt.Sprintf("INSERT INTO Ad (%s) VALUES (%s) RETURNING ID;", adColumns, queryValue)
	extendCondition := getExtendCondition(ad)
	if !checkCountryValid(ad) {
		return false
	}
	if !checkAgeValid(ad) {
		return false
	}
	// Start transection
	db := connect.GetDBconnection()
	tx, err := db.Begin()
	if !utils.TransectionCheckError(err, tx) {
		return false
	}

	// Insert new ad
	var newAdID int
	err = tx.QueryRow(query, ad.Title, ad.StartAt, ad.EndAt, ad.Conditions.AgeStart, ad.Conditions.AgeEnd,
		extendCondition.Male, extendCondition.Female,
		extendCondition.PlatformAndroid, extendCondition.PlatformIos, extendCondition.PlatformWeb).Scan(&newAdID)
	if !utils.TransectionCheckError(err, tx) {
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
	if !utils.TransectionCheckError(err, tx) {
		fmt.Println("Commit error.")
		return false
	}

	fmt.Println("Transection commit success.")

	return true
}

// Function: Handle wrong format json
func postFormatError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["message"] = "Failed"
	jsonResp, err := json.Marshal(resp)
	utils.CheckError(err)
	w.Write(jsonResp)
}

// Function: Parse gender and platform to boolean variable.
func getExtendCondition(ad *Ad) *extendCondition {
	extend := extendCondition{false, false, false, false, false}
	for _, gender := range ad.Conditions.Gender {
		switch gender {
		case "M", "m":
			extend.Male = true
		case "F", "f":
			extend.Female = true
		}
	}
	for _, platform := range ad.Conditions.Platform {
		lowerPlatform := strings.ToLower(platform)
		switch lowerPlatform {
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

// Function: Insert country data in decomposed table.
func handleCountryInsert(countrys []string, id int, tx *sql.Tx) bool {
	query := "INSERT INTO Country (ID, Country) VALUES ($1, $2);"
	stmt, err := tx.Prepare(query)
	if !utils.TransectionCheckError(err, tx) {
		return false
	}
	countrys = removeDuplicatedCountries(countrys)
	for _, country := range countrys {
		_, err := stmt.Exec(id, country)
		if !utils.TransectionCheckError(err, tx) {
			return false
		}
	}
	return true
}

// Function: Check if country is exactly two characters and if exist in ISO 3166 to prevent invalid data inserting.
func checkCountryValid(ad *Ad) bool {
	for _, country := range ad.Conditions.Country {
		if len(country) != 2 {
			return false
		}
		if !countries.ByName(country).IsValid() {
			return false
		}
	}
	return true
}

func checkNotNull(ad *Ad) bool {
	if ad.Title == "" {
		return false
	}
	if ad.StartAt.IsZero() {
		return false
	}
	if ad.EndAt.IsZero() {
		return false
	}
	return true
}

// Function: AgeStart, AgeEnd must be both set or both unset, and AgeEnd > AgeStart
func checkAgeValid(ad *Ad) bool {
	if ad.Conditions.AgeStart == 0 && ad.Conditions.AgeEnd != 0 {
		return false
	}
	if ad.Conditions.AgeStart != 0 && ad.Conditions.AgeEnd == 0 {
		return false
	}

	return ad.Conditions.AgeEnd >= ad.Conditions.AgeStart
}

// Function: Remove duplicated countries.
func removeDuplicatedCountries(countrys []string) []string {
	allKeys := make(map[string]bool)
	newListofCountries := []string{}
	for _, item := range countrys {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			newListofCountries = append(newListofCountries, item)
		}
	}
	return newListofCountries
}
