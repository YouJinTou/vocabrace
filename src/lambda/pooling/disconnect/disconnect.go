package main

import (
	"context"

	"github.com/YouJinTou/vocabrace/pooling"

	dynamodbpooling "github.com/YouJinTou/vocabrace/pooling/providers/dynamodb"

	ws "github.com/YouJinTou/vocabrace/lambda/pooling"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handle)
}

func handle(_ context.Context, req events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	c := ws.GetConfig()
	provider := dynamodbpooling.NewDynamoDBProvider(c.Stage)
	pool, err := provider.Leave(&pooling.LeaveInput{
		ConnectionID: req.RequestContext.ConnectionID,
		Bucket:       pooling.Novice,
	})

	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	ws.SendToPeers(pool.ConnectionIDs, ws.Message{
		Domain:  req.RequestContext.DomainName,
		Stage:   c.Stage,
		Message: "Client has left.",
	})

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}
