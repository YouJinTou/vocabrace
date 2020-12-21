package facts

import (
	"encoding/json"
	"log"

	"github.com/YouJinTou/vocabrace/tools"
)

// Fact encapsulates a fact.
type Fact struct {
	ID        string `json:"i"`
	Timestamp int    `json:"u"`
	Type      string `json:"t"`
	Data      string `json:"d"`
}

// NewFact creates a new fact.
func NewFact(ID, factType string, data interface{}) Fact {
	b, _ := json.Marshal(data)
	s := string(b)
	return Fact{
		ID:        ID,
		Timestamp: tools.FutureTimestamp(0),
		Type:      factType,
		Data:      s,
	}
}

// Publish publishes an fact.
func Publish(f Fact) {
	if _, err := tools.PutItem(tools.Table("facts"), f); err != nil {
		log.Printf(err.Error())
	}
}
