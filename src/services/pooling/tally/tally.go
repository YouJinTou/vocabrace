package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/YouJinTou/vocabrace/tools"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func main() {
	lambda.Start(handle)
}

func handle(ctx context.Context, e events.DynamoDBEvent) error {
	sess := session.Must(session.NewSession())
	dynamo := dynamodb.New(sess)
	for _, r := range e.Records {
		division, joined, players := division(r.Change)
		d, _ := tools.GetItem(tools.Table("tallies"), "ID", division, nil, nil, nil)
		shouldPool := willReachCapacity(d, joined)
		var i *dynamodb.UpdateItemInput

		if shouldPool {
			i = poolInput(division, r.Change.NewImage["ID"].String())
		} else if joined {
			i = connectInput(division, r.Change.NewImage["ID"].String(), players)
		} else {
			i = disconnectInput(division, r.Change.OldImage["ID"].String())
		}
		if _, err := dynamo.UpdateItem(i); err != nil {
			log.Print(err)
		}
	}
	return nil
}

func poolInput(division, connectionID string) *dynamodb.UpdateItemInput {
	ue := "REMOVE ConnectionIDs SET LastConnectionID = :cid, ShouldPool = :s"
	i := &dynamodb.UpdateItemInput{
		TableName:        tools.Table("tallies"),
		Key:              map[string]*dynamodb.AttributeValue{"ID": {S: aws.String(division)}},
		UpdateExpression: aws.String(ue),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":cid": &dynamodb.AttributeValue{S: aws.String(connectionID)},
			":s":   &dynamodb.AttributeValue{BOOL: aws.Bool(true)},
		},
	}
	return i
}

func connectInput(division, connectionID string, players int) *dynamodb.UpdateItemInput {
	ue := fmt.Sprintf("%s %s, %s", "ADD ConnectionIDs :c", "SET ShouldPool = :s", "#c = :cap")
	i := &dynamodb.UpdateItemInput{
		TableName:        tools.Table("tallies"),
		Key:              map[string]*dynamodb.AttributeValue{"ID": {S: aws.String(division)}},
		UpdateExpression: aws.String(ue),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":c":   &dynamodb.AttributeValue{SS: []*string{aws.String(connectionID)}},
			":cap": &dynamodb.AttributeValue{N: aws.String(strconv.Itoa(players))},
			":s":   &dynamodb.AttributeValue{BOOL: aws.Bool(false)},
		},
		ExpressionAttributeNames: map[string]*string{"#c": aws.String("Capacity")},
	}
	return i
}

func disconnectInput(division, connectionID string) *dynamodb.UpdateItemInput {
	ue := fmt.Sprintf("%s %s", "DELETE ConnectionIDs :c", "SET ShouldPool = :s")
	i := &dynamodb.UpdateItemInput{
		TableName:        tools.Table("tallies"),
		Key:              map[string]*dynamodb.AttributeValue{"ID": {S: aws.String(division)}},
		UpdateExpression: aws.String(ue),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":c":   &dynamodb.AttributeValue{SS: []*string{aws.String(connectionID)}},
			":s":   &dynamodb.AttributeValue{BOOL: aws.Bool(false)},
			":cid": &dynamodb.AttributeValue{S: aws.String(connectionID)},
		},
		ConditionExpression: aws.String("contains(ConnectionIDs, :cid)"),
	}
	return i
}

func division(r events.DynamoDBStreamRecord) (string, bool, int) {
	_, playerJoined := r.NewImage["LiveUntil"]
	var bucket, game, language string
	var players int

	if playerJoined {
		bucket, game, language, players = extract(r.NewImage)
	} else {
		bucket, game, language, players = extract(r.OldImage)
	}

	division := fmt.Sprintf("%s_%s_%s_%s_%d", os.Getenv("STAGE"), game, bucket, language, players)
	return division, playerJoined, players
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
