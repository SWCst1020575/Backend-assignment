package adTest

import (
	"time"
)

// Define ad information
type Ad struct {
	Title      string     `json:"title,omitempty"`
	StartAt    time.Time  `json:"startAt,omitempty"`
	EndAt      time.Time  `json:"endAt,omitempty"`
	Conditions *Condition `json:"conditions,omitempty"`
}

// Define reading testing data strcture for post testing
type AdPostTest struct {
	Title      string     `json:"title,omitempty"`
	StartAt    time.Time  `json:"startAt,omitempty"`
	EndAt      time.Time  `json:"endAt,omitempty"`
	Conditions *Condition `json:"conditions,omitempty"`
	Output     string     `json:"output,omitempty"`
}

// Define reading testing data strctur for get testing
type AdGetTest struct {
	Offset   int         `json:"offset,omitempty"`
	Limit    int         `json:"limit,omitempty"`
	Age      int         `json:"age,omitempty"`
	Gender   string      `json:"gender,omitempty"`
	Country  string      `json:"country,omitempty"`
	Platform string      `json:"platform,omitempty"`
	Output   getResponse `json:"output,omitempty"`
}

type Condition struct {
	AgeStart int      `json:"ageStart,omitempty"`
	AgeEnd   int      `json:"ageEnd,omitempty"`
	Gender   []string `json:"gender,omitempty"`
	Country  []string `json:"country,omitempty"`
	Platform []string `json:"platform,omitempty"`
}

// Define ad get method interface
type SearchAd struct {
	Offset   int    `json:"offset,omitempty"`
	Limit    int    `json:"limit,omitempty"`
	Age      int    `json:"age,omitempty"`
	Gender   string `json:"gender,omitempty"`
	Country  string `json:"country,omitempty"`
	Platform string `json:"platform,omitempty"`
}

type responseMessage struct {
	Msg string `json:"message"`
}

type getResponseData struct {
	Title string    `json:"title"`
	EndAt time.Time `json:"endAt"`
}

type getResponse struct {
	Items []getResponseData `json:"items"`
}
