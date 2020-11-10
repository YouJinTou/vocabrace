package main

import (
	"context"
	"encoding/json"

	lambdapooling "github.com/YouJinTou/vocabrace/lambda/pooling"

	"github.com/YouJinTou/vocabrace/pooling"
	"github.com/YouJinTou/vocabrace/ws"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handle)
}

func handle(_ context.Context, req *events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	data := getData(req.Body)
	c := lambdapooling.GetConfig()
	pool, err := pooling.GetPool(req.RequestContext.ConnectionID, c.Stage)

	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	ws.OnAction(&ws.ReceiverData{
		Game:          getGame(data),
		Domain:        req.RequestContext.DomainName,
		Stage:         c.Stage,
		PoolID:        pool.ID,
		Body:          data,
		ConnectionIDs: pool.ConnectionIDs,
		Initiator:     req.RequestContext.ConnectionID,
	})

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func getData(body string) string {
	type payload struct {
		Data string
	}
	p := payload{}
	err := json.Unmarshal([]byte(body), &p)

	if err != nil {
		panic(err.Error())
	}

	return p.Data
}

func getGame(body string) string {
	type game struct {
		Game string
	}
	g := game{}
	err := json.Unmarshal([]byte(body), &g)

	if err != nil {
		panic(err.Error())
	}

	return g.Game
}
