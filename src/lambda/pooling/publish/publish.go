package main

import (
	"context"

	ws "github.com/YouJinTou/vocabrace/lambda/pooling"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handle)
}

func handle(_ context.Context, req *events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	c := ws.GetConfig()
	pool, err := ws.GetPool(req.RequestContext.ConnectionID, c)

	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, err
	}

	ws.SendMany(pool.ConnectionIDs, ws.Message{
		Domain:  req.RequestContext.DomainName,
		Stage:   c.Stage,
		Message: req.Body,
	})

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}
