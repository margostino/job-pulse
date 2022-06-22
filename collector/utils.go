package collector

import (
	"github.com/margostino/job-pulse/domain"
	"github.com/margostino/job-pulse/utils"
	"strconv"
	"strings"
	"time"
)

func buildJobPost(jobTextParts []string, link string, rawPostDate string, postDate time.Time) *domain.JobPost {
	return &domain.JobPost{
		Position:    utils.GetOrDefault(0, jobTextParts),
		Company:     utils.GetOrDefault(1, jobTextParts),
		Location:    utils.GetOrDefault(2, jobTextParts),
		Benefit:     utils.GetOrDefault(3, jobTextParts),
		Link:        link,
		RawPostDate: rawPostDate,
		PostDate:    postDate,
	}
}

func timeDurationSince(startTime time.Time) float64 {
	return time.Now().UTC().Sub(startTime).Minutes()
}

func calculateJobPostDate(rawPostDate string) time.Time {
	dateParts := strings.SplitN(rawPostDate, " ", -1)
	value, _ := strconv.Atoi(dateParts[0])
	unit := dateParts[1]
	duration, _ := utils.GetSubtractTime(value, unit)
	return utils.ToBeginningOfDay(time.Now().UTC().Add(duration))
}
