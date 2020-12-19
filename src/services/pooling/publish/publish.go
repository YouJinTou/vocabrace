package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/YouJinTou/vocabrace/services/com/state"
	"github.com/YouJinTou/vocabrace/services/com/state/data"
	"github.com/YouJinTou/vocabrace/tools"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type pool struct {
	ID            string
	ConnectionIDs []string
}

type payload struct {
	PoolID string `json:"pid"`
	Body   string `json:"d"`
}

func main() {
	lambda.Start(handle)
}

func handle(_ context.Context, req *events.APIGatewayWebsocketProxyRequest) (
	events.APIGatewayProxyResponse, error) {
	p, dErr := getPayload(req.Body)
	if dErr != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, dErr
	}

	pool, err := getPool(p.PoolID)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	cons, gErr := data.GetConnections(pool.ConnectionIDs)
	if gErr != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, gErr
	}

	aErr := state.OnAction(data.OnActionInput{
		PoolID:          pool.ID,
		Body:            p.Body,
		Connections:     cons,
		Initiator:       req.RequestContext.ConnectionID,
		InitiatorUserID: *cons.UserIDByID(req.RequestContext.ConnectionID),
	})
	if aErr != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, aErr
	}

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func getPayload(body string) (*payload, error) {
	d := &payload{}
	pErr := json.Unmarshal([]byte(body), d)

	if pErr != nil {
		return nil, pErr
	}

	gErr := json.Unmarshal([]byte(d.Body), d)

	if gErr != nil {
		return nil, gErr
	}

	return d, nil
}

func getPool(poolID string) (*pool, error) {
	p, pErr := tools.GetItem(
		tools.Table("pools"),
		"ID",
		poolID,
		nil,
		nil,
		aws.String("ID, ConnectionIDs"))
	if pErr != nil {
		return nil, pErr
	}
	pool := &pool{}
	dynamodbattribute.UnmarshalMap(p.Item, pool)
	return pool, nil
}
