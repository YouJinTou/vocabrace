package core

// SliceRemoveString removes a target string.
func SliceRemoveString(strings []string, key string) []string {
	for i, curr := range strings {
		if curr == key {
			strings = append(strings[:i], strings[i+1:]...)
		}
	}

	return strings
}
