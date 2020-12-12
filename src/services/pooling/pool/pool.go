package main

import (
	"context"
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

func main() {
	lambda.Start(handle)
}

func handle(ctx context.Context, e events.DynamoDBEvent) error {
	for _, r := range e.Records {
		if !r.Change.NewImage["ShouldPool"].Boolean() {
			continue
		}
		if i, err := getInput(r.Change, tools.BatchGetItem); err == nil {
			ws.OnStart(i)
		} else {
			log.Printf(err.Error())
		}
	}
	return nil
}

func getInput(
	r events.DynamoDBStreamRecord,
	batchGetItem func(*string, string, []string) (*dynamodb.BatchGetItemOutput, error)) (
	*ws.OnStartInput, error) {
	lastConnectionID := r.NewImage["LastConnectionID"].String()
	currentConnectionIDs := r.OldImage["ConnectionIDs"].StringSet()
	connectionIDs := []string{lastConnectionID}
	connectionIDs = append(connectionIDs, currentConnectionIDs...)
	table := tools.Table("connections")
	o, err := batchGetItem(table, "ID", connectionIDs)
	users := []*ws.User{}
	input := &ws.OnStartInput{}

	for _, response := range o.Responses[*table] {
		con := connection{}
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

	return input, err
}
