package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/YouJinTou/vocabrace/pooling"
	"github.com/YouJinTou/vocabrace/tools"

	ws "github.com/YouJinTou/vocabrace/lambda/pooling"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
)

func main() {
	lambda.Start(handle)
}

func handle(ctx context.Context, req *events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	c := ws.GetConfig()
	con := pooling.NewMemcachedContext(c.MemcachedHost, c.MemcachedUsername, c.MemcachedPassword)
	p, err := con.JoinOrCreate(&pooling.Request{
		ConnectionID: req.RequestContext.ConnectionID,
		UserID:       uuid.New().String(),
		PoolLimit:    3})

	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, nil
	}

	peers, peersErr := con.GetPeers(req.RequestContext.ConnectionID)

	if peersErr != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, nil
	}

	notifyConductor(p, c)

	ws.SendToPeers(peers, ws.Message{
		Domain:  req.RequestContext.DomainName,
		Stage:   c.Stage,
		Message: "Just connected.",
	})

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func notifyConductor(p *pooling.Pool, c *ws.Config) {
	sess := session.Must(session.NewSession())
	svc := sqs.New(sess)
	queueName := fmt.Sprintf("%s_conductor", c.Stage)
	svc.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    aws.String(tools.BuildSqsURL(c.Region, c.AccountID, queueName)),
		MessageBody: aws.String(p.ID),
	})
}
