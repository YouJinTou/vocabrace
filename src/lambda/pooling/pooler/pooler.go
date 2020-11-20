package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/YouJinTou/vocabrace/ws"

	lambdapooling "github.com/YouJinTou/vocabrace/lambda/pooling"

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
	ID          string
	Connections []*connection
	Domain      string
	Bucket      string
	Game        string
	Stage       string
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

func handle(ctx context.Context, sqsEvent events.SQSEvent) error {
	c := lambdapooling.GetConfig()

	for {
		poolReady, batch := prepareBatch(c)

		if !poolReady {
			time.Sleep(time.Second * 1)
			continue
		}

		pool := createPool(batch, c)

		flagDisconnections(pool, c)

		if handleDisconnections(pool, c) {
			continue
		}

		setPool(pool, c)

		ws.OnStart(&ws.ReceiverData{
			ConnectionIDs: pool.ConnectionIDs(),
			Domain:        pool.Domain,
			Stage:         pool.Stage,
			Game:          pool.Game,
			PoolID:        pool.ID,
		})

		clearQueue(pool, c, false)
	}
}

func prepareBatch(c *lambdapooling.Config) (bool, []*sqs.Message) {
	queueName := fmt.Sprintf("%s_%s_pooler", c.Stage, "scrabble")
	messages := []*sqs.Message{}
	maxMessages := c.PoolLimit

	for i := 0; i < 3; i++ {
		o, _ := svc().ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(tools.BuildSqsURL(c.Region, c.AccountID, queueName)),
			WaitTimeSeconds:     aws.Int64(8),
			MaxNumberOfMessages: aws.Int64(int64(maxMessages)),
			VisibilityTimeout:   aws.Int64(25),
		})
		messages = append(messages, o.Messages...)
		maxMessages = maxMessages - len(messages)

		if len(messages) >= c.PoolLimit {
			messages = messages[0:c.PoolLimit]
			break
		}
	}

	if len(messages) < c.PoolLimit {
		return false, []*sqs.Message{}
	}

	return true, messages
}

func createPool(batch []*sqs.Message, c *lambdapooling.Config) *pool {
	p := pool{Connections: []*connection{}, Stage: c.Stage}
	for _, message := range batch {
		payload := lambdapooling.PoolerPayload{}
		json.Unmarshal([]byte(*message.Body), &payload)
		p.Bucket = payload.Bucket
		p.Domain = payload.Domain
		p.Game = payload.Game
		p.Connections = append(p.Connections, &connection{
			ID:            payload.ConnectionID,
			ReceiptHandle: *message.ReceiptHandle,
		})
	}
	return &p
}

func setPool(p *pool, c *lambdapooling.Config) {
	p.ID = uuid.New().String()
	_, pErr := dynamo().PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(fmt.Sprintf("%s_pools", c.Stage)),
		Item: map[string]*dynamodb.AttributeValue{
			"ID":            {S: aws.String(p.ID)},
			"ConnectionIDs": {SS: p.ConnectionIDsPtr()},
			"Bucket":        {S: aws.String(p.Bucket)},
			"Limit":         {N: aws.String(c.PoolLimitStr)},
			"LiveUntil":     {N: aws.String(tools.FutureTimestampStr(36000))},
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

func flagDisconnections(p *pool, c *lambdapooling.Config) {
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

func handleDisconnections(p *pool, c *lambdapooling.Config) bool {
	return clearQueue(p, c, true)
}

func clearQueue(p *pool, c *lambdapooling.Config, disconnectionsOnly bool) bool {
	svc := svc()
	queueName := fmt.Sprintf("%s_%s_pooler", c.Stage, p.Game)
	entries := []*sqs.DeleteMessageBatchRequestEntry{}
	hasDeleted := false

	for _, c := range p.Connections {
		entry := &sqs.DeleteMessageBatchRequestEntry{
			Id:            aws.String(uuid.New().String()),
			ReceiptHandle: aws.String(c.ReceiptHandle),
		}
		if disconnectionsOnly && c.IsDisconnected {
			entries = append(entries, entry)
		} else if !disconnectionsOnly {
			entries = append(entries, entry)
		}
	}

	if len(entries) > 0 {
		_, err := svc.DeleteMessageBatch(&sqs.DeleteMessageBatchInput{
			Entries:  entries,
			QueueUrl: aws.String(tools.BuildSqsURL(c.Region, c.AccountID, queueName)),
		})

		if err != nil {
			fmt.Println(err.Error())
		} else {
			hasDeleted = true
		}
	}

	return hasDeleted
}

func dynamo() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSession())
	return dynamodb.New(sess)
}

func svc() *sqs.SQS {
	sess := session.Must(session.NewSession())
	return sqs.New(sess)
}
