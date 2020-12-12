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
	Game   string `json:"g"`
	PoolID string `json:"pid"`
	Body   string `json:"d"`
}

func main() {
	lambda.Start(handle)
}

func handle(_ context.Context, req *events.APIGatewayWebsocketProxyRequest) (
	events.APIGatewayProxyResponse, error) {
	userID, uErr := getInitiatorUserID(req.RequestContext.ConnectionID)
	if uErr != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, uErr
	}

	data, dErr := getData(req.Body)
	if dErr != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, dErr
	}

	pool, err := getPool(data.PoolID)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	ws.OnAction(&ws.OnActionInput{
		Game:            data.Game,
		Domain:          req.RequestContext.DomainName,
		PoolID:          pool.ID,
		Body:            data.Body,
		ConnectionIDs:   pool.ConnectionIDs,
		Initiator:       req.RequestContext.ConnectionID,
		InitiatorUserID: *userID,
	})

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

func getInitiatorUserID(connectionID string) (*string, error) {
	o, err := tools.GetItem(
		tools.Table("connections"), "ID", connectionID, nil, nil, aws.String("UserID"))
	if err != nil || o.Item == nil {
		return nil, err
	}
	return o.Item["UserID"].S, nil
}
