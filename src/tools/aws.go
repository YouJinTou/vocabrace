package tools

import (
	"fmt"
)

// BuildSqsURL builds an SQS URL.
func BuildSqsURL(region, accountID, name string) string {
	url := fmt.Sprintf("https://sqs.%s.amazonaws.com/%s/%s", region, accountID, name)
	return url
}
