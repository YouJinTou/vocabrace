package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handle)
}

func handle(ctx context.Context, e events.DynamoDBEvent) error {
	// sess := session.Must(session.NewSession())
	// dynamo := dynamodb.New(sess)
	for _, r := range e.Records {
		log.Println(r)
	}
	return nil
}
