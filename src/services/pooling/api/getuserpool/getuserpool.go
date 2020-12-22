package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/YouJinTou/vocabrace/tools"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type userPool struct {
	UserID   string
	PoolID   string
	Players  int
	Language string
}

func handler(ctx context.Context, r events.APIGatewayProxyRequest) (
	events.APIGatewayProxyResponse, error) {
	response, err := getUserPool(r)
	setResponseHeaders(&response)
	return response, err
}

func main() {
	lambda.Start(handler)
}

func setResponseHeaders(r *events.APIGatewayProxyResponse) {
	if r.Headers == nil {
		r.Headers = make(map[string]string)
	}
	r.Headers["Access-Control-Allow-Origin"] = "*"
}

func getUserPool(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	o, err := tools.GetItem(
		tools.Table("user_pools"), "UserID", r.PathParameters["userID"], nil, nil, nil)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, err
	}
	if _, ok := o.Item["UserID"]; !ok {
		return events.APIGatewayProxyResponse{StatusCode: 404}, nil
	}
	up := &userPool{}
	dynamodbattribute.UnmarshalMap(o.Item, up)
	b, _ := json.Marshal(up)
	s := string(b)
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: s}, nil
}
