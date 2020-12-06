package main

import (
	"context"

	"github.com/YouJinTou/vocabrace/services/pooling/ws"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handle)
}

func handle(ctx context.Context, req *events.APIGatewayWebsocketProxyRequest) (
	events.APIGatewayProxyResponse, error) {
	p := pooler{OnStart: ws.OnStart}
	switch req.RequestContext.RouteKey {
	case "$connect":
		o, err := p.joinWaitlist(req.RequestContext.ConnectionID, req.QueryStringParameters)

		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500}, err
		}

		p.onWaitlistFull(o, req)
	case "$disconnect":
		err := p.leaveWaitlist(req.RequestContext.ConnectionID, req.QueryStringParameters)

		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500}, err
		}
	}
	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}
