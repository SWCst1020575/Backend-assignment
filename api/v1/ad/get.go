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
	s := SearchAd{}
	isParameterEmpty := false
	isOffsetExist := false
	isLimitExist := false
	for key, elements := range *search {
		lowerElement := strings.ToLower(elements[0])
		switch key {
		case "offset":
			val, err := strconv.Atoi(elements[0])
			if err != nil {
				return nil
			}
			isOffsetExist = true
			s.Offset = val
		case "age":
			val, err := strconv.Atoi(elements[0])
			if err != nil {
				return nil
			}
			isLimitExist = true
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
	if !isParameterEmpty || !isLimitExist || !isOffsetExist {
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
		condition += fmt.Sprintf(" AND A.agestart <= %d AND A.ageend >= %d", search.Age, search.Age)
	}
	if search.Country != "" {
		condition += fmt.Sprintf(" AND C.country ILIKE '%s'", search.Country)
	}
	if search.Platform != "" {
		platform := "platform" + search.Platform
		condition += fmt.Sprintf(" AND A.%s = true", platform)
	}

	if search.Gender == "m" {
		condition += " AND A.male = true"
	} else if search.Gender == "f" {
		condition += " AND A.female = true"
	}
	query := fmt.Sprintf("SELECT A.title, A.endat FROM Ad A JOIN Country C ON C.id = A.id %s ORDER BY A.endat ASC LIMIT %d OFFSET %d;", condition, search.Limit, search.Offset)
	return query
}
