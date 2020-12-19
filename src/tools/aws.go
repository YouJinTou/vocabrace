package tools

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/sns"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// BuildSqsURL builds an SQS URL.
func BuildSqsURL(region, accountID, name string) string {
	url := fmt.Sprintf("https://sqs.%s.amazonaws.com/%s/%s", region, accountID, name)
	return url
}

// BuildSnsArn builds an SNS topic arn.
func BuildSnsArn(region, accountID, name string) string {
	url := fmt.Sprintf("arn:aws:sns:%s:%s:%s", region, accountID, name)
	return url
}

// SnsPublish publishes to an SNS topic.
func SnsPublish(topic string, payload interface{}) error {
	sess := session.Must(session.NewSession())
	svc := sns.New(sess)
	b, _ := json.Marshal(payload)
	s := string(b)
	arn := BuildSnsArn(os.Getenv("REGION"), os.Getenv("ACCOUNT_ID"), topic)
	_, err := svc.Publish(&sns.PublishInput{
		Message:  aws.String(s),
		TopicArn: aws.String(arn),
	})
	if err != nil {
		log.Print(err)
	}
	return err
}

// GetItem gets an item from AWS DynamoDB.
func GetItem(table *string, pkName, pkValue string, skName, skValue *string, projection *string) (
	*dynamodb.GetItemOutput, error) {
	sess := session.Must(session.NewSession())
	dynamo := dynamodb.New(sess)
	key := map[string]*dynamodb.AttributeValue{
		pkName: {S: aws.String(pkValue)},
	}
	if skName != nil && skValue != nil {
		key[*skName] = &dynamodb.AttributeValue{
			S: skValue,
		}
	}
	i := &dynamodb.GetItemInput{
		TableName: table,
		Key:       key,
	}
	if projection != nil {
		i.ProjectionExpression = projection
	}
	o, err := dynamo.GetItem(i)
	return o, err
}

// PutItem puts an item in AWS DynamoDB.
func PutItem(table *string, v interface{}) (*dynamodb.PutItemOutput, error) {
	sess := session.Must(session.NewSession())
	dynamo := dynamodb.New(sess)
	item, err := dynamodbattribute.MarshalMap(v)
	if err != nil {
		return nil, err
	}
	o, pErr := dynamo.PutItem(&dynamodb.PutItemInput{
		TableName: table,
		Item:      item,
	})
	return o, pErr
}

// DeleteItem deletes an item from AWS DynamoDB.
func DeleteItem(table *string, pkName, pkValue string, skName, skValue *string) (
	*dynamodb.DeleteItemOutput, error) {
	sess := session.Must(session.NewSession())
	dynamo := dynamodb.New(sess)
	key := map[string]*dynamodb.AttributeValue{
		pkName: {S: aws.String(pkValue)},
	}
	if skName != nil && skValue != nil {
		key[*skName] = &dynamodb.AttributeValue{
			S: skValue,
		}
	}
	o, err := dynamo.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: table,
		Key:       key,
	})
	return o, err
}

// BatchGetItem gets items in batches from AWS DynamoDB.
func BatchGetItem(table *string, partitionKey string, keys []string) (
	*dynamodb.BatchGetItemOutput, error) {
	sess := session.Must(session.NewSession())
	dynamo := dynamodb.New(sess)
	kaa := &dynamodb.KeysAndAttributes{}
	keysMap := []map[string]*dynamodb.AttributeValue{}

	for _, k := range keys {
		keysMap = append(keysMap, map[string]*dynamodb.AttributeValue{
			partitionKey: {S: aws.String(k)}})
	}

	kaa.SetKeys(keysMap)

	o, err := dynamo.BatchGetItem(&dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			*table: kaa,
		},
	})

	return o, err
}

// Table builds a stage-prepended table name.
func Table(name string) *string {
	var stage = os.Getenv("STAGE")
	if stage == "" {
		stage = "dev"
	}
	return aws.String(fmt.Sprintf("%s_%s", stage, name))
}
