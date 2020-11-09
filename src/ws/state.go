package ws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type onStart func(*ReceiverData)
type onAction func(*ReceiverData)

// OnStart executes communication logic at the start of a game.
func OnStart(data *ReceiverData) {
	switch data.Game {
	case "scrabble":
		s := scrabblews{}
		s.OnStart(data)
	default:
		panic(fmt.Sprintf("invalid game %s", data.Game))
	}
}

// OnAction executes communication logic when a player takes an action.
func OnAction(data *ReceiverData) {
	switch data.Game {
	case "scrabble":
		s := scrabblews{}
		s.OnAction(data)
	default:
		panic(fmt.Sprintf("invalid game %s", data.Game))
	}
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
