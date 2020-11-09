package main

import (
	"context"
	"encoding/json"
	"fmt"

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
	game := req.QueryStringParameters["game"]
	bucket := req.QueryStringParameters["bucket"]
	queueName := fmt.Sprintf("%s_%s_pooler", c.Stage, game)
	marshalled, _ := json.Marshal(lambdapooling.PoolerPayload{
		Domain:       req.RequestContext.DomainName,
		ConnectionID: req.RequestContext.ConnectionID,
		Bucket:       bucket,
		Game:         game,
	})
	_, err := svc.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    aws.String(tools.BuildSqsURL(c.Region, c.AccountID, queueName)),
		MessageBody: aws.String(string(marshalled)),
	})

	return events.APIGatewayProxyResponse{StatusCode: 200}, err
}
