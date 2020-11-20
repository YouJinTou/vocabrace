package scrabble

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// WordChecker performs validation checks on words.
type WordChecker interface {
	ValidateWords(language string, words []string) error
}

// DynamoChecker implements WordChecker.
type DynamoChecker struct{}

// NewDynamoChecker returns a DynamoDB checker.
func NewDynamoChecker() WordChecker {
	return DynamoChecker{}
}

// ValidateWords checks if the given words are valid for the target language.
func (dc DynamoChecker) ValidateWords(language string, words []string) error {
	sess := session.Must(session.NewSession())
	dynamo := dynamodb.New(sess)
	kaa := &dynamodb.KeysAndAttributes{}
	keys := []map[string]*dynamodb.AttributeValue{}

	for _, w := range words {
		keys = append(keys, map[string]*dynamodb.AttributeValue{
			"Word": {S: aws.String(strings.ToLower(w))}})
	}

	kaa.SetKeys(keys)

	_, err := dynamo.BatchGetItem(&dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			fmt.Sprintf("scrabble_%s", strings.ToLower(language)): kaa,
		},
	})

	return err
}
