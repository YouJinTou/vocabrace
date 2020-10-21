package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Item The DynamoDB item for the 'connections' table
type Item struct {
	ConnectionID string
	Timestamp    int64
}

func handle(ctx context.Context, req events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	session := session.Must(session.NewSession())
	dynamo := dynamodb.New(session)
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

	endpoint := fmt.Sprintf(
		"https://%s.execute-api.%s.amazonaws.com/%s",
		req.RequestContext.DomainName,
		"eu-central-1",
		req.RequestContext.Stage)
	apiClient := apigatewaymanagementapi.New(session, aws.NewConfig().WithEndpoint(endpoint))
	connectionInput := apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: aws.String(req.RequestContext.ConnectionID),
		Data:         []byte("{ poolID: \"123\" }"),
	}
	request, _ := apiClient.PostToConnectionRequest(&connectionInput)
	postError := request.Send()

	if postError != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, postError
	}

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func main() {
	lambda.Start(handle)
}
