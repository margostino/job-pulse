package utils

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

var timeRegEx, _ = regexp.Compile("^[A-Za-z,.;\\s_-]*[0-9]{1,2}\\s(hour|hours|day|days|week|weeks|month|months)\\sago")

func GetSubtractTime(value int, unit string) (time.Duration, error) {
	switch unit {
	case "hour":
	case "hours":
		return -time.Hour * time.Duration(value), nil
	case "day":
	case "days":
		return -time.Hour * time.Duration(value) * 24, nil
	case "week":
	case "weeks":
		return -time.Hour * time.Duration(value) * 24 * 7, nil
	case "month":
	case "months":
		return -time.Hour * time.Duration(value) * 24 * 7 * 30, nil
	}

	return 0, errors.New("unit not valid")
}

func ToBeginningOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location()).UTC()
}

func GetRawPostDateOrDefault(index int, list []string) string {
	if len(list) >= index+1 {
		return list[index]
	} else if len(list) > 1 && timeRegEx.MatchString(list[2]) {
		allParts := strings.SplitN(list[2], ",", -1)
		for _, value := range allParts {
			if strings.Contains(value, " ago") {
				parts := strings.SplitN(strings.TrimSpace(value), " ", -1)
				if len(parts) < 4 {
					println("")
				}
				return strings.Join(parts[1:4], " ")
			}
		}
	}
	return "1 hour ago"
}
