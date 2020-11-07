package ws

import (
	"fmt"
	"sync"

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

// PoolPayload encapsulated data needed to push a message to a pool of connections.
type PoolPayload struct {
	Domain string
	PoolID string
	Bucket string
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

// SendToPeers sends a message to all peers.
func SendToPeers(connectionIDs []string, m Message) {
	var wg sync.WaitGroup

	for _, cid := range connectionIDs {
		wg.Add(1)

		m.ConnectionID = cid

		go send(&wg, m)
	}

	wg.Wait()
}

func send(wg *sync.WaitGroup, m Message) {
	defer wg.Done()

	Send(&m)
}
