package scrabble

import (
	"math/rand"
	"strings"
	"time"
)

const _StartingTileCount = 100

// Bag represents a bag.
type Bag struct {
	Tiles     *Tiles `json:"t"`
	lastDrawn *Tiles
}

// NewBag creates a new bag.
func NewBag(language string) *Bag {
	rand.Seed(time.Now().UTC().UnixNano())
	switch strings.ToLower(language) {
	case "bulgarian":
		return NewBulgarianBag()
	case "english":
		return NewEnglishBag()
	default:
		panic("invalid language")
	}
}

// NewBulgarianBag creates a new English bag.
func NewBulgarianBag() *Bag {
	tiles := NewTiles(BlankTile(), BlankTile())
	tiles.Append(CreateMany("А", 1, 9).Value...)
	tiles.Append(CreateMany("Б", 2, 3).Value...)
	tiles.Append(CreateMany("В", 2, 4).Value...)
	tiles.Append(CreateMany("Г", 3, 3).Value...)
	tiles.Append(CreateMany("Д", 2, 4).Value...)
	tiles.Append(CreateMany("Е", 1, 8).Value...)
	tiles.Append(CreateMany("Ж", 4, 2).Value...)
	tiles.Append(CreateMany("З", 4, 2).Value...)
	tiles.Append(CreateMany("И", 1, 8).Value...)
	tiles.Append(CreateMany("Й", 5, 1).Value...)
	tiles.Append(CreateMany("К", 2, 3).Value...)
	tiles.Append(CreateMany("Л", 2, 3).Value...)
	tiles.Append(CreateMany("М", 2, 4).Value...)
	tiles.Append(CreateMany("Н", 1, 4).Value...)
	tiles.Append(CreateMany("О", 1, 9).Value...)
	tiles.Append(CreateMany("П", 1, 4).Value...)
	tiles.Append(CreateMany("Р", 1, 4).Value...)
	tiles.Append(CreateMany("С", 1, 4).Value...)
	tiles.Append(CreateMany("Т", 1, 5).Value...)
	tiles.Append(CreateMany("У", 5, 3).Value...)
	tiles.Append(CreateMany("Ф", 10, 1).Value...)
	tiles.Append(CreateMany("Х", 5, 1).Value...)
	tiles.Append(CreateMany("Ц", 8, 1).Value...)
	tiles.Append(CreateMany("Ч", 5, 2).Value...)
	tiles.Append(CreateMany("Ш", 8, 1).Value...)
	tiles.Append(CreateMany("Щ", 10, 1).Value...)
	tiles.Append(CreateMany("Ъ", 3, 2).Value...)
	tiles.Append(CreateMany("Ы", 10, 1).Value...)
	tiles.Append(CreateMany("Ю", 8, 1).Value...)
	tiles.Append(CreateMany("Я", 5, 2).Value...)
	bag := Bag{
		Tiles: tiles,
	}
	return &bag
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
