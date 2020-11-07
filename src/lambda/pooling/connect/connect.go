package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/YouJinTou/vocabrace/tools"

	ws "github.com/YouJinTou/vocabrace/lambda/pooling"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handle)
}

func handle(ctx context.Context, req *events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	c := ws.GetConfig()
	sess := session.Must(session.NewSession())
	svc := sqs.New(sess)
	game := req.QueryStringParameters["game"]
	bucket := req.QueryStringParameters["bucket"]
	queueName := fmt.Sprintf("%s_%s_pooler", c.Stage, game)
	marshalled, _ := json.Marshal(ws.PoolerPayload{
		Domain:       req.RequestContext.DomainName,
		ConnectionID: req.RequestContext.ConnectionID,
		Bucket:       bucket,
		Game:         game,
	})

	svc.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    aws.String(tools.BuildSqsURL(c.Region, c.AccountID, queueName)),
		MessageBody: aws.String(string(marshalled)),
	})

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}
