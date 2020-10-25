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
		Region:   aws.String("eu-central-1"),
		Endpoint: aws.String(endpoint),
	}))
	apiClient := apigatewaymanagementapi.New(session)
	connectionInput := apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: aws.String(context.ConnectionID),
		Data:         []byte(message),
	}
	request, _ := apiClient.PostToConnectionRequest(&connectionInput)
	postError := request.Send()

	fmt.Println(postError)
	if postError != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, postError
	}

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}
