package adTest

import (
	"bytes"
	"dcard-assignment/cmd/utils"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"testing"

	"github.com/schollz/progressbar/v3"
)

func TestGet(t *testing.T) {
	dataGen()

}

const TEST_DATA_NUM = 10000
const CHAR_SET = "aAbBcCdDeEfFgGhHiIjJkKlLmMnNoOpPqQrRsStTuUvVwWxXyYzZ"

// Function: Generate random data to ad for testing get method.
func dataGen() {
	var testingData [TEST_DATA_NUM]Ad
	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < TEST_DATA_NUM; i++ {
		testingData[i].Title = textGen(16, randGen)
		testingData[i].StartAt = time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)
		testingData[i].EndAt = time.Date(2024, 4, randGen.Intn(30)+1, randGen.Intn(24), randGen.Intn(60), randGen.Intn(30)+1, 0, time.UTC)
		testingData[i].Conditions = &Condition{}
		testingData[i].Conditions.AgeStart = randGen.Intn(30) + 15
		testingData[i].Conditions.AgeEnd = randGen.Intn(30) + 1 + testingData[i].Conditions.AgeStart
		if randGen.Intn(2) == 1 {
			testingData[i].Conditions.Gender = append(testingData[i].Conditions.Gender, "F")
		}
		if randGen.Intn(2) == 1 {
			testingData[i].Conditions.Gender = append(testingData[i].Conditions.Gender, "M")
		}
		if randGen.Intn(2) == 1 {
			testingData[i].Conditions.Platform = append(testingData[i].Conditions.Platform, "ios")
		}
		if randGen.Intn(2) == 1 {
			testingData[i].Conditions.Platform = append(testingData[i].Conditions.Platform, "android")
		}
		if randGen.Intn(2) == 1 {
			testingData[i].Conditions.Platform = append(testingData[i].Conditions.Platform, "web")
		}
		if randGen.Float32() > 0.8 {
			testingData[i].Conditions.Country = append(testingData[i].Conditions.Country, "TW")
		}
		if randGen.Float32() > 0.5 {
			testingData[i].Conditions.Country = append(testingData[i].Conditions.Country, "JP")
		}
		if randGen.Float32() > 0.5 {
			testingData[i].Conditions.Country = append(testingData[i].Conditions.Country, "HK")
		}
	}

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
	fmt.Println("Generating complete.")

}

// Function: Generate random text.
func textGen(n int, randGen *rand.Rand) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = CHAR_SET[randGen.Intn(len(CHAR_SET))]
	}
	return string(b)
}
