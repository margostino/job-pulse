package scrapper

import (
	"github.com/margostino/job-pulse/utils"
	"strconv"
	"strings"
	"time"
)

func CalculateJobPostDate(rawPostDate string) time.Time {
	dateParts := strings.SplitN(rawPostDate, " ", -1)
	value, _ := strconv.Atoi(dateParts[0])
	unit := dateParts[1]
	duration, _ := utils.GetSubtractTime(value, unit)
	return utils.ToBeginningOfDay(time.Now().UTC().Add(duration))
}
