package scrabble

import (
	"math/rand"
)

const _StartingTileCount = 100

// Bag represents a bag.
type Bag struct {
	Tiles     []*Tile `json:"t"`
	lastDrawn []*Tile
}

// NewBag creates a new bag.
func NewBag(language language) *Bag {
	alphabet := language()
	tiles := []*Tile{BlankTile(), BlankTile()}
	for t := 0; t < 98; t++ {
		tiles = append(tiles, alphabet[rand.Intn(len(alphabet))].Copy(false))
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

	b.lastDrawn = make([]*Tile, len(tiles))
	copy(b.lastDrawn, tiles)

	return tiles
}

// GetLastDrawn returns a copy of the tiles that were last drawn.
func (b *Bag) GetLastDrawn() []*Tile {
	return b.lastDrawn
}

// Put puts tiles into the bag.
func (b *Bag) Put(tiles []*Tile) {
	for _, t := range tiles {
		b.Tiles = append(b.Tiles, t.Copy(true))
	}
}
