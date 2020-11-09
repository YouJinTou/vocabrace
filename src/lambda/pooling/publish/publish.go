package main

import (
	"context"

	lambdapooling "github.com/YouJinTou/vocabrace/lambda/pooling"

	"github.com/YouJinTou/vocabrace/pooling"
	"github.com/YouJinTou/vocabrace/ws"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handle)
}

func handle(_ context.Context, req *events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	c := lambdapooling.GetConfig()
	pool, err := pooling.GetPool(req.RequestContext.ConnectionID, c.Stage)

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
