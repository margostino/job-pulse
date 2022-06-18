package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func SanitizeUrl(rawUrl string) string {
	return strings.SplitN(rawUrl, "?", -1)[0]
}

func IndexFrom(url string, value int) string {
	param := strconv.Itoa(value)
	return fmt.Sprintf("%s&start=%s", url, param)
}
