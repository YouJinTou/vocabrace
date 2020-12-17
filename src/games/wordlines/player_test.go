package wordlines

import (
	"testing"

	"github.com/YouJinTou/vocabrace/tools"

	"github.com/google/uuid"
)

var toReceive = NewTiles(&Tile{Letter: "Q"}, &Tile{Letter: "X"})

func TestExchangeTilesReturnsErrorAboutMissingTiles(t *testing.T) {
	p := testPlayer()
	_, err := p.ExchangeTiles([]string{"ZZ"}, NewTiles(&Tile{}))

	if err == nil || err.Error() != "ZZ tile not found" {
		t.Errorf("Should return error about missing tile.")
	}
}

func TestExchangeTilesCountsMatch(t *testing.T) {
	p := testPlayer()
	toRemove := []string{p.Tiles.GetAt(1).ID, p.Tiles.GetAt(5).ID}
	tiles, _ := p.ExchangeTiles(toRemove, toReceive)

	if tiles.Count() != toReceive.Count() {
		t.Errorf("Received %d tile(s) instead of %d", tiles.Count(), toReceive.Count())
	}
}

func TestExchangeTilesReturnTilesMatch(t *testing.T) {
	p := testPlayer()
	toRemove := []string{p.Tiles.GetAt(1).ID, p.Tiles.GetAt(5).ID}
	tiles, _ := p.ExchangeTiles(toRemove, toReceive)

	for i := 0; i < tiles.Count(); i++ {
		if toRemove[i] != tiles.GetAt(i).ID {
			t.Errorf("Invalid return tile. Expected %s, got %s.", toRemove[i], tiles.GetAt(i).ID)
		}
	}
}

func TestExchangeTilesBeforeAfterCountsMatch(t *testing.T) {
	p := testPlayer()
	originalTiles := NewTiles(p.Tiles.Value...)
	toRemove := []string{p.Tiles.GetAt(1).ID, p.Tiles.GetAt(5).ID}

	p.ExchangeTiles(toRemove, toReceive)

	if p.Tiles.Count() != originalTiles.Count() {
		t.Errorf("Before/after count mismatch.")
	}
}

func TestExchangeTilesPlayerReceivesTiles(t *testing.T) {
	for xx := 0; xx < 10000; xx++ {
		p := testPlayer()
		originalTiles := NewTiles(p.Tiles.Value...)
		toRemove := []string{p.Tiles.GetAt(1).ID, p.Tiles.GetAt(5).ID}

		p.ExchangeTiles(toRemove, toReceive)

		newTiles := NewTiles()
		for _, t := range originalTiles.Value {
			if !tools.ContainsString(toRemove, t.ID) {
				newTiles.Append(t)
			}
		}
		newTiles.Append(toReceive.Value...)

		for _, nt := range newTiles.Value {
			if fbi := p.Tiles.FindByID(nt.ID); fbi == nil {
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
