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
type extendCondition struct {
	Male            bool
	Female          bool
	PlatformAndroid bool
	PlatformIos     bool
	PlatformWeb     bool
}
