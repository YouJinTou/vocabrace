package main

import (
	"context"

	"github.com/YouJinTou/vocabrace/tools"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func main() {
	lambda.Start(handle)
}

func handle(ctx context.Context, req *events.APIGatewayWebsocketProxyRequest) (
	events.APIGatewayProxyResponse, error) {
	disconnect(tools.DeleteItem, req)
	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func disconnect(
	deleteItem func(*string, string, string, *string, *string) (*dynamodb.DeleteItemOutput, error),
	r *events.APIGatewayWebsocketProxyRequest) {
	deleteItem(tools.Table("connections"), "ID", r.RequestContext.ConnectionID, nil, nil)
}
