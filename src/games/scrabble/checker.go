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
	IsValidWord(language, w string) bool
}

// DynamoChecker implements WordChecker.
type DynamoChecker struct{}

// NewDynamoChecker returns a DynamoDB checker.
func NewDynamoChecker() WordChecker {
	return DynamoChecker{}
}

// IsValidWord checks if a given word is valid for the target language.
func (dc DynamoChecker) IsValidWord(language, w string) bool {
	if len(w) == 0 {
		return false
	}

	sess := session.Must(session.NewSession())
	dynamo := dynamodb.New(sess)
	_, err := dynamo.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(fmt.Sprintf("scrabble_%s", strings.ToLower(language))),
		Key: map[string]*dynamodb.AttributeValue{
			"Word": {S: aws.String(strings.ToLower(w))},
		},
	})

	return err == nil
}
