package main

import (
	"context"

	dynamodbpooling "github.com/YouJinTou/vocabrace/pooling/providers/dynamodb"

	ws "github.com/YouJinTou/vocabrace/lambda/pooling"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handle)
}

func handle(_ context.Context, req *events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	c := ws.GetConfig()
	provider := dynamodbpooling.NewDynamoDBProvider()
	connectionIDs, err := provider.GetPeers(req.RequestContext.ConnectionID)

	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, nil
	}

	ws.SendToPeers(connectionIDs, ws.Message{
		Domain:  req.RequestContext.DomainName,
		Stage:   c.Stage,
		Message: req.Body,
	})

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}
