package ws

import (
	"fmt"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
)

// Message is used in sending data over a websocket.
type Message struct {
	Domain       string
	ConnectionID string
	Message      string
}

// Send sends a message to a connection ID.
func Send(m *Message) error {
	session := session.Must(session.NewSession())
	endpoint := fmt.Sprintf("https://%s/%s/", m.Domain, os.Getenv("STAGE"))
	apiClient := apigatewaymanagementapi.New(session, aws.NewConfig().WithEndpoint(endpoint))
	connectionInput := apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: aws.String(m.ConnectionID),
		Data:         []byte(m.Message),
	}

	for i := 0; i < 3; i++ {
		_, err := apiClient.PostToConnection(&connectionInput)

		if err == nil {
			return nil
		}

		fmt.Println(err.Error())
	}

	return nil
}

// SendManyUnique sends a unique message per connection.
func SendManyUnique(messages []*Message) {
	var wg sync.WaitGroup

	for _, m := range messages {
		wg.Add(1)

		go send(&wg, *m)
	}

	wg.Wait()
}

// SendMany sends a message to all peers.
func SendMany(connectionIDs []string, m *Message) {
	var wg sync.WaitGroup

	for _, cid := range connectionIDs {
		wg.Add(1)

		m.ConnectionID = cid

		go send(&wg, *m)
	}

	wg.Wait()
}

func send(wg *sync.WaitGroup, m Message) {
	defer wg.Done()

	Send(&m)
}
