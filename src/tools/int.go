package tools

// ContainsInt checks if a slice of ints contains a given int.
func ContainsInt(s []int, t int) bool {
	for _, i := range s {
		if t == i {
			return true
		}
	}
	return false
}
