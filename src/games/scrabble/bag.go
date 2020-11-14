package scrabble

import (
	"math/rand"
)

const _StartingTileCount = 100

// Bag represents a bag.
type Bag struct {
	Tiles     *Tiles `json:"t"`
	lastDrawn *Tiles
}

// NewBag creates a new bag.
func NewBag(language language) *Bag {
	alphabet := language()
	tiles := NewTiles(BlankTile(), BlankTile())
	for t := 0; t < 98; t++ {
		i := rand.Intn(alphabet.Count())
		tiles.Append(alphabet.GetAt(i).Copy(false))
	}
	bag := Bag{
		Tiles: tiles,
	}
	return &bag
}

// Count returns the current number of tiles in the bag.
func (b *Bag) Count() int {
	return b.Tiles.Count()
}

// Draw selects a number of random tiles from the bag.
func (b *Bag) Draw(n int) *Tiles {
	if n > b.Count() {
		n = b.Count()
	}

	tiles := NewTiles()

	for i := 0; i < n; i++ {
		idx := rand.Intn(b.Count())
		t := b.Tiles.RemoveAt(idx)
		tiles.Append(t)
	}

	b.lastDrawn = NewTiles(tiles.Value...)

	return tiles
}

// GetLastDrawn returns a copy of the tiles that were last drawn.
func (b *Bag) GetLastDrawn() *Tiles {
	return b.lastDrawn
}

// Put puts tiles into the bag.
func (b *Bag) Put(tiles *Tiles) {
	for _, t := range tiles.Value {
		b.Tiles.Append(t.Copy(true))
	}
}
