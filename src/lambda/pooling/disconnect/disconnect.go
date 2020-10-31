package main

import (
	"context"

	ws "github.com/YouJinTou/vocabrace/lambda/pooling"

	"github.com/YouJinTou/vocabrace/pooling"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handle)
}

func handle(_ context.Context, req events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	c := ws.GetConfig()
	con := pooling.NewMemcachedContext(c.MemcachedHost, c.MemcachedUsername, c.MemcachedPassword)
	pool, err := con.Leave(req.RequestContext.ConnectionID)

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
