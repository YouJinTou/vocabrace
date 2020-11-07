package main

import (
	"context"
	"fmt"

	ws "github.com/YouJinTou/vocabrace/lambda/pooling"
	"github.com/YouJinTou/vocabrace/tools"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func main() {
	lambda.Start(handle)
}

func handle(_ context.Context, req events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	markDisconnection(req.RequestContext.ConnectionID)

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func markDisconnection(ID string) {
	c := ws.GetConfig()
	sess := session.Must(session.NewSession())
	dynamo := dynamodb.New(sess)
	dynamo.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(fmt.Sprintf("%s_disconnections", c.Stage)),
		Item: map[string]*dynamodb.AttributeValue{
			"ID":        {S: aws.String(ID)},
			"LiveUntil": {N: aws.String(tools.FutureTimestampStr(7200))},
		},
	})
}
