package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/YouJinTou/vocabrace/services/com/state/data"

	"github.com/YouJinTou/vocabrace/services/com/state"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type payload struct {
	ConnectionIDs []string
}

func main() {
	lambda.Start(handle)
}

func handle(ctx context.Context, e events.SNSEvent) {
	for _, r := range e.Records {
		if c, err := getInput(r.SNS.Message, data.GetConnections); err == nil {
			if sErr := state.OnStart(c); sErr != nil {
				log.Printf(sErr.Error())
			}
		} else {
			log.Printf(err.Error())
		}
	}
}

func getInput(body string, get func([]string) (*data.Connections, error)) (data.OnStartInput, error) {
	p := &payload{}
	json.Unmarshal([]byte(body), p)
	connections, err := get(p.ConnectionIDs)
	return data.OnStartInput{Connections: connections}, err
}
