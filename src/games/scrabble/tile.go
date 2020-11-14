package scrabble

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
)

// Tile represents a tile.
type Tile struct {
	Letter string `json:"l"`
	Value  int    `json:"v"`
	ID     string `json:"i"`
}

func tileID() string {
	return uuid.New().String()[0:5]
}

// NewTile creates a new tile.
func NewTile(letter string, value int) *Tile {
	return &Tile{letter, value, tileID()}
}

// BlankTile creates a blank tile.
func BlankTile() *Tile {
	return &Tile{
		Letter: "",
		Value:  0,
		ID:     tileID(),
	}
}

// String returns a string representation of the tile.
func (t *Tile) String() string {
	return fmt.Sprintf("%s|%s|%d", t.ID, t.Letter, t.Value)
}

// Copy copies a tile.
func (t *Tile) Copy(preserveIndex bool) *Tile {
	var idx string
	if preserveIndex {
		idx = t.ID
	} else {
		idx = tileID()
	}
	return &Tile{
		Letter: t.Letter,
		Value:  t.Value,
		ID:     idx,
	}
}

// MarshalDynamoDBAttributeValue marshals the tile to a DynamoDB string.
func (t *Tile) MarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	av.S = aws.String(t.String())
	return nil
}

// UnmarshalDynamoDBAttributeValue unmarshals a DynamoDB string into a Tile.
func (t *Tile) UnmarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	tokens := strings.Split(*av.S, "|")

	if len(tokens) != 3 {
		return errors.New("expected three tile tokens")
	}

	value, err := strconv.Atoi(tokens[2])

	if err != nil {
		return errors.New("expected third token to be a number")
	}

	t.ID = tokens[0]
	t.Letter = tokens[1]
	t.Value = value

	return nil
}
