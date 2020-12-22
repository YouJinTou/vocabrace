package state

import (
	"encoding/json"
	"errors"
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
func OnStart(i data.OnStartInput) {
	handler, err := load(i.Connections.Game())
	if err != nil {
		ws.SendManyUnique(startError(&i, "load", err))
		return
	}

	o := handler.OnStart(i)
	if o.Error != nil {
		ws.SendManyUnique(startError(&i, "start", nil))
		return
	}

	ssi := &saveStateInput{
		PoolID:        o.PoolID,
		ConnectionIDs: i.Connections.IDs(),
		Game:          o.Game,
	}
	if err := saveState(ssi); err != nil {
		ws.SendManyUnique(startError(&i, "save", err))
		return
	}

	ws.SendManyUnique(o.Messages)

	appendHistory(o.PoolID, o.Messages)

	updateUserPools(ssi, false)
}

// OnAction executes communication logic when a player takes an action.
func OnAction(i data.OnActionInput) {
	handler, err := load(i.Connections.Game())
	if err != nil {
		ws.Send(actionError(&i, "load", err))
		return
	}

	i.State = loadState(i.PoolID)
	o := handler.OnAction(i)

	if o.Error != nil {
		ws.Send(o.Error)
		return
	}

	si := &saveStateInput{
		PoolID:        i.PoolID,
		ConnectionIDs: i.Connections.IDs(),
		Game:          o.Game,
	}
	if err := saveState(si); err != nil {
		ws.Send(actionError(&i, "save", err))
		return
	}

	ws.SendManyUnique(o.Messages)

	appendHistory(i.PoolID, o.Messages)

	if o.IsOver {
		updateUserPools(si, true)
	}
}

// OnReconnect executes communication logic during a reconnect.
func OnReconnect(i data.OnReconnectInput) error {
	playerSpecificHistory := []string{}
	for _, turn := range i.History {
		var found *ws.Message
		for _, turnMessage := range turn {
			if *turnMessage.UserID == i.Connection.UserID {
				found = turnMessage
			}
		}
		if found == nil {
			return errors.New("could not load history")
		}
		playerSpecificHistory = append(playerSpecificHistory, found.Message)
	}

	payload, _ := json.Marshal(playerSpecificHistory)
	message := &ws.Message{
		ConnectionID: i.Connection.ID,
		Domain:       i.Connection.Domain,
		Message:      string(payload),
	}

	ws.Send(message)
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
			SET GameState = :s, 
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

func updateUserPools(i *saveStateInput, delete bool) {
	table := tools.Table("connections")
	o, _ := tools.BatchGetItem(table, "ID", i.ConnectionIDs)
	for _, r := range o.Responses[*table] {
		if delete {
			tools.DeleteItem(tools.Table("user_pools"), "UserID", *r["UserID"].S, nil, nil)
		} else {
			players, _ := strconv.Atoi(*r["Players"].N)
			v := struct {
				UserID    string
				PoolID    string
				Players   int
				Language  string
				LiveUntil int
			}{*r["UserID"].S, i.PoolID, players, *r["Language"].S, tools.FutureTimestamp(10000)}
			tools.PutItem(tools.Table("user_pools"), v)
		}
	}
}

func startError(i *data.OnStartInput, m string, err error) []*ws.Message {
	messages := []*ws.Message{}
	for _, c := range i.Connections.IDs() {
		messages = append(messages, &ws.Message{
			ConnectionID: c,
			Domain:       i.Connections.Domain(),
			Message:      fmt.Sprintf("%s: %s", m, err.Error()),
			UserID:       i.Connections.UserIDByID(c),
		})
	}
	return messages
}

func actionError(i *data.OnActionInput, m string, err error) *ws.Message {
	return &ws.Message{
		ConnectionID: i.Initiator,
		Domain:       i.Connections.Domain(),
		Message:      fmt.Sprintf("%s: %s", m, err.Error()),
		UserID:       &i.InitiatorUserID,
	}
}
