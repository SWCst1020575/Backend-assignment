package ad

import (
	"encoding/json"
	"io"
	"net/http"
)

// Handle post method of ad
func Post(writer http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)
	checkError(err)

	var newAd Ad
	err = json.Unmarshal(body, &newAd)
	checkError(err)

	// TODO: Save newAd to database

	defer request.Body.Close()

	response, err := json.Marshal(newAd)
	checkError(err)

	// Return success
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(response)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
