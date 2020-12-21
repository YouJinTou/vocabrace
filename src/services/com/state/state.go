package state

import (
	"fmt"
	"strconv"

	"github.com/YouJinTou/vocabrace/services/com/state/data"

	wordlines "github.com/YouJinTou/vocabrace/services/com/wordlines"
	"github.com/YouJinTou/vocabrace/services/com/ws"
	"github.com/YouJinTou/vocabrace/tools"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type saveStateInput struct {
	PoolID        string
	ConnectionIDs []string
	Game          interface{}
}

func load(game string) (data.State, error) {
	switch game {
	case "wordlines":
		return wordlines.Controller{}, nil
	default:
		return nil, fmt.Errorf("invalid game %s", game)
	}
}

// OnStart executes communication logic at the start of a game.
func OnStart(i data.OnStartInput) error {
	handler, err := load(i.Connections.Game())
	if err != nil {
		return err
	}

	o, sErr := handler.OnStart(i)
	if sErr != nil {
		return sErr
	}

	ssi := &saveStateInput{
		PoolID:        o.PoolID,
		ConnectionIDs: i.Connections.IDs(),
		Game:          o.Game,
	}
	if err := saveState(ssi); err != nil {
		return err
	}

	ws.SendManyUnique(o.Messages)

	appendHistory(o.PoolID, o.Messages)
	return nil
}

// OnAction executes communication logic when a player takes an action.
func OnAction(i data.OnActionInput) error {
	handler, err := load(i.Connections.Game())
	if err != nil {
		return err
	}

	i.State = loadState(i.PoolID)
	o, aErr := handler.OnAction(i)
	if aErr != nil {
		return aErr
	}

	si := &saveStateInput{
		PoolID:        i.PoolID,
		ConnectionIDs: i.Connections.IDs(),
		Game:          o.Game,
	}
	if err := saveState(si); err != nil {
		return err
	}

	ws.SendManyUnique(o.Messages)

	appendHistory(i.PoolID, o.Messages)
	return nil
}

// OnReconnect executes communication logic during a reconnect.
func OnReconnect(i data.OnReconnectInput) error {
	handler, err := load(i.Connection.Game)
	if err != nil {
		return err
	}

	o, rErr := handler.OnReconnect(i)
	if rErr != nil {
		return err
	}

	ws.Send(o.Message)
	return nil
}

func saveState(i *saveStateInput) error {
	sess := session.Must(session.NewSession())
	dynamo := dynamodb.New(sess)
	m, mErr := dynamodbattribute.MarshalMap(i.Game)

	if mErr != nil {
		return mErr
	}

	_, uErr := dynamo.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: tools.Table("pools"),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {S: aws.String(i.PoolID)},
		},
		UpdateExpression: aws.String(`
			"SET GameState = :s, 
			ConnectionIDs = :c, 
			History = if_not_exists(History, :h), 
			LiveUntil = if_not_exists(LiveUntil, :l)`),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":s": {M: m},
			":c": {SS: tools.ToStringPtrs(i.ConnectionIDs)},
			":h": {L: []*dynamodb.AttributeValue{}},
			":l": {N: aws.String(strconv.Itoa(tools.FutureTimestamp(24 * 3600 * 7)))},
		},
	})

	return uErr
}

func loadState(poolID string) map[string]*dynamodb.AttributeValue {
	i, err := tools.GetItem(tools.Table("pools"), "ID", poolID, nil, nil, nil)
	if err != nil {
		panic(err.Error())
	}

	return i.Item["GameState"].M
}

func appendHistory(poolID string, messages []*ws.Message) {
	sess := session.Must(session.NewSession())
	dynamo := dynamodb.New(sess)
	l, err := dynamodbattribute.MarshalList(messages)
	if err != nil {
		return
	}
	dynamo.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: tools.Table("pools"),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {S: aws.String(poolID)},
		},
		UpdateExpression: aws.String("SET History = list_append(History, :d)"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":d": {L: []*dynamodb.AttributeValue{
				&dynamodb.AttributeValue{L: l},
			}},
		},
	})
}
