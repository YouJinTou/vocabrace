package ws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type state interface {
	OnStart(*ReceiverData) PoolID
	OnAction(*ReceiverData)
}

func load(data *ReceiverData) state {
	switch data.Game {
	case "scrabble":
		return scrabblews{}
	default:
		panic(fmt.Sprintf("invalid game %s", data.Game))
	}
}

// PoolID is the pool ID.
type PoolID string

// OnStart executes communication logic at the start of a game.
func OnStart(data *ReceiverData) PoolID {
	return load(data).OnStart(data)
}

// OnAction executes communication logic when a player takes an action.
func OnAction(data *ReceiverData) {
	load(data).OnAction(data)
}

func saveState(data *ReceiverData, v interface{}) error {
	sess := session.Must(session.NewSession())
	dynamo := dynamodb.New(sess)
	m, mErr := dynamodbattribute.MarshalMap(v)

	if mErr != nil {
		return mErr
	}

	_, uErr := dynamo.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String(fmt.Sprintf("%s_pools", data.Stage)),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {S: aws.String(data.PoolID)},
		},
		UpdateExpression:          aws.String("SET GameState = :s"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{":s": {M: m}},
	})

	return uErr
}

func loadState(data *ReceiverData, v interface{}) {
	sess := session.Must(session.NewSession())
	dynamo := dynamodb.New(sess)
	i, err := dynamo.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(fmt.Sprintf("%s_pools", data.Stage)),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {S: aws.String(data.PoolID)},
		},
	})

	if err != nil {
		panic(err.Error())
	}

	dynamodbattribute.UnmarshalMap(i.Item["GameState"].M, v)
}
