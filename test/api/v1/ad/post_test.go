package adTest

import (
	"bytes"
	"dcard-assignment/cmd/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPost(t *testing.T) {
	testData := getTestData()

	for _, data := range testData {
		expectedOutput := data.Output
		requestAd := getAdType(&data)
		requestJson, err := json.Marshal(requestAd)
		utils.CheckError(err)
		response, err := http.Post("http://localhost:3000/api/v1/ad", "application/json", bytes.NewReader(requestJson))
		utils.CheckError(err)

		body, err := io.ReadAll(response.Body)
		utils.CheckError(err)

		msg := responseMessage{}

		err = json.Unmarshal(body, &msg)
		utils.CheckError(err)

		// Check if response match expected output
		assertion := assert.New(t)
		pass := assertion.Equal(expectedOutput, msg.Msg)
		if !pass {
			fmt.Println(string(requestJson))
			break
		}

	}

	// req := httptest.NewRequest("POST", "localhost:3000/api/v1/ad", nil)
}

// Function: Remove output field to match ad type.
func getAdType(adTest *AdTest) *Ad {
	requestAd := Ad{}
	requestAd.Title = adTest.Title
	requestAd.StartAt = adTest.StartAt
	requestAd.EndAt = adTest.EndAt
	requestAd.Conditions = adTest.Conditions
	return &requestAd
}

// Function: Read testing data from json file.
func getTestData() []AdTest {
	file := "post.json"
	testData := []AdTest{}

	data, err := os.ReadFile(file)
	utils.CheckError(err)

	err = json.Unmarshal(data, &testData)
	utils.CheckError(err)
	return testData
}
