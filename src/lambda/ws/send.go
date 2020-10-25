package lambdaws

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
)

// Send sends a message to a connection ID.
func Send(context *events.APIGatewayWebsocketProxyRequestContext, message string) (events.APIGatewayProxyResponse, error) {
	endpoint := fmt.Sprintf("https://%s", context.DomainName)
	session := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1"),
	}))
	apiClient := apigatewaymanagementapi.New(session, &aws.Config{
		Endpoint: aws.String(endpoint),
	})
	connectionInput := apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: aws.String(context.ConnectionID),
		Data:         []byte(message),
	}
	request, _ := apiClient.PostToConnectionRequest(&connectionInput)
	postError := request.Send()

	if postError != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, postError
	}

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}
