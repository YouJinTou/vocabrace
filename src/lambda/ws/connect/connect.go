package main

import (
	"context"
	"fmt"
	"os"

	"github.com/YouJinTou/vocabrace/pool"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	"github.com/tkanos/gonfig"
)

func main() {
	lambda.Start(handle)
}

func handle(ctx context.Context, req events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	config := getConfig()
	p := pool.New(config)
	p.JoinOrCreate(&pool.Request{
		ConnectionID: req.RequestContext.ConnectionID,
		UserID:       uuid.New().String(),
		PoolLimit:    5})
	// endpoint := fmt.Sprintf(
	// 	"https://%s.execute-api.%s.amazonaws.com/%s",
	// 	req.RequestContext.DomainName,
	// 	"eu-central-1",
	// 	req.RequestContext.Stage)
	// session := session.Must(session.NewSession(&aws.Config{Region: aws.String("eu-central-1")}))
	// apiClient := apigatewaymanagementapi.New(session, aws.NewConfig().WithEndpoint(endpoint))
	// connectionInput := apigatewaymanagementapi.PostToConnectionInput{
	// 	ConnectionId: aws.String(req.RequestContext.ConnectionID),
	// 	Data:         []byte("{ poolID: \"123\" }"),
	// }
	// request, _ := apiClient.PostToConnectionRequest(&connectionInput)
	// postError := request.Send()

	// if postError != nil {
	// 	return events.APIGatewayProxyResponse{StatusCode: 500}, postError
	// }

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func getConfig() *pool.Config {
	config := pool.Config{}
	stage := os.Getenv("STAGE")
	file := fmt.Sprintf("config.%s.json", stage)
	err := gonfig.GetConf(file, &config)

	if err != nil {
		panic(err)
	}

	return &config
}
