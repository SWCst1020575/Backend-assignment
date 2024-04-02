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

// Define ad information for testing
type AdTest struct {
	Title      string     `json:"title,omitempty"`
	StartAt    time.Time  `json:"startAt,omitempty"`
	EndAt      time.Time  `json:"endAt,omitempty"`
	Conditions *Condition `json:"conditions,omitempty"`
	Output     string     `json:"output,omitempty"`
}

type AdTestRead struct {
	Test []AdTest `json:"test"`
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

type responseMessage struct {
	Msg string `json:"message"`
}
