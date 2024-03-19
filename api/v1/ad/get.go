package ad

import (
	"dcard-assignment/cmd/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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

	fmt.Println(search)
}

// Function: Parse str to int and block invalid request, and prevent sql injection as well
func parseSearch(search *url.Values) *SearchAd {
	s := SearchAd{}
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
			if s.Gender != "f" && s.Gender != "m" {
				return nil
			}
		case "country":
			s.Country = lowerElement
			if len(s.Country) != 2 {
				return nil
			}
		case "platform":
			s.Platform = lowerElement
			if s.Platform != "ios" && s.Platform != "android" && s.Platform != "web" {
				return nil
			}
		}
		isParameterEmpty = true
	}
	if !isParameterEmpty {
		return nil
	}
	return &s
}

// Function: Handle wrong url format
func getFormatError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["message"] = "Query error"
	jsonResp, err := json.Marshal(resp)
	utils.CheckError(err)
	w.Write(jsonResp)
}
