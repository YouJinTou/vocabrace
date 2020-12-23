package wordlines

import (
	"fmt"
	"strings"

	"github.com/YouJinTou/vocabrace/tools"
	"github.com/aws/aws-sdk-go/aws"
)

// WordChecker performs validation checks on words.
type WordChecker interface {
	ValidateWords(language string, words []string) ([]string, error)
}

// DynamoChecker implements WordChecker.
type DynamoChecker struct{}

// NewDynamoChecker returns a DynamoDB checker.
func NewDynamoChecker() WordChecker {
	return DynamoChecker{}
}

// ValidateWords checks if the given words are valid for the target language.
func (dc DynamoChecker) ValidateWords(language string, words []string) ([]string, error) {
	table := fmt.Sprintf("wordlines_%s", strings.ToLower(language))
	o, err := tools.BatchGetItem(&table, "Word", words)
	notFound := []string{}
	for _, w := range words {
		found := false
		for _, tables := range o.Responses {
			for _, kv := range tables {
				if w == *kv["Word"].S {
					found = true
					break
				}
			}
		}
		if !found {
			notFound = append(notFound, w)
		}
	}

	saveMissing(notFound)

	return notFound, err
}

func saveMissing(notFound []string) {
	if len(notFound) == 0 {
		return
	}
	type missing struct {
		Word      string
		LiveUntil int
	}
	for _, w := range notFound {
		item := missing{Word: w, LiveUntil: tools.FutureTimestamp(170000)}
		tools.PutItem(aws.String("wordlines_missing_words"), item)
	}
}
