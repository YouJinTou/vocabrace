package tools

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// BuildSqsURL builds an SQS URL.
func BuildSqsURL(region, accountID, name string) string {
	url := fmt.Sprintf("https://sqs.%s.amazonaws.com/%s/%s", region, accountID, name)
	return url
}

// BatchGetItem gets items in batches from AWS DynamoDB.
func BatchGetItem(tableName, partitionKey string, keys []string) (*dynamodb.BatchGetItemOutput, error) {
	sess := session.Must(session.NewSession())
	dynamo := dynamodb.New(sess)
	kaa := &dynamodb.KeysAndAttributes{}
	keysMap := []map[string]*dynamodb.AttributeValue{}

	for _, k := range keys {
		keysMap = append(keysMap, map[string]*dynamodb.AttributeValue{
			partitionKey: {S: aws.String(strings.ToLower(k))}})
	}

	kaa.SetKeys(keysMap)

	o, err := dynamo.BatchGetItem(&dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			tableName: kaa,
		},
	})

	return o, err
}
