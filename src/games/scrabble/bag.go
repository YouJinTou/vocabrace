package scrabble

import (
	"math/rand"
)

const _StartingTileCount = 100

// Bag represents a bag.
type Bag struct {
	Tiles []*Tile `json:"t"`
}

// NewBag creates a new bag.
func NewBag(language language) *Bag {
	alphabet := language()
	tiles := []*Tile{BlankTile(), BlankTile()}
	for t := 0; t < 98; t++ {
		tiles = append(tiles, alphabet[rand.Intn(len(alphabet))].Copy())
	}
	bag := Bag{
		Tiles: tiles,
	}
	return &bag
}

// Count returns the current number of tiles in the bag.
func (b *Bag) Count() int {
	return len(b.Tiles)
}

// Draw selects a number of random tiles from the bag.
func (b *Bag) Draw(n int) []*Tile {
	if n > b.Count() {
		n = b.Count()
	}

	tiles := []*Tile{}

	for i := 0; i < n; i++ {
		idx := rand.Intn(b.Count())
		t := b.Tiles[idx]
		b.Tiles[idx] = b.Tiles[b.Count()-1]
		b.Tiles = b.Tiles[:b.Count()-1]

		tiles = append(tiles, t)
	}

	return tiles
}
