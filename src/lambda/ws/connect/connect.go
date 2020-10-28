package main

import (
	"context"

	lambdaws "github.com/YouJinTou/vocabrace/lambda/ws"

	"github.com/YouJinTou/vocabrace/pool"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
)

func main() {
	lambda.Start(handle)
}

func handle(ctx context.Context, req *events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	c := lambdaws.GetConfig()
	p := pool.NewMemcached(c.MemcachedHost, c.MemcachedUsername, c.MemcachedPassword)
	err := p.JoinOrCreate(&pool.Request{
		ConnectionID: req.RequestContext.ConnectionID,
		UserID:       uuid.New().String(),
		PoolLimit:    5})

	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}
