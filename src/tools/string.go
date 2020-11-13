package tools

import "strings"

// ContainsString checks if a slice of strings contains a given string.
func ContainsString(s []string, t string) bool {
	for _, i := range s {
		if strings.ToLower(i) == strings.ToLower(t) {
			return true
		}
	}
	return false
}
