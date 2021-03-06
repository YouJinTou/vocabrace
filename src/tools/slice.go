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

// ContainsStringPtr checks if a slice of strings contains a given string.
func ContainsStringPtr(s []*string, t string) bool {
	for _, i := range s {
		if i == nil {
			continue
		}
		if *i == t {
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

// ToStringPtrs takes a slice of strings and returns a slice of pointers to strings.
func ToStringPtrs(s []string) []*string {
	result := make([]*string, len(s), len(s))
	for i := 0; i < len(s); i++ {
		result[i] = &s[i]
	}
	return result
}

// FromStringPtrs takes a slice of pointer strings and returns a slice of strings.
func FromStringPtrs(s []*string) []string {
	result := make([]string, len(s), len(s))
	for i := 0; i < len(s); i++ {
		result[i] = *s[i]
	}
	return result
}

// UniqueStrs takes a slice and returns it without the duplicates.
func UniqueStrs(s []string) []string {
	m := make(map[string]bool)
	unique := []string{}
	for _, t := range s {
		if _, ok := m[t]; !ok {
			m[t] = true
			unique = append(unique, t)
		}
	}
	return unique
}
