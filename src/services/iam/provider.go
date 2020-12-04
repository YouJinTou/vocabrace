package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/YouJinTou/vocabrace/tools"

	"github.com/aws/aws-lambda-go/events"
)

func providerAuth(r *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	u := &User{}
	json.Unmarshal([]byte(r.Body), u)
	table := fmt.Sprintf("%s_iam_users", os.Getenv("STAGE"))
	_, err := tools.PutItem(table, u)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, err
	}
	return events.APIGatewayProxyResponse{StatusCode: 200}, err
}
