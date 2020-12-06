package tools

import (
	"time"
)

// FutureTimestamp takes the current UNIX epoch and adds a certain number of seconds to it.
func FutureTimestamp(add int) int {
	return int(time.Now().Unix()) + add
}
