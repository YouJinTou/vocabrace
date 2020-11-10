package scrabble

import "testing"

func TestPut(t *testing.T) {
	b := NewBag(English)
	originalCount := b.Count()
	toPut := []*Tile{&Tile{Letter: "ZZ"}, &Tile{Letter: "QQ"}, &Tile{Letter: "FF"}}

	b.Draw(len(toPut))
	b.Put(toPut)

	if b.Count() != originalCount {
		t.Errorf("Draw/Put count mismatch.")
	}

	for tp := 0; tp < len(toPut); tp++ {
		found := false
		for ti := 0; ti < b.Count(); ti++ {
			if b.Tiles[ti].Letter == toPut[tp].Letter {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Tile not replaced.")
		}
	}
}
