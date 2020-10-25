package lambdaws

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
)

// Send sends a message to a connection ID.
func Send(connectionID string, message string) (events.APIGatewayProxyResponse, error) {
	session := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1"),
	}))
	apiClient := apigatewaymanagementapi.New(session)
	connectionInput := apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: aws.String(connectionID),
		Data:         []byte(message),
	}
	output, err := apiClient.PostToConnection(&connectionInput)

	fmt.Println(output.String())

	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}
