package scrabble

import (
	"strconv"
	"testing"
)

var counts = []int{0, 1, 2, 100, 1000}
var puts = []*Tile{
	&Tile{Letter: "ZZ", ID: "abc"},
	&Tile{Letter: "QQ", ID: "qqf"},
	&Tile{Letter: "FF", ID: "ffx"},
}

func TestVerifyUniqueTiles(t *testing.T) {
	b := NewBag(English)
	m := map[string]int{}
	for _, t := range b.Tiles {
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
			for j, bt := range b.Tiles {
				for _, d := range drawn {
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

	b.Draw(len(puts))
	b.Put(puts)

	if b.Count() != originalCount {
		t.Errorf("Draw/Put count mismatch.")
	}
}

func TestPutAddsTiles(t *testing.T) {
	b := NewBag(English)
	originalCount := b.Count()
	expected := originalCount + len(puts)
	b.Put(puts)

	if expected != b.Count() {
		t.Errorf("expected count %d, got %d", expected, b.Count())
	}

	for _, p := range puts {
		found := false
		for _, bt := range b.Tiles {
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
			result := originalCount - len(drawn)
			if result != b.Count() {
				t.Errorf("got %q, want %q", result, len(drawn))
			}
		})
	}
}
