package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/YouJinTou/vocabrace/services/pooling/ws"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type payload struct {
	ConnectionIDs []string
}

func main() {
	lambda.Start(handle)
}

func handle(ctx context.Context, e events.SNSEvent) error {
	for _, r := range e.Records {
		if c, err := getInput(r.SNS.Message, ws.GetConnections); err == nil {
			if sErr := ws.OnStart(c); sErr != nil {
				log.Printf(sErr.Error())
			}
		} else {
			log.Printf(err.Error())
		}
	}
	return nil
}

func getInput(body string, get func([]string) (*ws.Connections, error)) (*ws.Connections, error) {
	p := &payload{}
	json.Unmarshal([]byte(body), p)
	connections, err := get(p.ConnectionIDs)
	return connections, err
}
