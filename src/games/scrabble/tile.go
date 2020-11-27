package scrabble

import (
	"encoding/json"
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

// Tiles is a set of tiles.
type Tiles struct {
	Value []*Tile `json:"t"`
}

func tileID() string {
	return uuid.New().String()[0:6]
}

// NewTile creates a new tile.
func NewTile(letter string, value int) *Tile {
	return &Tile{letter, value, tileID()}
}

// NewTileWithID creates a new tile with a given ID.
func NewTileWithID(ID, letter string, value int) *Tile {
	t := NewTile(letter, value)
	t.ID = ID
	return t
}

// NewTiles creates new tiles.
func NewTiles(tiles ...*Tile) *Tiles {
	value := []*Tile{}
	for _, t := range tiles {
		value = append(value, &Tile{
			Letter: t.Letter,
			Value:  t.Value,
			ID:     t.ID,
		})
	}
	return &Tiles{Value: value}
}

// BlankTile creates a blank tile.
func BlankTile() *Tile {
	return &Tile{
		Letter: "",
		Value:  0,
		ID:     tileID(),
	}
}

// CreateMany creates many tiles of the same type.
func CreateMany(letter string, value, count int) Tiles {
	result := Tiles{}
	for i := 0; i < count; i++ {
		result.Append(NewTile(letter, value))
	}
	return result
}

// FromString creates a tile from a string.
func FromString(s string) *Tile {
	tokens := strings.Split(s, "|")

	if len(tokens) != 3 {
		return nil
	}

	value, err := strconv.Atoi(tokens[2])

	if err != nil {
		return nil
	}

	return &Tile{
		ID:     tokens[0],
		Letter: tokens[1],
		Value:  value,
	}
}

// Append appends a tile to the current tiles.
func (t *Tiles) Append(tiles ...*Tile) *Tiles {
	for _, it := range tiles {
		t.Value = append(t.Value, it)
	}
	return t
}

// RemoveAt removes a tile at a given index.
func (t *Tiles) RemoveAt(i int) *Tile {
	curr := t.Value[i]
	t.Value[i] = t.Value[len(t.Value)-1]
	t.Value = t.Value[:len(t.Value)-1]
	return curr
}

// GetAt gets a tile by index.
func (t *Tiles) GetAt(i int) *Tile {
	return t.Value[i]
}

// Count returns a count of all the tiles in the set.
func (t *Tiles) Count() int {
	return len(t.Value)
}

// RemoveByID removes a tile by an ID.
func (t *Tiles) RemoveByID(ID string) *Tile {
	newTiles := NewTiles()
	var toReturn *Tile
	for _, i := range t.Value {
		if i.ID == ID {
			toReturn = i
		} else {
			newTiles.Append(i)
		}
	}

	t.Value = newTiles.Value

	return toReturn
}

// FindByID finds a tile by an ID.
func (t *Tiles) FindByID(ID string) *Tile {
	for _, i := range t.Value {
		if i.ID == ID {
			return i
		}
	}
	return nil
}

// IsBlank returns whether a tile is blank by default.
func (t *Tile) IsBlank() bool {
	return t.Value == 0
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

// MarshalJSON serializes Tiles as a list of strings.
func (t Tiles) MarshalJSON() ([]byte, error) {
	tiles := []string{}
	for _, i := range t.Value {
		tiles = append(tiles, i.String())
	}
	return json.Marshal(tiles)
}

// UnmarshalJSON deserializes Tiles back into its shape.
func (t *Tiles) UnmarshalJSON(b []byte) error {
	tileStrings := []string{}
	tiles := NewTiles()
	json.Unmarshal(b, &tileStrings)

	for _, s := range tileStrings {
		tiles.Append(FromString(s))
	}

	t.Value = tiles.Value

	return nil
}
