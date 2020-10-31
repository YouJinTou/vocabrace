package main

import (
	"context"

	"github.com/YouJinTou/vocabrace/pooling"

	ws "github.com/YouJinTou/vocabrace/lambda/pooling"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
)

func main() {
	lambda.Start(handle)
}

func handle(ctx context.Context, req *events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	c := ws.GetConfig()
	con := pooling.NewMemcachedContext(c.MemcachedHost, c.MemcachedUsername, c.MemcachedPassword)
	_, err := con.JoinOrCreate(&pooling.Request{
		ConnectionID: req.RequestContext.ConnectionID,
		UserID:       uuid.New().String(),
		PoolLimit:    3})

	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, nil
	}

	peers, peersErr := con.GetPeers(req.RequestContext.ConnectionID)

	if peersErr != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, nil
	}

	ws.SendToPeers(peers, ws.Message{
		Domain:  req.RequestContext.DomainName,
		Stage:   c.Stage,
		Message: "Just connected.",
	})

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}