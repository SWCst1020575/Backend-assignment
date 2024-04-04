package ad

import (
	"dcard-assignment/cmd/connect"
	"dcard-assignment/cmd/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/biter777/countries"
)

// Function: Handle get method of ad
func Get(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	// Parse url in get
	urlQuery := request.URL.Query()
	search := parseSearch(&urlQuery)
	if search == nil {
		getFormatError(writer)
		return
	}

	query := parseQuery(search)

	db := connect.GetDBconnection()
	rows, err := db.Query(query)
	utils.CheckError(err)
	defer rows.Close()

	responseData := []getResponseData{}
	for rows.Next() {
		var title string
		var endat time.Time
		err := rows.Scan(&title, &endat)
		utils.CheckError(err)
		responseData = append(responseData, getResponseData{title, endat})
	}
	response := getResponse{responseData}
	jsonResp, err := json.Marshal(response)
	utils.CheckError(err)

	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(jsonResp)
	utils.CheckError(err)
	fmt.Println("Get success.")
}

// Function: Parse str to int and block invalid request, and prevent sql injection as well
func parseSearch(search *url.Values) *SearchAd {
	s := SearchAd{0, 0, 0, "", "", ""}
	isParameterEmpty := false

	for key, elements := range *search {
		lowerElement := strings.ToLower(elements[0])
		switch key {
		case "offset":
			val, err := strconv.Atoi(elements[0])
			if err != nil {
				return nil
			}
			s.Offset = val
		case "age":
			val, err := strconv.Atoi(elements[0])
			if err != nil {
				return nil
			}
			s.Age = val

		case "limit":
			val, err := strconv.Atoi(elements[0])
			if err != nil {
				return nil
			}
			s.Limit = val
		case "gender":
			s.Gender = lowerElement
			if s.Gender != "f" && s.Gender != "m" && s.Gender != "" {
				return nil
			}
		case "country":
			s.Country = lowerElement
			if len(s.Country) != 2 {
				return nil
			}
			if !countries.ByName(s.Country).IsValid() {
				return nil
			}
		case "platform":
			s.Platform = lowerElement
			if s.Platform != "ios" && s.Platform != "android" && s.Platform != "web" && s.Platform != "" {
				return nil
			}
		}
		isParameterEmpty = true
	}
	if !isParameterEmpty || s.Offset == 0 || s.Limit == 0 {
		return nil
	}
	return &s
}

// Function: Handle wrong url format
func getFormatError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["message"] = "Failed"
	jsonResp, err := json.Marshal(resp)
	utils.CheckError(err)
	w.Write(jsonResp)
}

// Function: Parse query string.
func parseQuery(search *SearchAd) string {
	condition := "WHERE NOW()>startat"
	if search.Age != 0 {
		condition += fmt.Sprintf(" AND ((A.agestart <= %d AND A.ageend >= %d) OR (A.agestart = 0 AND A.ageend = 0))", search.Age, search.Age)
	}
	if search.Country != "" {
		condition += fmt.Sprintf(" AND ((C.country ILIKE '%s') OR (C.country IS NULL))", search.Country)
	}
	if search.Platform != "" {
		platform := "platform" + search.Platform
		condition += fmt.Sprintf(" AND (A.%s = true OR (A.platformandroid = false AND A.platformios = false AND A.platformweb = false))", platform)
	}

	if search.Gender == "m" {
		condition += " AND (A.male = true OR (A.male = false AND a.female = false))"
	} else if search.Gender == "f" {
		condition += " AND (A.female = true OR (A.male = false AND a.female = false))"
	}
	var query string
	if search.Country != "" {
		query = fmt.Sprintf("SELECT A.title, A.endat FROM Ad A JOIN Country C ON C.id = A.id %s ORDER BY A.endat ASC LIMIT %d OFFSET %d;", condition, search.Limit, search.Offset)
	} else {
		query = fmt.Sprintf("SELECT A.title, A.endat FROM Ad A %s ORDER BY A.endat ASC LIMIT %d OFFSET %d;", condition, search.Limit, search.Offset)
	}
	return query
}
