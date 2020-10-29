package lambdaws

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
)

// Message is used in Send().
type Message struct {
	Domain       string
	Stage        string
	ConnectionID string
	Message      string
}

// Send sends a message to a connection ID.
func Send(m *Message) (events.APIGatewayProxyResponse, error) {
	session := session.Must(session.NewSession())
	endpoint := fmt.Sprintf("https://%s/%s/", m.Domain, m.Stage)
	apiClient := apigatewaymanagementapi.New(session, aws.NewConfig().WithEndpoint(endpoint))
	connectionInput := apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: aws.String(m.ConnectionID),
		Data:         []byte(m.Message),
	}
	_, err := apiClient.PostToConnection(&connectionInput)

	if err != nil {
		fmt.Println(err.Error())

		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}
