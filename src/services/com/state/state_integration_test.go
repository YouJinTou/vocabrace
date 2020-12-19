package state

import (
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type game struct {
	A string
}

func Test_SavesLoadsState(t *testing.T) {
	os.Setenv("AWS_REGION", "eu-central-1")
	os.Setenv("STAGE", "dev")
	os.Setenv("AWS_PROFILE", "vocabrace")

	pid := "testing_pid"
	g := game{A: "A"}

	if err := saveState(&saveStateInput{
		PoolID:        pid,
		ConnectionIDs: []string{"123", "456"},
		Game:          g,
	}); err != nil {
		t.Errorf(err.Error())
	}

	game := &game{}
	m := loadState(pid)
	dynamodbattribute.UnmarshalMap(m, game)

	if game.A != g.A {
		t.Errorf("values mismatch after load")
	}
}
