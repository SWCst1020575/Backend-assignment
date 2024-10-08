package ad

import (
	"time"
)

// Define ad information
type Ad struct {
	Title      string    `json:"title"`
	StartAt    time.Time `json:"startAt"`
	EndAt      time.Time `json:"endAt"`
	Conditions Condition `json:"conditions"`
}

type Condition struct {
	AgeStart int      `json:"ageStart"`
	AgeEnd   int      `json:"ageEnd"`
	Gender   []string `json:"gender"`
	Country  []string `json:"country"`
	Platform []string `json:"platform"`
}

// Define ad get method interface
type SearchAd struct {
	Offset   int    `json:"offset"`
	Limit    int    `json:"limit"`
	Age      int    `json:"age"`
	Gender   string `json:"gender"`
	Country  string `json:"country"`
	Platform string `json:"platform"`
}

// For translation from string to boolean
type extendCondition struct {
	Male            bool
	Female          bool
	PlatformAndroid bool
	PlatformIos     bool
	PlatformWeb     bool
}

type getResponseData struct {
	Title string    `json:"title"`
	EndAt time.Time `json:"endAt"`
}

type getResponse struct {
	Items []getResponseData `json:"items"`
}
