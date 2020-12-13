package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/YouJinTou/vocabrace/tools"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/YouJinTou/vocabrace/services/pooling/ws"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type connection struct {
	ID       string
	Domain   string
	Game     string
	Language string
	UserID   string
}

type payload struct {
	ConnectionIDs []*string
	Game          string
	Bucket        string
	Language      string
}

func main() {
	lambda.Start(handle)
}

func handle(ctx context.Context, e events.SQSEvent) error {
	for _, r := range e.Records {
		if i, err := getInput(r.Body, tools.BatchGetItem); err == nil {
			if sErr := ws.OnStart(i); sErr != nil {
				log.Printf(sErr.Error())
			}
		} else {
			log.Printf(err.Error())
		}
	}
	return nil
}

func getInput(
	body string,
	batchGetItem func(*string, string, []string) (*dynamodb.BatchGetItemOutput, error)) (
	*ws.OnStartInput, error) {
	p := &payload{}
	json.Unmarshal([]byte(body), p)
	table := tools.Table("connections")
	o, err := batchGetItem(table, "ID", tools.FromStringPtrs(p.ConnectionIDs))
	if err != nil {
		return nil, err
	}

	users := []*ws.User{}
	input := &ws.OnStartInput{}

	for _, response := range o.Responses[*table] {
		con := &connection{}
		dynamodbattribute.UnmarshalMap(response, con)
		users = append(users, &ws.User{
			ConnectionID: con.ID,
			UserID:       con.UserID,
			Username:     con.UserID,
		})
		input.Language = con.Language
		input.Domain = con.Domain
		input.Game = con.Game
	}

	input.Users = users
	return input, nil
}
