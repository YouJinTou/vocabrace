package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type connection struct {
	ID     string
	PoolID string
}

// Leave removes a connection from the pool
func (dpp DynamoDBPoolingProvider) Leave(connectionID string) (*Pool, error) {
	stage := "dev"
	bucket := "novice"
	dpp.dynamo().UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String(fmt.Sprintf("%s_buckets", stage)),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {S: aws.String(bucket)},
		},
		UpdateExpression: aws.String(""),
	})
}

func (dpp DynamoDBPoolingProvider) getPoolID(connectionID string) (string, error) {
	result, err := dpp.dynamo().DeleteItem(&dynamodb.DeleteItemInput{
		TableName:    aws.String(fmt.Sprintf("%s_connections", "dev")),
		Key:          map[string]*dynamodb.AttributeValue{"ID": {S: aws.String(connectionID)}},
		ReturnValues: aws.String("ALL_OLD"),
	})
	c := connection{}
	dynamodbattribute.UnmarshalMap(result.Attributes, &c)

	return c.PoolID, err
}
