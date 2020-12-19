package data

import (
	"github.com/YouJinTou/vocabrace/tools"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Connection encapsulates connection information.
type Connection struct {
	ID       string
	Domain   string
	Game     string
	Language string
	UserID   string
}

// Connections encapsulates a slice of connections.
type Connections struct{ c []*Connection }

// NewConnections creates a set of connections with shared characteristics.
func NewConnections(c []*Connection) *Connections {
	return &Connections{c: c}
}

// IDs gets all the connection IDs.
func (v *Connections) IDs() []string {
	ids := []string{}
	for _, c := range v.c {
		ids = append(ids, c.ID)
	}
	return ids
}

// UserIDs gets all the connection user IDs.
func (v *Connections) UserIDs() []string {
	ids := []string{}
	for _, c := range v.c {
		ids = append(ids, c.UserID)
	}
	return ids
}

// IDByUserID returns the corresponding connection ID.
func (v *Connections) IDByUserID(ID string) *string {
	for _, c := range v.c {
		if c.UserID == ID {
			return &c.ID
		}
	}
	return nil
}

// UserIDByID returns the corresponding user ID.
func (v *Connections) UserIDByID(ID string) *string {
	for _, c := range v.c {
		if c.ID == ID {
			return &c.UserID
		}
	}
	return nil
}

// Game gets the connections game.
func (v *Connections) Game() string {
	return v.c[0].Game
}

// Domain gets the connections domain.
func (v *Connections) Domain() string {
	return v.c[0].Domain
}

// Language gets the connections language.
func (v *Connections) Language() string {
	return v.c[0].Language
}

// OtherIDs gets all connection IDs except for the target.
func (v *Connections) OtherIDs(target string) []string {
	connections := []string{}
	for _, c := range v.c {
		if c.ID != target {
			connections = append(connections, c.ID)
		}
	}
	return connections
}

// GetConnections gets the entire connection data based on a set of connection IDs.
func GetConnections(connectionIDs []string) (*Connections, error) {
	table := tools.Table("connections")
	o, err := tools.BatchGetItem(table, "ID", connectionIDs)
	if err != nil {
		return nil, err
	}
	connections := []*Connection{}
	for _, response := range o.Responses[*table] {
		con := &Connection{}
		dynamodbattribute.UnmarshalMap(response, con)
		connections = append(connections, con)
	}
	return NewConnections(connections), nil
}
