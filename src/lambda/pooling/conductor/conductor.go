package main

import (
	"context"
	"encoding/json"
	"fmt"

	ws "github.com/YouJinTou/vocabrace/lambda/pooling"
	"github.com/YouJinTou/vocabrace/pooling"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handle)
}

func handle(ctx context.Context, sqsEvent events.SQSEvent) error {
	c := ws.GetConfig()
	con := pooling.NewMemcachedContext(c.MemcachedHost, c.MemcachedUsername, c.MemcachedPassword)

	for _, message := range sqsEvent.Records {
		payload := ws.PoolPayload{}

		json.Unmarshal([]byte(message.Body), &payload)

		pool, err := con.GetPool(payload.PoolID)

		if err != nil {
			fmt.Println(fmt.Sprintf("Pool %s not found.", payload.PoolID))

			continue
		}

		ws.SendToPeers(pool.ConnectionIDs, ws.Message{
			Domain:  payload.Domain,
			Stage:   c.Stage,
			Message: "GAME MAP",
		})
	}

	return nil
}
