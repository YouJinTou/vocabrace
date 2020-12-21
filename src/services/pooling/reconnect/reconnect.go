package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/YouJinTou/vocabrace/services/com/ws"

	"github.com/YouJinTou/vocabrace/services/com/state"

	"github.com/YouJinTou/vocabrace/services/com/state/data"

	"github.com/YouJinTou/vocabrace/tools"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type payload struct {
	Connection data.Connection
	PoolID     string
}

func main() {
	lambda.Start(handle)
}

func handle(ctx context.Context, e events.SNSEvent) {
	for _, r := range e.Records {
		p := unmarshalPayload(r.SNS.Message)
		result, err := updatePool(p)
		if err == nil {
			state.OnReconnect(data.OnReconnectInput{
				Connection: p.Connection,
				History:    toHistory(result),
			})
		} else {
			log.Printf(err.Error())
		}
	}
}

func unmarshalPayload(body string) payload {
	p := &payload{}
	json.Unmarshal([]byte(body), p)
	return *p
}

func updatePool(p payload) ([]*dynamodb.AttributeValue, error) {
	sess := session.Must(session.NewSession())
	dynamo := dynamodb.New(sess)
	o, err := dynamo.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: tools.Table("pools"),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {S: aws.String(p.PoolID)},
		},
		UpdateExpression: aws.String("ADD ConnectionIDs :c"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":c": {SS: []*string{&p.Connection.ID}},
		},
		ReturnValues: aws.String("ALL_NEW"),
	})
	return o.Attributes["History"].L, err
}

func toHistory(list []*dynamodb.AttributeValue) [][]*ws.Message {
	history := &[][]*ws.Message{}
	dynamodbattribute.UnmarshalList(list, history)
	return *history
}
