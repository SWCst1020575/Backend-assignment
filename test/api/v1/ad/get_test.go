package adTest

import (
	"bytes"
	"dcard-assignment/cmd/utils"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"

	"testing"

	"github.com/schollz/progressbar/v3"
	"github.com/stretchr/testify/assert"
)

const AD_TESTING_DATA_FILE = "adTestingData.json"

// False when data has been in database already
const IS_DATA_PREPARE = true

func TestGet(t *testing.T) {
	dataPrepare()

	testData := getGetTestData()
	for _, data := range testData {
		expectedOutput := data.Output
		requestAd := getSearchAdType(&data)
		requestUrl := parseGetUrl(requestAd)

		requestJson, err := json.Marshal(requestAd)
		utils.CheckError(err)
		response, err := http.Get(*requestUrl)
		utils.CheckError(err)

		body, err := io.ReadAll(response.Body)
		utils.CheckError(err)

		resp := getResponse{}

		err = json.Unmarshal(body, &resp)
		utils.CheckError(err)

		// Check if response match expected output
		assertion := assert.New(t)
		pass := assertion.Equal(expectedOutput, resp)
		if !pass {
			fmt.Println(string(requestJson))
			break
		}

	}
}

const TEST_DATA_NUM = 5000
const CHAR_SET = "aAbBcCdDeEfFgGhHiIjJkKlLmMnNoOpPqQrRsStTuUvVwWxXyYzZ"

func dataPrepare() {
	// If data has generated then read them, or generate random data.
	// Then inserts them into database.

	if !IS_DATA_PREPARE {
		return
	}
	_, err := os.Stat(AD_TESTING_DATA_FILE)
	if err == nil {
		fmt.Println("Using exist data.")
		dataRead()
	} else {
		fmt.Println("Generate random data.")
		dataGen()
	}
}

// Function: Generate random ad data for testing get method.
func dataGen() {
	var testingData [TEST_DATA_NUM]Ad
	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < TEST_DATA_NUM; i++ {
		testingData[i].Title = textGen(16, randGen)
		testingData[i].StartAt = time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)
		testingData[i].EndAt = time.Date(2024, 4, randGen.Intn(30)+1, randGen.Intn(24), randGen.Intn(60), randGen.Intn(30)+1, 0, time.UTC)
		testingData[i].Conditions = &Condition{}

		if randGen.Float32() < 0.7 {
			testingData[i].Conditions.AgeStart = randGen.Intn(30) + 15
			testingData[i].Conditions.AgeEnd = randGen.Intn(30) + 1 + testingData[i].Conditions.AgeStart
		}
		if randGen.Float32() < 0.6 {
			testingData[i].Conditions.Gender = append(testingData[i].Conditions.Gender, "F")
		}
		if randGen.Float32() < 0.6 {
			testingData[i].Conditions.Gender = append(testingData[i].Conditions.Gender, "M")
		}
		if randGen.Float32() < 0.75 {
			if randGen.Float32() < 0.5 {
				testingData[i].Conditions.Platform = append(testingData[i].Conditions.Platform, "ios")
			}
			if randGen.Float32() < 0.5 {
				testingData[i].Conditions.Platform = append(testingData[i].Conditions.Platform, "android")
			}
			if randGen.Float32() < 0.5 {
				testingData[i].Conditions.Platform = append(testingData[i].Conditions.Platform, "web")
			}
		}
		if randGen.Float32() < 0.75 {
			if randGen.Float32() < 0.8 {
				testingData[i].Conditions.Country = append(testingData[i].Conditions.Country, "TW")
			}
			if randGen.Float32() < 0.5 {
				testingData[i].Conditions.Country = append(testingData[i].Conditions.Country, "JP")
			}
			if randGen.Float32() < 0.5 {
				testingData[i].Conditions.Country = append(testingData[i].Conditions.Country, "HK")
			}
		}
	}

	bar := progressbar.Default(TEST_DATA_NUM)

	outputJson, err := json.Marshal(testingData)
	utils.CheckError(err)

	err = os.WriteFile(AD_TESTING_DATA_FILE, outputJson, 0644)
	utils.CheckError(err)

	for _, data := range testingData {
		bar.Add(1)
		requestJson, err := json.Marshal(data)
		utils.CheckError(err)
		response, err := http.Post("http://localhost:3000/api/v1/ad", "application/json", bytes.NewReader(requestJson))
		utils.CheckError(err)

		body, err := io.ReadAll(response.Body)
		utils.CheckError(err)

		msg := responseMessage{}

		err = json.Unmarshal(body, &msg)
		utils.CheckError(err)

	}
	fmt.Println("Insertion complete.")

}

// Function: Read data from random generating ad data previouslyfor testing get method.
func dataRead() {
	var testingData [TEST_DATA_NUM]Ad

	data, err := os.ReadFile(AD_TESTING_DATA_FILE)
	utils.CheckError(err)

	err = json.Unmarshal(data, &testingData)
	utils.CheckError(err)

	bar := progressbar.Default(TEST_DATA_NUM)

	for _, data := range testingData {
		bar.Add(1)
		requestJson, err := json.Marshal(data)
		utils.CheckError(err)
		response, err := http.Post("http://localhost:3000/api/v1/ad", "application/json", bytes.NewReader(requestJson))
		utils.CheckError(err)

		body, err := io.ReadAll(response.Body)
		utils.CheckError(err)

		msg := responseMessage{}

		err = json.Unmarshal(body, &msg)
		utils.CheckError(err)

	}

	fmt.Println("Insertion complete.")
}

// Function: Generate random text.
func textGen(n int, randGen *rand.Rand) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = CHAR_SET[randGen.Intn(len(CHAR_SET))]
	}
	return string(b)
}

// Function: Read testing data from json file.
func getGetTestData() []AdGetTest {
	file := "get.json"
	testData := []AdGetTest{}

	data, err := os.ReadFile(file)
	utils.CheckError(err)

	err = json.Unmarshal(data, &testData)
	utils.CheckError(err)
	return testData
}
func getSearchAdType(adTest *AdGetTest) *SearchAd {
	requestAd := SearchAd{}
	requestAd.Offset = adTest.Offset
	requestAd.Limit = adTest.Limit
	requestAd.Age = adTest.Age
	requestAd.Gender = adTest.Gender
	requestAd.Country = adTest.Country
	requestAd.Platform = adTest.Platform
	return &requestAd
}

func parseGetUrl(searchAd *SearchAd) *string {
	url := fmt.Sprintf("http://localhost:3000/api/v1/ad?offset=%d&limit=%d", searchAd.Offset, searchAd.Limit)
	if searchAd.Age != 0 {
		url += fmt.Sprintf("&age=%d", searchAd.Age)
	}
	if searchAd.Gender != "" {
		url += fmt.Sprintf("&gender=%s", searchAd.Gender)
	}
	if searchAd.Country != "" {
		url += fmt.Sprintf("&country=%s", searchAd.Country)
	}
	if searchAd.Platform != "" {
		url += fmt.Sprintf("&platform=%s", searchAd.Platform)
	}

	return &url
}
