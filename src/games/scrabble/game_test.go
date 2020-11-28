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

	if _, err := g.Exchange(toExchange); err != nil {
		t.Errorf(err.Error())
	}
}

func TestPlace(t *testing.T) {
	g, _, tiles := setupGame(2)

	if _, err := g.Place(tiles); err != nil {
		t.Errorf(err.Error())
	}
}

func TestPlaceAssignsCorrectValueToBlanks(t *testing.T) {
	g, _, w := setupGame(2)
	w.Cells[0].Tile.Value = 0
	w.Cells[0].Tile.Letter = "a"
	g.ToMove().Tiles.GetAt(2).Value = 0
	g.Place(w)

	if g.Board.GetAt(w.Cells[0].Index).Tile.Letter != w.Cells[0].Tile.Letter {
		t.Errorf("expected blank to have been replaced")
	}
}

func TestPlaceReturnsErrorOnInvalidTileIndices(t *testing.T) {
	g, _, w := setupGame(2)
	idx := "qqqqqqq"
	w.Cells[0].Tile.ID = idx
	_, err := g.Place(w)

	if err == nil || err.Error() != fmt.Sprintf("tile with ID %s not found", idx) {
		t.Errorf("passed invalid tile ID, expected error")
	}
}

func TestPlaceSetsBoard(t *testing.T) {
	g, _, w := setupGame(2)

	g.Place(w)

	if len(g.Board.Cells) != w.Length() {
		t.Errorf("Board not set.")
	}
}

func TestPlaceAwardsPoints(t *testing.T) {
	g, _, tiles := setupGame(2)

	g.Place(tiles)

	if g.LastToMove().Points <= 0 {
		t.Errorf("Points not awarded.")
	}
}

func TestPlaceRemovesTilesFromBag(t *testing.T) {
	g, _, w := setupGame(2)
	bagStartingCount := g.Bag.Count()

	g.Place(w)

	if g.Bag.Count() != bagStartingCount-w.Length() {
		t.Errorf("Bag untouched.")
	}
}

func TestPlaceGivesTilesBackToPlayer(t *testing.T) {
	g, _, tiles := setupGame(2)

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
	g, _, tiles := setupGame(2)

	g.Place(tiles)

	if g.delta.LastAction != "Place" {
		t.Errorf("Delta not set.")
	}
}

func TestDeltaStateContainsPlayerPoints(t *testing.T) {
	g, _, tiles := setupGame(2)

	g.Place(tiles)

	delta := g.GetDelta()
	for name, p := range delta.Points {
		if g.GetPlayerByName(name).Points != p {
			t.Errorf("invalid points")
		}
	}
}

