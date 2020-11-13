package scrabble

import (
	"testing"

	"github.com/YouJinTou/vocabrace/tools"

	"github.com/google/uuid"
)

func TestExchangeTilesReturnsErrorAboutCountsMismatch(t *testing.T) {
	p := testPlayer()
	_, err := p.ExchangeTiles([]string{"A"}, []*Tile{&Tile{}, &Tile{}})

	if err == nil || err.Error() != "exchange and receive tile counts should match" {
		t.Errorf("Should return error about counts mismatch.")
	}
}

func TestExchangeTilesReturnsErrorAboutMissingTiles(t *testing.T) {
	p := testPlayer()
	_, err := p.ExchangeTiles([]string{"ZZ"}, []*Tile{&Tile{}})

	if err == nil || err.Error() != "ZZ tile not found" {
		t.Errorf("Should return error about missing tile.")
	}
}

func TestExchangeTilesExchangesTiles(t *testing.T) {
	for xx := 0; xx < 10000; xx++ {
		p := testPlayer()
		originalTiles := make([]*Tile, len(p.Tiles))
		copy(originalTiles, p.Tiles)
		toRemove := []string{p.Tiles[0].Letter, p.Tiles[1].Letter}
		toReceive := []*Tile{&Tile{Letter: "Q"}, &Tile{Letter: "X"}}
		tiles, _ := p.ExchangeTiles(toRemove, toReceive)

		if len(tiles) != len(toReceive) {
			t.Errorf("Received %d tile(s) instead of %d", len(tiles), len(toReceive))
		}

		for i := 0; i < len(tiles); i++ {
			if toRemove[i] != tiles[i].Letter {
				t.Errorf("Invalid return tile. Expected %s, got %s.", toRemove[i], tiles[i].Letter)
			}
		}

		if len(p.Tiles) != len(originalTiles) {
			t.Errorf("Before/after count mismatch.")
		}

		remainingAfterRemoval := []*Tile{}
		for _, t := range originalTiles {
			if !tools.ContainsString(toRemove, t.Letter) {
				remainingAfterRemoval = append(remainingAfterRemoval, t)
			}
		}
		newTiles := append(remainingAfterRemoval, toReceive...)

		usedIndices := []int{}
		for _, nt := range newTiles {
			found := false

			for i, pt := range p.Tiles {
				if nt.Letter == pt.Letter && !tools.ContainsInt(usedIndices, i) {
					found = true
					usedIndices = append(usedIndices, i)
					break
				}
			}

			if !found {
				t.Errorf("New tile %s not found.", nt.Letter)
			}
		}
	}
}

func testPlayer() *Player {
	b := NewBag(English)
	p := Player{
		ID:     uuid.New().String(),
		Name:   uuid.New().String(),
		Points: 0,
		Tiles:  b.Draw(7),
	}
	return &p
}

func testPlayerArgs(ID string) *Player {
	p := testPlayer()
	p.ID = ID
	return p
}
