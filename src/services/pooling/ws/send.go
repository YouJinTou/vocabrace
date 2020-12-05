package ws

import (
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
)

// ReceiverData encapsulates receiver data for AWS API Gateway websockets.
type ReceiverData struct {
	Initiator     string
	ConnectionIDs []string
	Domain        string
	Stage         string
	PoolID        string
	Game          string
	Body          string
}

func (rd *ReceiverData) otherConnections() []string {
	connections := []string{}
	for _, cid := range rd.ConnectionIDs {
		if cid != rd.Initiator {
			connections = append(connections, cid)
		}
	}
	return connections
}

// Message is used in sending data over a websocket.
type Message struct {
	Domain       string
	Stage        string
	ConnectionID string
	Message      string
}

// Send sends a message to a connection ID.
func Send(m *Message) error {
	session := session.Must(session.NewSession())
	endpoint := fmt.Sprintf("https://%s/%s/", m.Domain, m.Stage)
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
