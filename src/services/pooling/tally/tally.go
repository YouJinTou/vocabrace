package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/YouJinTou/vocabrace/services/com/state/data"

	"github.com/YouJinTou/vocabrace/tools"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type div struct {
	Name     string
	Bucket   string
	Joined   bool
	Players  int
	Language string
	Game     string
}

func main() {
	lambda.Start(handle)
}

func handle(ctx context.Context, e events.DynamoDBEvent) {
	sess := session.Must(session.NewSession())
	dynamo := dynamodb.New(sess)
	for _, r := range e.Records {
		if reconnect(r.Change.NewImage) {
			continue
		}

		div := division(r.Change)
		d, _ := tools.GetItem(tools.Table("tallies"), "ID", div.Name, nil, nil, nil)
		shouldPool := willReachCapacity(d, div.Joined)
		var i *dynamodb.UpdateItemInput

		if shouldPool {
			i = poolInput(div.Name, r.Change.NewImage["ID"].String())
		} else if div.Joined {
			i = connectInput(div, r.Change.NewImage["ID"].String())
		} else {
			i = disconnectInput(div.Name, r.Change.OldImage["ID"].String())
		}

		if _, err := dynamo.UpdateItem(i); err == nil {
			if shouldPool {
				pool(r.Change.NewImage["ID"].String(), d.Item["ConnectionIDs"].SS, div)
			}
		} else {
			log.Print(err)
		}
	}
}

func pool(connectionID string, connectionIDs []*string, division div) {
	cids := []*string{&connectionID}
	cids = append(cids, connectionIDs...)
	payload := struct {
		ConnectionIDs []*string
		Game          string
		Bucket        string
		Language      string
	}{
		cids,
		division.Game,
		division.Bucket,
		division.Language,
	}
	tools.SnsPublish(fmt.Sprintf("%s_pools", os.Getenv("STAGE")), payload)
}

func poolInput(division, connectionID string) *dynamodb.UpdateItemInput {
	i := &dynamodb.UpdateItemInput{
		TableName:        tools.Table("tallies"),
		Key:              map[string]*dynamodb.AttributeValue{"ID": {S: aws.String(division)}},
		UpdateExpression: aws.String("REMOVE ConnectionIDs"),
	}
	return i
}

func connectInput(division div, connectionID string) *dynamodb.UpdateItemInput {
	i := &dynamodb.UpdateItemInput{
		TableName:        tools.Table("tallies"),
		Key:              map[string]*dynamodb.AttributeValue{"ID": {S: aws.String(division.Name)}},
		UpdateExpression: aws.String("ADD ConnectionIDs :c SET #c = :cap"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":c":   &dynamodb.AttributeValue{SS: []*string{aws.String(connectionID)}},
			":cap": &dynamodb.AttributeValue{N: aws.String(strconv.Itoa(division.Players))},
		},
		ExpressionAttributeNames: map[string]*string{"#c": aws.String("Capacity")},
	}
	return i
}

func disconnectInput(division, connectionID string) *dynamodb.UpdateItemInput {
	i := &dynamodb.UpdateItemInput{
		TableName:        tools.Table("tallies"),
		Key:              map[string]*dynamodb.AttributeValue{"ID": {S: aws.String(division)}},
		UpdateExpression: aws.String("DELETE ConnectionIDs :c"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":c":   &dynamodb.AttributeValue{SS: []*string{aws.String(connectionID)}},
			":cid": &dynamodb.AttributeValue{S: aws.String(connectionID)},
		},
		ConditionExpression: aws.String("contains(ConnectionIDs, :cid)"),
	}
	return i
}

func division(r events.DynamoDBStreamRecord) div {
	_, playerJoined := r.NewImage["LiveUntil"]
	var bucket, game, language string
	var players int

	if playerJoined {
		bucket, game, language, players = extract(r.NewImage)
	} else {
		bucket, game, language, players = extract(r.OldImage)
	}

	division := fmt.Sprintf("%s_%s_%s_%s_%d", os.Getenv("STAGE"), game, bucket, language, players)
	return div{
		Name:     division,
		Bucket:   bucket,
		Game:     game,
		Players:  players,
		Language: language,
		Joined:   playerJoined,
	}
}

func extract(m map[string]events.DynamoDBAttributeValue) (string, string, string, int) {
	players, _ := m["Players"].Integer()
	return m["Bucket"].String(), m["Game"].String(), m["Language"].String(), int(players)
}

func willReachCapacity(d *dynamodb.GetItemOutput, joined bool) bool {
	if d.Item == nil || !joined {
		return false
	}

	var waiting = 0
	if val, ok := d.Item["ConnectionIDs"]; ok {
		waiting = len(val.SS)
	}
	capacity, _ := strconv.Atoi(*d.Item["Capacity"].N)
	return waiting+1 >= capacity
}

func reconnect(m map[string]events.DynamoDBAttributeValue) bool {
	pid, _ := m["PoolID"]
	if pid.String() == "" {
		return false
	}

	ID, _ := m["ID"]
	domain, _ := m["Domain"]
	game, _ := m["Game"]
	language, _ := m["Language"]
	userID, _ := m["UserID"]
	connection := data.Connection{
		ID:       ID.String(),
		Domain:   domain.String(),
		Game:     game.String(),
		Language: language.String(),
		UserID:   userID.String(),
	}
	payload := struct {
		Connection data.Connection
		PoolID     string
	}{connection, pid.String()}

	err := tools.SnsPublish(fmt.Sprintf("%s_reconnect", os.Getenv("STAGE")), payload)
	return err == nil
}
