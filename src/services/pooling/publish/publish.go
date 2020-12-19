package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/YouJinTou/vocabrace/tools"

	"github.com/YouJinTou/vocabrace/services/pooling/ws"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type pool struct {
	ID            string
	ConnectionIDs []string
}

type data struct {
	PoolID string `json:"pid"`
	Body   string `json:"d"`
}

func main() {
	lambda.Start(handle)
}

func handle(_ context.Context, req *events.APIGatewayWebsocketProxyRequest) (
	events.APIGatewayProxyResponse, error) {
	data, dErr := getData(req.Body)
	if dErr != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, dErr
	}

	pool, err := getPool(data.PoolID)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	cons, gErr := ws.GetConnections(pool.ConnectionIDs)
	if gErr != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, gErr
	}

	aErr := ws.OnAction(&ws.OnActionInput{
		PoolID:          pool.ID,
		Body:            data.Body,
		Connections:     cons,
		Initiator:       req.RequestContext.ConnectionID,
		InitiatorUserID: *cons.UserIDByID(req.RequestContext.ConnectionID),
	})
	if aErr != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, aErr
	}

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func getData(body string) (*data, error) {
	d := &data{}
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
