package main

import (
	"context"
	"encoding/json"
	"fmt"

	dynamodbpooling "github.com/YouJinTou/vocabrace/pooling/providers/dynamodb"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/YouJinTou/vocabrace/pooling"
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
	provider := dynamodbpooling.NewDynamoDBProvider(c.Stage)
	p, err := provider.JoinOrCreate(&pooling.JoinOrCreateInput{
		ConnectionID: req.RequestContext.ConnectionID,
		PoolLimit:    3,
		Bucket:       pooling.Novice,
	})

	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, nil
	}

	notifyConductor(p, c, req.RequestContext.DomainName)

	peers := p.GetPeers(req.RequestContext.ConnectionID)

	ws.SendToPeers(peers, ws.Message{
		Domain:  req.RequestContext.DomainName,
		Stage:   c.Stage,
		Message: "Just connected.",
	})

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func notifyConductor(p *pooling.Pool, c *ws.Config, domain string) {
	sess := session.Must(session.NewSession())
	svc := sqs.New(sess)
	queueName := fmt.Sprintf("%s_conductor", c.Stage)
	marshalled, _ := json.Marshal(ws.PoolPayload{
		Domain: domain,
		PoolID: p.ID,
		Bucket: p.Bucket,
	})

	svc.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    aws.String(tools.BuildSqsURL(c.Region, c.AccountID, queueName)),
		MessageBody: aws.String(string(marshalled)),
	})
}
