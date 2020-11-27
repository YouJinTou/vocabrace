package scrabble

import (
	"strconv"
	"testing"
)

var counts = []int{0, 1, 2, 100, 1000}
var puts = NewTiles(
	&Tile{Letter: "ZZ", ID: "abc"},
	&Tile{Letter: "QQ", ID: "qqf"},
	&Tile{Letter: "FF", ID: "ffx"},
)

func TestVerifyUniqueTiles(t *testing.T) {
	b := NewBag(English)
	m := map[string]int{}
	for _, t := range b.Tiles.Value {
		if _, ok := m[t.ID]; ok {
			m[t.ID]++
		} else {
			m[t.ID] = 1
		}
	}

	for k, v := range m {
		if v > 1 {
			t.Errorf("id %s should be there once, actually %d times", k, v)
		}
	}
}

func TestDrawRemovesTiles(t *testing.T) {
	for _, c := range counts {
		t.Run(strconv.Itoa(c), func(t *testing.T) {
			b := NewBag(English)
			drawn := b.Draw(c)
			for j, bt := range b.Tiles.Value {
				for _, d := range drawn.Value {
					if bt.ID == d.ID {
						t.Errorf("%s still there when count is %d (at %d)", bt.ID, c, j)
					}
				}
			}
		})
	}
}

func TestDrawPutCount(t *testing.T) {
	b := NewBag(English)
	originalCount := b.Count()

	b.Draw(puts.Count())
	b.Put(puts)

	if b.Count() != originalCount {
		t.Errorf("Draw/Put count mismatch.")
	}
}

func TestPutAddsTiles(t *testing.T) {
	b := NewBag(English)
	originalCount := b.Count()
	expected := originalCount + puts.Count()
	b.Put(puts)

	if expected != b.Count() {
		t.Errorf("expected count %d, got %d", expected, b.Count())
	}

	for _, p := range puts.Value {
		found := false
		for _, bt := range b.Tiles.Value {
			if p.ID == bt.ID {
				found = true
			}
		}
		if !found {
			t.Errorf("tile not found")
		}
	}
}

func TestDrawCount(t *testing.T) {
	for _, c := range counts {
		t.Run(strconv.Itoa(c), func(t *testing.T) {
			b := NewBag(English)
			originalCount := b.Count()
			drawn := b.Draw(c)
			result := originalCount - drawn.Count()
			if result != b.Count() {
				t.Errorf("got %q, want %q", result, drawn.Count())
			}
		})
	}
}

func TestBulgarianBagHasCorrectDistribution(t *testing.T) {
	b := NewBulgarianBag()
	tests := []*dTile{
		{"", 2, 0},
		{"А", 9, 1},
		{"О", 9, 1},
		{"Е", 8, 1},
		{"И", 8, 1},
		{"Т", 5, 1},
		{"Н", 4, 1},
		{"П", 4, 1},
		{"Р", 4, 1},
		{"С", 4, 1},
		{"В", 4, 2},
		{"Д", 4, 2},
		{"М", 4, 2},
		{"Б", 3, 2},
		{"К", 3, 2},
		{"Л", 3, 2},
		{"Г", 3, 3},
		{"Ъ", 2, 3},
		{"Ж", 2, 4},
		{"З", 2, 4},
		{"У", 3, 5},
		{"Ч", 2, 5},
		{"Я", 2, 5},
		{"Й", 1, 5},
		{"Х", 1, 5},
		{"Ц", 1, 8},
		{"Ш", 1, 8},
		{"Ю", 1, 8},
		{"Ф", 1, 10},
		{"Щ", 1, 10},
		{"Ь", 1, 10},
	}
	bagTestDistribution(tests, b, t)
}

func TestBagsHaveCorrectCount(t *testing.T) {
	bagTestCount(102, NewBulgarianBag(), t)
	bagTestCount(100, NewEnglishBag(), t)
}

func TestBagsTilesHaveUniqueIDs(t *testing.T) {
	bagTestUniqueIDs(NewBulgarianBag(), t)
	bagTestUniqueIDs(NewEnglishBag(), t)
}

func TestEnglishBagHasCorrectDistribution(t *testing.T) {
	b := NewEnglishBag()
	tests := []*dTile{
		&dTile{"", 2, 0},
		&dTile{"A", 9, 1},
		&dTile{"B", 2, 3},
		&dTile{"C", 2, 3},
		&dTile{"D", 4, 2},
		&dTile{"E", 12, 1},
		&dTile{"F", 2, 4},
		&dTile{"G", 3, 2},
		&dTile{"H", 2, 4},
		&dTile{"I", 9, 1},
		&dTile{"J", 1, 8},
		&dTile{"K", 1, 5},
		&dTile{"L", 4, 1},
		&dTile{"M", 2, 3},
		&dTile{"N", 6, 1},
		&dTile{"O", 8, 1},
		&dTile{"P", 2, 3},
		&dTile{"Q", 1, 10},
		&dTile{"R", 6, 1},
		&dTile{"S", 4, 1},
		&dTile{"T", 6, 1},
		&dTile{"U", 4, 1},
		&dTile{"V", 2, 4},
		&dTile{"W", 2, 4},
		&dTile{"X", 1, 8},
		&dTile{"Y", 2, 4},
		&dTile{"Z", 1, 10},
	}
	bagTestDistribution(tests, b, t)
}

func bagTestCount(expected int, b *Bag, t *testing.T) {
	if b.Count() != expected {
		t.Errorf("expected %d, got %d", expected, b.Count())
	}
}

func bagTestUniqueIDs(b *Bag, t *testing.T) {
	m := map[string]int{}
	for _, tile := range b.Tiles.Value {
		if _, ok := m[tile.ID]; ok {
			t.Errorf("tile ID repeats")
		} else {
			m[tile.ID] = 1
		}
	}
}

type dTile struct {
	Letter string
	Count  int
	Value  int
}

func bagTestDistribution(tests []*dTile, b *Bag, t *testing.T) {
	m := map[string][]*Tile{}
	for _, tile := range b.Tiles.Value {
		if _, ok := m[tile.Letter]; ok {
			m[tile.Letter] = append(m[tile.Letter], tile)
		} else {
			m[tile.Letter] = []*Tile{tile}
		}
	}
	for _, tt := range tests {
		t.Run(tt.Letter, func(t *testing.T) {
			if _, ok := m[tt.Letter]; !ok {
				t.Errorf("letter not found %s", tt.Letter)
			}
			if len(m[tt.Letter]) != tt.Count {
				t.Errorf("letter %s count got %d, expected %d", tt.Letter, len(m[tt.Letter]), tt.Count)
			}
			if m[tt.Letter][0].Value != tt.Value {
				t.Errorf("letter %s value got %d, expected %d", tt.Letter, m[tt.Letter][0].Value, tt.Value)
			}
		})
	}
}
