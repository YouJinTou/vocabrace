package ws

import (
	"fmt"

	"github.com/YouJinTou/vocabrace/tools"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type state interface {
	OnStart(*Connections)
	OnAction(*OnActionInput)
}

func load(game string) (state, error) {
	switch game {
	case "scrabble":
		return scrabblews{
			saveState:      saveState,
			loadState:      loadState,
			send:           Send,
			sendManyUnique: SendManyUnique,
		}, nil
	default:
		return nil, fmt.Errorf("invalid game %s", game)
	}
}

// OnStart executes communication logic at the start of a game.
func OnStart(c *Connections) error {
	handler, err := load(c.Game())
	if err == nil {
		handler.OnStart(c)
	}
	return err
}

// OnActionInput encapsulates data for each turn.
type OnActionInput struct {
	Body            string
	Initiator       string
	InitiatorUserID string
	PoolID          string
	Connections     *Connections
}

// OnAction executes communication logic when a player takes an action.
func OnAction(data *OnActionInput) error {
	handler, err := load(data.Connections.Game())
	if err == nil {
		handler.OnAction(data)
	}
	return err
}

type saveStateInput struct {
	PoolID        string
	ConnectionIDs []string
	V             interface{}
}

func saveState(i *saveStateInput) error {
	sess := session.Must(session.NewSession())
	dynamo := dynamodb.New(sess)
	m, mErr := dynamodbattribute.MarshalMap(i.V)

	if mErr != nil {
		return mErr
	}

	_, uErr := dynamo.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: tools.Table("pools"),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {S: aws.String(i.PoolID)},
		},
		UpdateExpression: aws.String("SET GameState = :s, ConnectionIDs = :c"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":s": {M: m},
			":c": {SS: tools.ToStringPtrs(i.ConnectionIDs)},
		},
	})

	return uErr
}

func loadState(poolID string, v interface{}) {
	i, err := tools.GetItem(tools.Table("pools"), "ID", poolID, nil, nil, nil)
	if err != nil {
		panic(err.Error())
	}

	dynamodbattribute.UnmarshalMap(i.Item["GameState"].M, v)
}
