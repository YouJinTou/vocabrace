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

func TestEnglishBagHasCorrectCount(t *testing.T) {
	b := NewEnglishBag()
	if b.Count() != 100 {
		t.Errorf("expected 100, got %d", b.Count())
	}
}

func TestEnglishBagTilesHaveUniqueIDs(t *testing.T) {
	b := NewEnglishBag()
	m := map[string]int{}
	for _, tile := range b.Tiles.Value {
		if _, ok := m[tile.ID]; ok {
			t.Errorf("tile ID repeats")
		} else {
			m[tile.ID] = 1
		}
	}
	if b.Count() != 100 {
		t.Errorf("expected 100, got %d", b.Count())
	}
}

func TestEnglishBagHasCorrectDistribution(t *testing.T) {
	b := NewEnglishBag()
	tests := []struct {
		letter string
		count  int
		value  int
	}{
		{"", 2, 0},
		{"A", 9, 1},
		{"B", 2, 3},
		{"C", 2, 3},
		{"D", 4, 2},
		{"E", 12, 1},
		{"F", 2, 4},
		{"G", 3, 2},
		{"H", 2, 4},
		{"I", 9, 1},
		{"J", 1, 8},
		{"K", 1, 5},
		{"L", 4, 1},
		{"M", 2, 3},
		{"N", 6, 1},
		{"O", 8, 1},
		{"P", 2, 3},
		{"Q", 1, 10},
		{"R", 6, 1},
		{"S", 4, 1},
		{"T", 6, 1},
		{"U", 4, 1},
		{"V", 2, 4},
		{"W", 2, 4},
		{"X", 1, 8},
		{"Y", 2, 4},
		{"Z", 1, 10},
	}
	m := map[string][]*Tile{}
	for _, tile := range b.Tiles.Value {
		if _, ok := m[tile.Letter]; ok {
			m[tile.Letter] = append(m[tile.Letter], tile)
		} else {
			m[tile.Letter] = []*Tile{tile}
		}
	}
	for _, tt := range tests {
		t.Run(tt.letter, func(t *testing.T) {
			if _, ok := m[tt.letter]; !ok {
				t.Errorf("letter not found %s", tt.letter)
			}
			if len(m[tt.letter]) != tt.count {
				t.Errorf("letter %s count got %d, expected %d", tt.letter, len(m[tt.letter]), tt.count)
			}
			if m[tt.letter][0].Value != tt.value {
				t.Errorf("letter %s value got %d, expected %d", tt.letter, m[tt.letter][0].Value, tt.value)
			}
		})
	}
}
