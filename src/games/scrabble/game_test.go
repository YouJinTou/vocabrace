package scrabble

import (
	"strconv"
	"testing"
)

func TestExchange(t *testing.T) {
	players := []*Player{testPlayer(), testPlayer()}
	g := NewGame(players)
	toExchange := []string{g.ToMove().Tiles[0].Index, g.ToMove().Tiles[1].Index}
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

func TestPlaceBoardSet(t *testing.T) {
	g, _, tiles := setupPlace()

	g.Place(tiles)

	if len(g.Board.Cells) != len(tiles) {
		t.Errorf("Board not set.")
	}
}

func TestPlacePointsAwarded(t *testing.T) {
	g, _, tiles := setupPlace()

	g.Place(tiles)

	if g.Players[0].Points <= 0 {
		t.Errorf("Points not awarded.")
	}
}

func TestPlaceTilesRemovedFromBag(t *testing.T) {
	g, _, tiles := setupPlace()
	bagStartingCount := g.Bag.Count()

	g.Place(tiles)

	if g.Bag.Count() != bagStartingCount-len(tiles) {
		t.Errorf("Bag untouched.")
	}
}

func TestPlacePlayerReceivesTilesBack(t *testing.T) {
	g, _, tiles := setupPlace()

	g.Place(tiles)

	for _, bt := range g.Bag.GetLastDrawn() {
		found := false
		for _, pt := range g.LastToMove().Tiles {
			if bt.Index == pt.Index {
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

	if g.ToMove().ID == previousToMode {
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
			g := NewGame(players)
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
		g := NewGame(players)

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

func setupPlace() (Game, []*Player, []*Cell) {
	players := []*Player{testPlayer(), testPlayer()}
	g := NewGame(players)
	tiles := []*Cell{
		&Cell{
			Tile:  *g.ToMove().Tiles[2].Copy(true),
			Index: 0,
		},
		&Cell{
			Tile:  *g.ToMove().Tiles[5].Copy(true),
			Index: 1,
		},
	}
	return *g, players, tiles
}
