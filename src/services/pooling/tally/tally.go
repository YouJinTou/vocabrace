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
		if err := tally(dynamo.GetItem, dynamo.UpdateItem, r); err != nil {
			log.Print(err)
		}
	}
	return nil
}

func tally(
	getItem func(*dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error),
	updateItem func(*dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error),
	r events.DynamoDBEventRecord) error {
	division, sign, players := division(r.Change)
	shouldPool := willReachCapacity(division, sign, getItem)
	var ue, w string
	if shouldPool {
		w = "0"
		ue = fmt.Sprintf("SET Waiting = :w + :z, #c = :c, ShouldPool = :s")
	} else {
		w = "1"
		ue = fmt.Sprintf(
			"SET Waiting = if_not_exists(Waiting, :z) %s :w, #c = :c, ShouldPool = :s", sign)
	}
	fmt.Println(ue)
	_, err := updateItem(&dynamodb.UpdateItemInput{
		TableName:        tools.Table("tallies"),
		Key:              map[string]*dynamodb.AttributeValue{"ID": {S: aws.String(division)}},
		UpdateExpression: aws.String(ue),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":w": &dynamodb.AttributeValue{N: aws.String(w)},
			":c": &dynamodb.AttributeValue{N: aws.String(strconv.Itoa(players))},
			":z": &dynamodb.AttributeValue{N: aws.String("0")},
			":s": &dynamodb.AttributeValue{BOOL: aws.Bool(shouldPool)},
		},
		ExpressionAttributeNames: map[string]*string{"#c": aws.String("Capacity")},
	})
	return err
}

func division(r events.DynamoDBStreamRecord) (string, string, int) {
	_, playerJoined := r.NewImage["LiveUntil"]
	var bucket, game, language string
	var players int
	var sign string

	if playerJoined {
		bucket, game, language, players = extract(r.NewImage)
		sign = "+"
	} else {
		bucket, game, language, players = extract(r.OldImage)
		sign = "-"
	}

	division := fmt.Sprintf("%s_%s_%s_%s_%d", os.Getenv("STAGE"), game, bucket, language, players)
	return division, sign, players
}

func extract(m map[string]events.DynamoDBAttributeValue) (string, string, string, int) {
	players, _ := m["Players"].Integer()
	return m["Bucket"].String(), m["Game"].String(), m["Language"].String(), int(players)
}

func willReachCapacity(
	division, sign string,
	getItem func(*dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)) bool {
	if sign == "-" {
		return false
	}

	d, _ := getItem(&dynamodb.GetItemInput{
		TableName: tools.Table("tallies"),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": &dynamodb.AttributeValue{S: aws.String(division)},
		},
	})

	if d.Item == nil {
		return false
	}

	waiting, _ := strconv.Atoi(*d.Item["Waiting"].N)
	capacity, _ := strconv.Atoi(*d.Item["Capacity"].N)
	return waiting+1 >= capacity
}