func TestPlaceSetsNextPlayer(t *testing.T) {
	g, _, tiles := setupGame(2)
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

func Test_AllLettersNotDrawn_NotOver(t *testing.T) {
	g, _, _ := setupGame(2)
	g.Bag.Draw(50)
	d := g.GetDelta()

	if d.WinnerID != nil {
		t.Errorf("should not be over")
	}
}

func Test_BagEmpty_PlayersStillHaveTiles_NotAllPassedTwice_NotOver(t *testing.T) {
	g, _, _ := setupGame(2)
	g.Bag.Draw(100)
	d := g.GetDelta()

	if d.WinnerID != nil {
		t.Errorf("should not be over")
	}
}

func Test_BagEmpty_PlayerExhaustsTiles_Over(t *testing.T) {
	g, _, _ := setupGame(2)
	g.Bag.Draw(100)
	w := testCreateWord(BoardOrigin, true, g.ToMove().Tiles.Value...)
	result, _ := g.Place(w)
	d := result.GetDelta()

	if d.WinnerID == nil {
		t.Errorf("should be over")
	}
}

func Test_BagEmpty_PlayerExhaustsTiles_AddsOtherPlayersTilesSumToLastPlaced(t *testing.T) {
	for i := 0; i < 50; i++ {
		g, _, _ := setupGame(3)
		g.Bag.Draw(100)
		w := testCreateWord(BoardOrigin, true, g.ToMove().Tiles.Value...)
		toMovePointsBeforePlace := g.ToMove().Points
		finalWordPoints := CalculatePoints(w, []*Word{w})
		sumOfAllOpponentTiles := 0
		for _, p := range g.Players {
			if p.ID != g.ToMoveID {
				sumOfAllOpponentTiles += p.Tiles.Sum()
			}
		}
		expectedFinal := toMovePointsBeforePlace + finalWordPoints + sumOfAllOpponentTiles

		result, _ := g.Place(w)
		d := result.GetDelta()
		actual := d.Points[result.ToMove().Name]

		if actual != expectedFinal {
			t.Errorf("expected %d, got %d", expectedFinal, actual)
			break
		}
	}
}

func Test_BagEmpty_PlayerExhaustsTiles_SubtractsTilesSumFromOtherPlayers(t *testing.T) {
	for i := 0; i < 50; i++ {
		g, _, _ := setupGame(3)
		g.Bag.Draw(100)
		w := testCreateWord(BoardOrigin, true, g.ToMove().Tiles.Value...)
		beforePlace := g.playerPoints()
		result, _ := g.Place(w)
		afterPlace := result.playerPoints()

		for name, points := range beforePlace {
			player := g.GetPlayerByName(name)
			leader := g.Leader()
			if leader == player {
				continue
			}

			actual := afterPlace[name]
			expected := points - player.Tiles.Sum()
			if expected != actual {
				t.Errorf("expected %d, got %d", expected, actual)
				return
			}
		}
	}
}

func Test_AllPassedTwice_Over(t *testing.T) {
	g, _, _ := setupGame(3)
	for i := 0; i < 2*len(g.Players); i++ {
		g.Pass()
	}

	d := g.GetDelta()
	if d.WinnerID == nil {
		t.Errorf("should be over")
	}
}

func Test_AllPassedButLast_NotOver(t *testing.T) {
	g, _, _ := setupGame(3)
	for i := 0; i < 2*len(g.Players)-1; i++ {
		g.Pass()
	}

	d := g.GetDelta()
	if d.WinnerID != nil {
		t.Errorf("should not be over")
	}
}

func Test_AllExchangedTwice_Over(t *testing.T) {
	g, _, _ := setupGame(3)
	for i := 0; i < 2*len(g.Players); i++ {
		g.Exchange([]string{g.ToMove().Tiles.GetAt(0).ID})
	}

	d := g.GetDelta()
	if d.WinnerID == nil {
		t.Errorf("should be over")
	}
}

func Test_AllExchangedButLast_NotOver(t *testing.T) {
	g, _, _ := setupGame(3)
	for i := 0; i < 2*len(g.Players)-1; i++ {
		g.Exchange([]string{g.ToMove().Tiles.GetAt(0).ID})
	}

	d := g.GetDelta()
	if d.WinnerID != nil {
		t.Errorf("should not be over")
	}
}

func Test_EndCounterExceeded_SubtractsPlayerTilesSum(t *testing.T) {
	g, _, _ := setupGame(3)
	for i := 0; i < 2*len(g.Players)-1; i++ {
		g.Pass()
	}
	beforePoints := g.playerPoints()

	g.Pass()

	afterPoints := g.playerPoints()
	for name, points := range beforePoints {
		player := g.GetPlayerByName(name)
		actual := afterPoints[name]
		expected := points - player.Tiles.Sum()
		if expected != actual {
			t.Errorf("expected %d, got %d", expected, actual)
			return
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

func setupGame(playerCount int) (Game, []*Player, *Word) {
	players := []*Player{}
	for i := 0; i < playerCount; i++ {
		players = append(players, testPlayer())
	}
	g := NewGame(English, players, v())
	w := testCreateWord(
		BoardOrigin,
		true,
		g.ToMove().Tiles.GetAt(2).Copy(true),
		g.ToMove().Tiles.GetAt(3).Copy(true))
	return *g, players, w
}

func testCreateWord(startIndex int, isAcross bool, t ...*Tile) *Word {
	cells := []*Cell{}
	for i, tile := range t {
		var idx int
		if isAcross {
			idx = startIndex + i
		} else {
			idx = startIndex + BoardHeight
		}
		cells = append(cells, NewCell(tile, idx))
	}
	w := NewWord(cells)
	return w
}
