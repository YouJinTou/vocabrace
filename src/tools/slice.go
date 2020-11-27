package tools

import (
	"sort"
	"strings"
)

// ContainsString checks if a slice of strings contains a given string.
func ContainsString(s []string, t string) bool {
	for _, i := range s {
		if i == t {
			return true
		}
	}
	return false
}

// ContainsInt checks if a slice of ints contains a given int.
func ContainsInt(s []int, t int) bool {
	for _, i := range s {
		if i == t {
			return true
		}
	}
	return false
}

// ToLowerStrings to-lowers all strings in a slice.
func ToLowerStrings(s []string) []string {
	result := []string{}
	for _, x := range s {
		result = append(result, strings.ToLower(x))
	}
	return result
}

// SortInts sorts ints in ascending order.
func SortInts(i ...int) []int {
	sort.Ints(i)
	return i
}
