package domain

import "time"

type JobPost struct {
	Position    string
	Company     string
	Location    string
	Benefit     string
	Link        string
	RawPostDate string
	PostDate    time.Time
	Latitude    float64
	Longitude   float64
}

type Stats struct {
	StartTime     time.Time
	PositionInput string
	LocationInput string
}
