package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/google/uuid"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/YouJinTou/vocabrace/services/pooling"
	"github.com/YouJinTou/vocabrace/tools"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handle)
}

func handle(ctx context.Context, req *events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	c := pooling.GetConfig()
	sess := session.Must(session.NewSession())
	svc := sqs.New(sess)
	queueName := buildQueueName(c, req.QueryStringParameters)
	players, _ := strconv.Atoi(req.QueryStringParameters["players"])
	marshalled, _ := json.Marshal(pooling.PoolerPayload{
		Domain:       req.RequestContext.DomainName,
		ConnectionID: req.RequestContext.ConnectionID,
		Bucket:       getBucket(req.QueryStringParameters),
		Game:         req.QueryStringParameters["game"],
		Players:      players,
		Language:     getParam("language", "english", req.QueryStringParameters),
		UserID:       getParam("userId", uuid.New().String(), req.QueryStringParameters),
		Username:     getParam("username", "Anonymous", req.QueryStringParameters),
	})
	_, err := svc.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    aws.String(tools.BuildSqsURL(c.Region, c.AccountID, queueName)),
		MessageBody: aws.String(string(marshalled)),
	})

	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func buildQueueName(c *pooling.Config, p map[string]string) string {
	queueName := fmt.Sprintf("%s_%s_%s_%s_pooler", c.Stage, p["game"], p["language"], p["players"])
	return queueName
}

func getParam(key, fallback string, params map[string]string) string {
	if val, ok := params[key]; ok {
		return val
	}
	return fallback
}

func getBucket(params map[string]string) string {
	return pooling.Novice
}
