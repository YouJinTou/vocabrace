package data

import (
	"github.com/YouJinTou/vocabrace/services/com/ws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// State exposes methods that allow logic to occur on start/action.
type State interface {
	OnStart(OnStartInput) (OnStartOutput, error)
	OnAction(OnActionInput) (OnActionOutput, error)
	OnReconnect(OnReconnectInput) (OnReconnectOutput, error)
}

// OnStartInput encapsulates data to be processed during the start of the game.
type OnStartInput struct {
	Connections *Connections
}

// OnStartOutput encapsulates data ready to be persisted/sent after the setup logic.
type OnStartOutput struct {
	PoolID   string
	Messages []*ws.Message
	Game     interface{}
}

// OnActionInput encapsulates data for each turn.
type OnActionInput struct {
	Body            string
	Initiator       string
	InitiatorUserID string
	PoolID          string
	Connections     *Connections
	State           map[string]*dynamodb.AttributeValue
}

// OnActionOutput encapsulates data to send after each turn.
type OnActionOutput struct {
	Messages []*ws.Message
	Game     interface{}
	Error    *ws.Message
}

// OnReconnectInput encapsulates data required to perform a reconnection.
type OnReconnectInput struct {
	Connection Connection
	State      map[string]*dynamodb.AttributeValue
}

// OnReconnectOutput encapsulates data required by the client to visualize the current state.
type OnReconnectOutput struct {
	Message *ws.Message
}
