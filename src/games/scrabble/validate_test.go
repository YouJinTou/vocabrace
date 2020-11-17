package scrabble

import (
	"fmt"
	"testing"
)

func TestValidatePlace_NoTiles_ReturnsError(t *testing.T) {
	err := v().ValidatePlace(testValidatorGame(), &Word{})
	if err == nil || err.Error() != "at least one tile required to form a word" {
		t.Errorf("expected an error")
	}
}

func TestValidatePlace_IndicesOutsideBounds_ReturnsError(t *testing.T) {
	invalidIndices := []int{-1000, -1, 225, 1000}
	vs := fmt.Sprintf("valid indices between %d and %d", BoardMinIndex, BoardMaxIndex)
	for _, i := range invalidIndices {
		w := NewWord([]*Cell{NewCell(BlankTile(), i)})
		err := v().ValidatePlace(testValidatorGame(), w)
		if err == nil || err.Error() != vs {
			t.Errorf("expected an error")
		}
	}
}

func TestValidatePlace_TilesExistOnIndices_ReturnsError(t *testing.T) {
	t.Run("old across, new across 1", testOverlap("test", "at", 16, 15, true, true, true))
	t.Run("old across, new across 2", testOverlap("test", "true", 16, 17, true, true, true))
	t.Run("old across, new across 3", testOverlap("test", "true", 16, 18, true, true, true))
	t.Run("old across, new across 4", testOverlap("test", "true", 16, 19, true, true, true))

	t.Run("old across, new down 1", testOverlap("test", "true", 67, 67, true, false, true))
	t.Run("old across, new down 2", testOverlap("test", "true", 67, 52, true, false, true))
	t.Run("old across, new down 3", testOverlap("test", "true", 67, 37, true, false, true))
	t.Run("old across, new down 4", testOverlap("test", "true", 67, 22, true, false, true))

	t.Run("old down, new across 1", testOverlap("test", "true", 0, 0, false, true, true))
	t.Run("old down, new across 2", testOverlap("test", "easy", 0, 15, false, true, true))
	t.Run("old down, new across 3", testOverlap("test", "stem", 0, 30, false, true, true))
	t.Run("old down, new across 4", testOverlap("test", "true", 0, 45, false, true, true))

	t.Run("old down, new down 1", testOverlap("test", "true", 7, 7, false, false, true))
	t.Run("old down, new down 2", testOverlap("test", "easy", 22, 7, false, false, true))
	t.Run("old down, new down 3", testOverlap("test", "stem", 37, 7, false, false, true))
	t.Run("old down, new down 4", testOverlap("test", "true", 52, 7, false, false, true))
}

func TestValidatePlace_TilesDoNotExistOnInidces_DoesNotReturnError(t *testing.T) {
	t.Run("old across, new across 1", testOverlap("test", "at", 17, 15, true, true, false))
	t.Run("old across, new across 2", testOverlap("test", "true", 21, 17, true, true, false))
	t.Run("old across, new across 3", testOverlap("test", "true", 22, 18, true, true, false))
	t.Run("old across, new across 4", testOverlap("test", "true", 23, 19, true, true, false))

	t.Run("old across, new down 1", testOverlap("test", "true", 67, 82, true, false, false))
	t.Run("old across, new down 2", testOverlap("test", "true", 67, 83, true, false, false))
	t.Run("old across, new down 3", testOverlap("test", "true", 67, 84, true, false, false))
	t.Run("old across, new down 4", testOverlap("test", "true", 67, 85, true, false, false))

	t.Run("old down, new across 1", testOverlap("test", "true", 50, 46, false, true, false))
	t.Run("old down, new across 2", testOverlap("test", "easy", 50, 61, false, true, false))
	t.Run("old down, new across 3", testOverlap("test", "stem", 50, 81, false, true, false))
	t.Run("old down, new across 4", testOverlap("test", "true", 50, 91, false, true, false))

	t.Run("old down, new down 1", testOverlap("test", "true", 7, 6, false, false, false))
	t.Run("old down, new down 2", testOverlap("test", "easy", 22, 23, false, false, false))
	t.Run("old down, new down 3", testOverlap("test", "stem", 37, 40, false, false, false))
	t.Run("old down, new down 4", testOverlap("test", "true", 52, 51, false, false, false))
}

func TestValidatePlace_PassesWhenPlayerHasCorrectTiles(t *testing.T) {
	g := testValidatorGame()
	cells := []*Cell{}
	for _, t := range g.ToMove().Tiles.Value {
		cells = append(cells, NewCell(t, 0))
	}
	w := NewWord(cells)
	err := v().ValidatePlace(g, w)
	if err != nil {
		t.Errorf("did not expect error")
	}
}

func TestValidatePlace_FailsWhenPlayerHasIncorrectTiles(t *testing.T) {
	g := testValidatorGame()
	blank := BlankTile()
	cells := []*Cell{NewCell(blank, 0)}
	w := NewWord(cells)
	err := v().ValidatePlace(g, w)
	if err == nil || err.Error() != fmt.Sprintf("tile with ID %s not found", blank.ID) {
		t.Errorf("expected error")
	}
}

func testOverlap(
	existing, new string,
	existingStart, newStart int,
	isAcrossStart, isAcrossNew, assertError bool) func(*testing.T) {
	return func(t *testing.T) {
		g := testValidatorGame()
		occupied := word(existing, existingStart, isAcrossStart, []int{}, []int{}, []int{})
		g.Board.SetCells(occupied.Cells)

		toPlace := word(new, newStart, isAcrossNew, []int{}, []int{}, []int{})
		err := v().ValidatePlace(g, toPlace)

		if assertError {
			if err == nil || err.Error() != "index overlap" {
				t.Errorf("expected an error")
			}
		} else {
			if err != nil && err.Error() == "index overlap" {
				t.Errorf("did not expect an error")
			}
		}
	}
}

func v() *Validator {
	return NewValidator()
}

func testValidatorGame() Game {
	players := []*Player{testPlayer(), testPlayer()}
	return *NewGame(players)
}
