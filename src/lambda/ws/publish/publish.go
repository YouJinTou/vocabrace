package main

import (
	"context"
	"fmt"
	"sync"

	lambdaws "github.com/YouJinTou/vocabrace/lambda/ws"
	"github.com/YouJinTou/vocabrace/pool"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handle)
}

func handle(_ context.Context, req *events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	c := lambdaws.GetConfig()
	p := pool.NewMemcached(c.MemcachedHost, c.MemcachedUsername, c.MemcachedPassword)
	connectionIDs, err := p.GetPeers(req.RequestContext.ConnectionID)

	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, nil
	}

	var wg sync.WaitGroup

	for _, cid := range connectionIDs {
		wg.Add(1)

		go send(&wg, req.RequestContext.DomainName, c.Stage, cid, req.Body)
	}

	wg.Wait()

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func send(wg *sync.WaitGroup, domain, stage, connectionID, body string) {
	defer wg.Done()

	if _, err := lambdaws.Send(domain, stage, connectionID, body); err != nil {
		fmt.Println(err.Error())
	}
}
