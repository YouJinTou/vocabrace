package wordlines

import (
	"fmt"
	"strings"

	"github.com/YouJinTou/vocabrace/tools"
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

	return notFound, err
}
