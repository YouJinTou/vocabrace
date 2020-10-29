package main

import (
	"context"

	lambdaws "github.com/YouJinTou/vocabrace/lambda/ws"

	"github.com/YouJinTou/vocabrace/pooling"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handle)
}

func handle(_ context.Context, req events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	c := lambdaws.GetConfig()
	con := pooling.NewMemcachedContext(c.MemcachedHost, c.MemcachedUsername, c.MemcachedPassword)
	err := con.Leave(req.RequestContext.ConnectionID)

	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}
