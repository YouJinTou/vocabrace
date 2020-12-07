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
	OnStart(*OnStartInput)
	OnAction(*OnActionInput)
}

func load(game string) state {
	switch game {
	case "scrabble":
		return scrabblews{
			saveState:      saveState,
			loadState:      loadState,
			send:           Send,
			sendManyUnique: SendManyUnique,
		}
	default:
		panic(fmt.Sprintf("invalid game %s", game))
	}
}

// User ecanpsulates user data.
type User struct {
	ConnectionID string
	UserID       string
	Username     string
}

func userByID(users []*User, ID string) *User {
	for _, u := range users {
		if u.UserID == ID {
			return u
		}
	}
	return nil
}

// OnStartInput encapsulates data for the start state.
type OnStartInput struct {
	Users    []*User
	Language string
	PoolID   string
	Domain   string
	Game     string
}

// OnStart executes communication logic at the start of a game.
func OnStart(data *OnStartInput) {
	load(data.Game).OnStart(data)
}

// OnActionInput encapsulates data for each turn.
type OnActionInput struct {
	Body          string
	Initiator     string
	Domain        string
	Game          string
	PoolID        string
	ConnectionIDs []string
}

func (i *OnActionInput) otherConnections() []string {
	connections := []string{}
	for _, cid := range i.ConnectionIDs {
		if cid != i.Initiator {
			connections = append(connections, cid)
		}
	}
	return connections
}

// OnAction executes communication logic when a player takes an action.
func OnAction(data *OnActionInput) {
	load(data.Game).OnAction(data)
}

func saveState(poolID string, v interface{}) error {
	sess := session.Must(session.NewSession())
	dynamo := dynamodb.New(sess)
	m, mErr := dynamodbattribute.MarshalMap(v)

	if mErr != nil {
		return mErr
	}

	_, uErr := dynamo.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: tools.Table("pools"),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {S: aws.String(poolID)},
		},
		UpdateExpression:          aws.String("SET GameState = :s"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{":s": {M: m}},
	})

	return uErr
}

func loadState(poolID string, v interface{}) {
	sess := session.Must(session.NewSession())
	dynamo := dynamodb.New(sess)
	i, err := dynamo.GetItem(&dynamodb.GetItemInput{
		TableName: tools.Table("pools"),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {S: aws.String(poolID)},
		},
	})

	if err != nil {
		panic(err.Error())
	}

	dynamodbattribute.UnmarshalMap(i.Item["GameState"].M, v)
}
