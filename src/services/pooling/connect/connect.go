package main

import (
	"context"
	"strconv"

	"github.com/YouJinTou/vocabrace/tools"
	"github.com/google/uuid"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func main() {
	lambda.Start(handle)
}

func handle(ctx context.Context, req *events.APIGatewayWebsocketProxyRequest) (
	events.APIGatewayProxyResponse, error) {
	connect(tools.PutItem, req)
	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func connect(
	putItem func(*string, interface{}) (*dynamodb.PutItemOutput, error),
	r *events.APIGatewayWebsocketProxyRequest) {
	pid, pidExists := r.QueryStringParameters["pid"]
	_, userIDExists := r.QueryStringParameters["userID"]
	connection := struct {
		ID             string
		Game           string
		Players        int
		Language       string
		UserID         string
		Bucket         string
		Domain         string
		LiveUntil      int
		IsReconnection bool
		PoolID         string
	}{
		r.RequestContext.ConnectionID,
		r.QueryStringParameters["game"],
		players(r.QueryStringParameters),
		r.QueryStringParameters["language"],
		userID(r.QueryStringParameters),
		"novice",
		r.RequestContext.DomainName,
		tools.FutureTimestamp(7200),
		pidExists && userIDExists,
		pid,
	}
	putItem(tools.Table("connections"), connection)
}

func players(params map[string]string) int {
	players, _ := strconv.Atoi(params["players"])
	return players
}

func userID(params map[string]string) string {
	if ID, ok := params["userID"]; ok {
		return ID
	}
	return uuid.New().String()
}
