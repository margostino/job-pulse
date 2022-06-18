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
}
