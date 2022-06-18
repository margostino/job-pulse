package utils

import "strings"

func GetOrDefault(index int, list []string) string {
	if len(list) >= index+1 {
		return strings.TrimSpace(list[index])
	}
	return ""
}
