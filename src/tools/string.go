package tools

// ContainsString checks if a slice of strings contains a given string.
func ContainsString(s []string, t string) bool {
	for _, i := range s {
		if i == t {
			return true
		}
	}
	return false
}
