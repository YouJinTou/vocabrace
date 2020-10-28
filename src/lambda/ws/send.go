package lambdaws

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
)

// Payload is the payload to post to a connection
type Payload struct {
	Action string `json:"action"`
	Data   string `json:"data"`
}

// Send sends a message to a connection ID.
func Send(domain, stage, connectionID, message string) (events.APIGatewayProxyResponse, error) {
	session := session.Must(session.NewSession())
	endpoint := fmt.Sprintf("https://%s/%s/", domain, stage)
	apiClient := apigatewaymanagementapi.New(session, aws.NewConfig().WithEndpoint(endpoint))
	connectionInput := apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: aws.String(connectionID),
		Data:         []byte(message),
	}
	_, err := apiClient.PostToConnection(&connectionInput)

	if err != nil {
		fmt.Println(err.Error())

		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}
