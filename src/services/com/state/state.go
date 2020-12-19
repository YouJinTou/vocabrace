package state

import (
	"fmt"

	sd "github.com/YouJinTou/vocabrace/services/com/state/data"

	"github.com/YouJinTou/vocabrace/services/com/state/ws"
	wordlines "github.com/YouJinTou/vocabrace/services/com/wordlines"
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

func load(game string) (sd.State, error) {
	switch game {
	case "wordlines":
		return wordlines.Controller{}, nil
	default:
		return nil, fmt.Errorf("invalid game %s", game)
	}
}

// OnStart executes communication logic at the start of a game.
func OnStart(i sd.OnStartInput) error {
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
	return nil
}

// OnAction executes communication logic when a player takes an action.
func OnAction(data sd.OnActionInput) error {
	handler, err := load(data.Connections.Game())
	if err != nil {
		return err
	}

	data.State = loadState(data.PoolID)
	o, aErr := handler.OnAction(data)
	if aErr != nil {
		return aErr
	}

	i := &saveStateInput{
		PoolID:        data.PoolID,
		ConnectionIDs: data.Connections.IDs(),
		Game:          o.Game,
	}
	if err := saveState(i); err != nil {
		return err
	}

	ws.SendManyUnique(o.Messages)
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
		UpdateExpression: aws.String("SET GameState = :s, ConnectionIDs = :c"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":s": {M: m},
			":c": {SS: tools.ToStringPtrs(i.ConnectionIDs)},
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
