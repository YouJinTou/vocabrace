package tools

import (
	"strconv"
	"time"
)

// FutureTimestamp takes the current UNIX epoch and adds a certain number of seconds to it.
func FutureTimestamp(add int) int {
	return int(time.Now().Unix()) + add
}

// FutureTimestampStr takes the current UNIX epoch, adds a certain number of seconds to it, and stringifies it.
func FutureTimestampStr(add int) string {
	return strconv.Itoa(FutureTimestamp(add))
}
