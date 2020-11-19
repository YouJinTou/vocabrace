package scrabble

import (
	"math/rand"
	"strings"
)

const _StartingTileCount = 100

// Bag represents a bag.
type Bag struct {
	Tiles     *Tiles `json:"t"`
	lastDrawn *Tiles
}

// NewBag creates a new bag.
func NewBag(language string) *Bag {
	switch strings.ToLower(language) {
	case "english":
		return NewEnglishBag()
	default:
		panic("invalid language")
	}
}

// NewEnglishBag creates a new English bag.
func NewEnglishBag() *Bag {
	tiles := NewTiles(BlankTile(), BlankTile())
	tiles.Append(CreateMany("A", 1, 9).Value...)
	tiles.Append(CreateMany("B", 3, 2).Value...)
	tiles.Append(CreateMany("C", 3, 2).Value...)
	tiles.Append(CreateMany("D", 2, 4).Value...)
	tiles.Append(CreateMany("E", 1, 12).Value...)
	tiles.Append(CreateMany("F", 4, 2).Value...)
	tiles.Append(CreateMany("G", 2, 3).Value...)
	tiles.Append(CreateMany("H", 4, 2).Value...)
	tiles.Append(CreateMany("I", 1, 9).Value...)
	tiles.Append(CreateMany("J", 8, 1).Value...)
	tiles.Append(CreateMany("K", 5, 1).Value...)
	tiles.Append(CreateMany("L", 1, 4).Value...)
	tiles.Append(CreateMany("M", 3, 2).Value...)
	tiles.Append(CreateMany("N", 1, 6).Value...)
	tiles.Append(CreateMany("O", 1, 8).Value...)
	tiles.Append(CreateMany("P", 3, 2).Value...)
	tiles.Append(CreateMany("Q", 10, 1).Value...)
	tiles.Append(CreateMany("R", 1, 6).Value...)
	tiles.Append(CreateMany("S", 1, 4).Value...)
	tiles.Append(CreateMany("T", 1, 6).Value...)
	tiles.Append(CreateMany("U", 1, 4).Value...)
	tiles.Append(CreateMany("V", 4, 2).Value...)
	tiles.Append(CreateMany("W", 4, 2).Value...)
	tiles.Append(CreateMany("X", 8, 1).Value...)
	tiles.Append(CreateMany("Y", 4, 2).Value...)
	tiles.Append(CreateMany("Z", 10, 1).Value...)
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
