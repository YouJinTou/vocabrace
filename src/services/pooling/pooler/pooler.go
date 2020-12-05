package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/YouJinTou/vocabrace/services/pooling/ws"

	"github.com/YouJinTou/vocabrace/services/pooling"

	"github.com/YouJinTou/vocabrace/tools"
	"github.com/aws/aws-sdk-go/service/sqs"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
)

type pool struct {
	ID           string
	Users        []*ws.User
	Connections  []*connection
	Domain       string
	Bucket       string
	Game         string
	Stage        string
	PlayersCount int
	Language     string
}

type connection struct {
	ID             string
	ReceiptHandle  string
	IsDisconnected bool
}

func (p *pool) ConnectionIDs() []string {
	IDs := []string{}
	for _, c := range p.Connections {
		IDs = append(IDs, c.ID)
	}
	return IDs
}

func (p *pool) ConnectionIDsPtr() []*string {
	IDs := []*string{}
	for _, c := range p.Connections {
		IDs = append(IDs, &c.ID)
	}
	return IDs
}

func main() {
	lambda.Start(handle)
}

func handle(ctx context.Context, event events.SQSEvent) error {
	c := pooling.GetConfig()

	pool, err := createPool(event.Records, c)
	if err != nil {
		return err
	}

	if err := handleDisconnections(event, pool, c); err != nil {
		return err
	}

	setPool(pool, c)

	ws.OnStart(&ws.ReceiverData{
		Users:         pool.Users,
		ConnectionIDs: pool.ConnectionIDs(),
		Domain:        pool.Domain,
		Stage:         pool.Stage,
		Game:          pool.Game,
		PoolID:        pool.ID,
		Language:      pool.Language,
	})

	return nil
}

func createPool(batch []events.SQSMessage, c *pooling.Config) (*pool, error) {
	p := pool{Connections: []*connection{}, Stage: c.Stage}
	users := []*ws.User{}
	for _, message := range batch {
		payload := pooling.PoolerPayload{}
		json.Unmarshal([]byte(message.Body), &payload)

		if payload.Players > len(batch) {
			return nil, errors.New("not enough players")
		}

		p.PlayersCount = payload.Players
		p.Bucket = payload.Bucket
		p.Domain = payload.Domain
		p.Game = payload.Game
		p.Connections = append(p.Connections, &connection{
			ID:            payload.ConnectionID,
			ReceiptHandle: message.ReceiptHandle,
		})
		p.Language = payload.Language
		users = append(users, &ws.User{
			ConnectionID: payload.ConnectionID,
			Username:     payload.Username,
			UserID:       payload.UserID,
		})
	}
	p.Users = users
	return &p, nil
}

func setPool(p *pool, c *pooling.Config) {
	p.ID = uuid.New().String()
	_, pErr := dynamo().PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(fmt.Sprintf("%s_pools", c.Stage)),
		Item: map[string]*dynamodb.AttributeValue{
			"ID":            {S: aws.String(p.ID)},
			"ConnectionIDs": {SS: p.ConnectionIDsPtr()},
			"Bucket":        {S: aws.String(p.Bucket)},
			"Limit":         {N: aws.String(strconv.Itoa(p.PlayersCount))},
			"LiveUntil":     {N: aws.String(tools.FutureTimestampStr(36000))},
			"Language":      {S: aws.String(p.Language)},
		},
	})

	if pErr != nil {
		panic(pErr.Error())
	}

	requests := []*dynamodb.WriteRequest{}
	for _, cid := range p.ConnectionIDs() {
		requests = append(requests, &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: map[string]*dynamodb.AttributeValue{
					"ID":        {S: aws.String(cid)},
					"PoolID":    {S: aws.String(p.ID)},
					"LiveUntil": {N: aws.String(tools.FutureTimestampStr(7200))},
				},
			},
		})
	}
	_, err := dynamo().BatchWriteItem(&dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			fmt.Sprintf("%s_connections", c.Stage): requests,
		},
	})

	if err != nil {
		panic(err.Error())
	}
}

func handleDisconnections(event events.SQSEvent, p *pool, c *pooling.Config) error {
	flagDisconnections(p, c)

	sess := session.Must(session.NewSession())
	svc := sqs.New(sess)
	queueName := extractQueueName(event)
	entries := []*sqs.DeleteMessageBatchRequestEntry{}

	for _, c := range p.Connections {
		if c.IsDisconnected {
			entry := &sqs.DeleteMessageBatchRequestEntry{
				Id:            aws.String(uuid.New().String()),
				ReceiptHandle: aws.String(c.ReceiptHandle),
			}
			entries = append(entries, entry)
		}
	}

	if len(entries) == 0 {
		return nil
	}

	_, err := svc.DeleteMessageBatch(&sqs.DeleteMessageBatchInput{
		Entries:  entries,
		QueueUrl: aws.String(tools.BuildSqsURL(c.Region, c.AccountID, queueName)),
	})

	if err != nil {
		return err
	}

	return errors.New("disconnections exist")
}

func flagDisconnections(p *pool, c *pooling.Config) {
	o, err := tools.BatchGetItem(fmt.Sprintf("%s_disconnections", c.Stage), "ID", p.ConnectionIDs())

	if err != nil {
		fmt.Println(err.Error())
	}

	for _, items := range o.Responses {
		for _, kv := range items {
			for _, c := range p.Connections {
				if c.ID == *kv["ID"].S {
					c.IsDisconnected = true
				}
			}
		}
	}
}

func extractQueueName(event events.SQSEvent) string {
	arn := event.Records[0].EventSourceARN
	tokens := strings.Split(arn, ":")
	queueName := tokens[len(tokens)-1]
	return queueName
}

func dynamo() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSession())
	return dynamodb.New(sess)
}
