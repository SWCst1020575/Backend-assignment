package ad

import (
	"dcard-assignment/cmd/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
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

func parseSearch(search *url.Values) *SearchAd {
	s := SearchAd{}
	isParameterEmpty := false
	for key, elements := range *search {
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
			s.Gender = elements[0]
		case "country":
			s.Country = elements[0]
		case "platform":
			s.Platform = elements[0]
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
