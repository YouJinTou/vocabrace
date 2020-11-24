package scrabble

import (
	"fmt"
	"strconv"
	"testing"
)

const English = "english"

func TestExchange(t *testing.T) {
	players := []*Player{testPlayer(), testPlayer()}
	g := NewGame(English, players, v())
	toExchange := []string{g.ToMove().Tiles.GetAt(0).ID, g.ToMove().Tiles.GetAt(1).ID}
	_, err := g.Exchange(toExchange)

	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestPlace(t *testing.T) {
	g, _, tiles := setupPlace()
	_, err := g.Place(tiles)

	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestPlaceAssignsCorrectValueToBlanks(t *testing.T) {
	g, _, w := setupPlace()
	w.Cells[0].Tile.Value = 0
	w.Cells[0].Tile.Letter = "a"
	g.ToMove().Tiles.GetAt(2).Value = 0
	g.Place(w)

	if g.Board.GetAt(w.Cells[0].Index).Tile.Letter != w.Cells[0].Tile.Letter {
		t.Errorf("expected blank to have been replaced")
	}
}

func TestPlaceReturnsErrorOnInvalidTileIndices(t *testing.T) {
	g, _, w := setupPlace()
	idx := "qqqqqqq"
	w.Cells[0].Tile.ID = idx
	_, err := g.Place(w)

	if err == nil || err.Error() != fmt.Sprintf("tile with ID %s not found", idx) {
		t.Errorf("passed invalid tile ID, expected error")
	}
}

func TestPlaceSetsBoard(t *testing.T) {
	g, _, w := setupPlace()

	g.Place(w)

	if len(g.Board.Cells) != w.Length() {
		t.Errorf("Board not set.")
	}
}

func TestPlaceAwardsPoints(t *testing.T) {
	g, _, tiles := setupPlace()

	g.Place(tiles)

	if g.LastToMove().Points <= 0 {
		t.Errorf("Points not awarded.")
	}
}

func TestPlaceRemovesTilesFromBag(t *testing.T) {
	g, _, w := setupPlace()
	bagStartingCount := g.Bag.Count()

	g.Place(w)

	if g.Bag.Count() != bagStartingCount-w.Length() {
		t.Errorf("Bag untouched.")
	}
}

func TestPlaceGivesTilesBackToPlayer(t *testing.T) {
	g, _, tiles := setupPlace()

	g.Place(tiles)

	for _, bt := range g.Bag.GetLastDrawn().Value {
		found := false
		for _, pt := range g.LastToMove().Tiles.Value {
			if bt.ID == pt.ID {
				found = true
			}
		}
		if !found {
			t.Errorf("Invalid tile assigned.")
		}
	}
}

func TestPlaceSetsDeltaState(t *testing.T) {
	g, _, tiles := setupPlace()

	g.Place(tiles)

	if g.delta.LastAction != "Place" {
		t.Errorf("Delta not set.")
	}
}

func TestPlaceSetsNextPlayer(t *testing.T) {
	g, _, tiles := setupPlace()
	previousToMode := g.ToMoveID

	g.Place(tiles)

	if g.ToMoveID == previousToMode {
		t.Errorf("Next player not set.")
	}
}

func TestOrderPlayers(t *testing.T) {
	for total := 1; total < 50; total++ {
		players := []*Player{}
		for i := 0; i < total; i++ {
			players = append(players, testPlayerArgs(strconv.Itoa(i)))
		}

		for x := 0; x < 50; x++ {
			g := NewGame(English, players, v())
			io, _ := strconv.Atoi(g.Order[0])
			expected := getExpectedOrder(io, total)

			if len(expected) != len(g.Order) {
				t.Errorf("Lengths mismatch.")
			}

			for i := 0; i < len(g.Order); i++ {
				if expected[i] != g.Order[i] {
					t.Errorf("Invalid order.")
				}
			}
		}
	}
}

func TestSetNext(t *testing.T) {
	p1 := testPlayerArgs("1")
	p2 := testPlayerArgs("2")
	p3 := testPlayerArgs("3")
	players := []*Player{p1, p2, p3}

	for j := 0; j < 10; j++ {
		g := NewGame(English, players, v())

		for i := 0; i < 50; i++ {
			toMoveID := g.ToMoveID
			g.setNext()
			failed := false

			if toMoveID == p1.ID {
				if g.ToMoveID != p2.ID {
					failed = true
				}
			} else if toMoveID == p2.ID {
				if g.ToMoveID != p3.ID {
					failed = true
				}
			} else {
				if g.ToMoveID != p1.ID {
					failed = true
				}
			}

			if failed {
				t.Errorf("could not set next player")
			}
		}
	}
}

func getExpectedOrder(idx, total int) []string {
	result := []string{strconv.Itoa(idx)}
	for i := idx + 1; i < total; i++ {
		result = append(result, strconv.Itoa(i))
	}
	for i := 0; i < idx; i++ {
		result = append(result, strconv.Itoa(i))
	}
	return result
}

func setupPlace() (Game, []*Player, *Word) {
	players := []*Player{testPlayer(), testPlayer()}
	g := NewGame(English, players, v())
	cells := []*Cell{
		&Cell{
			Tile:  *g.ToMove().Tiles.GetAt(2).Copy(true),
			Index: BoardOrigin,
		},
		&Cell{
			Tile:  *g.ToMove().Tiles.GetAt(3).Copy(true),
			Index: BoardOrigin + 1,
		},
	}
	return *g, players, NewWord(cells)
}
