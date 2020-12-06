package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

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

func main() {
	lambda.Start(handle)
}

func handle(_ context.Context, req *events.APIGatewayWebsocketProxyRequest) (
	events.APIGatewayProxyResponse, error) {
	data, game := getData(req.Body)
	pool, err := getPool(req.RequestContext.ConnectionID)

	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	ws.OnAction(&ws.ReceiverData{
		Game:          game,
		Domain:        req.RequestContext.DomainName,
		Stage:         os.Getenv("STAGE"),
		PoolID:        pool.ID,
		Body:          data,
		ConnectionIDs: pool.ConnectionIDs,
		Initiator:     req.RequestContext.ConnectionID,
	})

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func getData(body string) (string, string) {
	type payload struct {
		Data string `json:"d"`
	}
	p := payload{}
	pErr := json.Unmarshal([]byte(body), &p)

	if pErr != nil {
		panic(pErr.Error())
	}

	type game struct {
		Game string `json:"g"`
	}
	g := game{}
	gErr := json.Unmarshal([]byte(p.Data), &g)

	if gErr != nil {
		panic(gErr.Error())
	}

	return p.Data, g.Game
}

func getPool(connectionID string) (*pool, error) {
	con, cErr := tools.GetItem(
		fmt.Sprintf("%s_connections", os.Getenv("STAGE")), "ID", connectionID, nil, nil, nil)
	if cErr != nil {
		return nil, cErr
	}
	p, pErr := tools.GetItem(fmt.Sprintf(
		"%s_pools", os.Getenv("STAGE")),
		"ID",
		*con.Item["PoolID"].S,
		nil,
		nil,
		aws.String("ID, ConnectionIDs"))
	if pErr != nil {
		return nil, pErr
	}
	pool := pool{}
	dynamodbattribute.UnmarshalMap(p.Item, &pool)
	return &pool, nil
}
