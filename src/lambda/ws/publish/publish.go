package main

import (
	"context"

	lambdaws "github.com/YouJinTou/vocabrace/lambda/ws"
	"github.com/YouJinTou/vocabrace/pool"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handle)
}

func handle(_ context.Context, req *events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	config := lambdaws.GetPoolConfig()
	p := pool.New(config)
	connectionIDs, err := p.GetPeers(req.RequestContext.ConnectionID)

	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, nil
	}

	for _, c := range connectionIDs {
		lambdaws.Send(c, "Testing.")
	}

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}
