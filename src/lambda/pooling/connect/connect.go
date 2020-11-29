package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	lambdapooling "github.com/YouJinTou/vocabrace/lambda/pooling"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/YouJinTou/vocabrace/tools"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handle)
}

func handle(ctx context.Context, req *events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	c := lambdapooling.GetConfig()
	sess := session.Must(session.NewSession())
	svc := sqs.New(sess)
	queueName := buildQueueName(c, req.QueryStringParameters)
	marshalled, _ := json.Marshal(lambdapooling.PoolerPayload{
		Domain:       req.RequestContext.DomainName,
		ConnectionID: req.RequestContext.ConnectionID,
		Bucket:       "novice",
		Game:         req.QueryStringParameters["game"],
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

func buildQueueName(c *lambdapooling.Config, params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	queueName := fmt.Sprintf("%s_%s", c.Stage, params["game"])
	for _, k := range keys {
		queueName += fmt.Sprintf("_%s", params[k])
	}
	queueName += "_pooler"
	return queueName
}
