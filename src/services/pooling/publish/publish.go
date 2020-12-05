package main

import (
	"context"
	"encoding/json"

	"github.com/YouJinTou/vocabrace/services/pooling"
	"github.com/YouJinTou/vocabrace/services/pooling/ws"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handle)
}

func handle(_ context.Context, req *events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	data, game := getData(req.Body)
	c := pooling.GetConfig()
	pool, err := pooling.GetPool(req.RequestContext.ConnectionID, c.Stage)

	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	ws.OnAction(&ws.ReceiverData{
		Game:          game,
		Domain:        req.RequestContext.DomainName,
		Stage:         c.Stage,
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
