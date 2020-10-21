package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Item The DynamoDB item for the 'connections' table
type Item struct {
	ConnectionID string
	Timestamp    int64
}

func handle(_ context.Context, req events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	dynamo := dynamodb.New(session.New())
	item := Item{
		ConnectionID: req.RequestContext.ConnectionID,
		Timestamp:    req.RequestContext.ConnectedAt,
	}
	marshalled, err := dynamodbattribute.MarshalMap(item)

	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	input := dynamodb.PutItemInput{
		Item:      marshalled,
		TableName: aws.String("connections"),
	}
	_, putError := dynamo.PutItem(&input)

	if putError != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, putError
	}

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func main() {
	lambda.Start(handle)
}
