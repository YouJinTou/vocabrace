package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/YouJinTou/vocabrace/pooling"

	ws "github.com/YouJinTou/vocabrace/lambda/pooling"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func main() {
	lambda.Start(handle)
}

func handle(_ context.Context, req *events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	c := ws.GetConfig()
	pool, err := getPool(req.RequestContext.ConnectionID, c)

	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, nil
	}

	ws.SendMany(pool.ConnectionIDs, ws.Message{
		Domain:  req.RequestContext.DomainName,
		Stage:   c.Stage,
		Message: req.Body,
	})

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func getPool(connectionID string, c *ws.Config) (*pooling.Pool, error) {
	sess := session.Must(session.NewSession())
	dynamo := dynamodb.New(sess)
	o, cErr := dynamo.GetItem(&dynamodb.GetItemInput{
		TableName:            aws.String(fmt.Sprintf("%s_connections", c.Stage)),
		Key:                  map[string]*dynamodb.AttributeValue{"ID": {S: aws.String(connectionID)}},
		ProjectionExpression: aws.String("PoolID"),
	})

	if o.Item == nil {
		return nil, cErr
	}

	i, pErr := dynamo.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(fmt.Sprintf("%s_pools", c.Stage)),
		Key:       map[string]*dynamodb.AttributeValue{"ID": {S: o.Item["PoolID"].S}},
	})
	pool := pooling.Pool{}
	dynamodbattribute.UnmarshalMap(i.Item, &pool)

	return &pool, pErr
}
