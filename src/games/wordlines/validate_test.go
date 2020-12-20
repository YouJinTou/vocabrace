package wordlines

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

type wordCheckerMock struct{}

var validateWordsMock func(a string, b []string) ([]string, error)

func (w wordCheckerMock) ValidateWords(language string, words []string) ([]string, error) {
	if validateWordsMock == nil {
		return []string{}, nil
	}
	return validateWordsMock(language, words)
}

func TestValidatePlace_StartDifferentFromOrigin_ReturnsError(t *testing.T) {
	g := testValidatorGame()
	err := v().ValidatePlace(g, word("test", 0, false, []int{}, []int{}, []int{}))
	if err == nil || err.Error() != "first place must cross the origin" {
		t.Errorf("expected error")
	}
}

func TestValidatePlace_WordNotHorizontal_ReturnsError(t *testing.T) {
	g := testValidatorGame()
	w1 := word("test", BoardOrigin, true, []int{}, []int{}, []int{})
	g.Board.SetCells(w1.Cells)

	w2 := playerWord(g, []int{BoardOrigin + 30, BoardOrigin + 15 + 1, BoardOrigin + 15 + 2})
	err := v().ValidatePlace(g, w2)

	if err == nil || err.Error() != "word must be straight" {
		t.Errorf("expected error")
	}
}

func TestValidatePlace_WordNotVertical_ReturnsError(t *testing.T) {
	g := testValidatorGame()
	w1 := word("test", BoardOrigin, false, []int{}, []int{}, []int{})
	g.Board.SetCells(w1.Cells)

	w2 := playerWord(
		g,
		[]int{BoardOrigin + 1, BoardOrigin + 1 + BoardHeight, BoardOrigin + 1 + BoardHeight + 1})
	err := v().ValidatePlace(g, w2)

	if err == nil || err.Error() != "word must be straight" {
		t.Errorf("expected error")
	}
}

// X _ _ _
// t e s t
// X _ _ _
func TestValidatePlace_WordVertical_HasGaps_DoesNotReturnError1(t *testing.T) {
	g := testValidatorGame()
	w1 := word("test", BoardOrigin, true, []int{}, []int{}, []int{})
	g.Board.SetCells(w1.Cells)

	w2 := playerWord(g, []int{BoardOrigin - BoardHeight, BoardOrigin + BoardHeight})

	if err := v().ValidatePlace(g, w2); err != nil {
		t.Errorf(err.Error())
	}
}

// X _ _ _
// t e s t
// o _ _ _
// X _ _ _
func TestValidatePlace_WordVertical_HasGaps_DoesNotReturnError2(t *testing.T) {
	g := testValidatorGame()
	w1 := word("test", BoardOrigin, true, []int{}, []int{}, []int{})
	g.Board.SetCells(w1.Cells)
	w2 := word("to", BoardOrigin, false, []int{BoardOrigin}, []int{}, []int{})
	g.Board.SetCells(w2.Cells)
	w3 := playerWord(g, []int{BoardOrigin - BoardHeight, BoardOrigin + 2*BoardHeight})

	if err := v().ValidatePlace(g, w3); err != nil {
		t.Errorf(err.Error())
	}
}

// X _ _ _
// t e s t
// o _ _ _
// _ _ _ _
// X _ _ _
func TestValidatePlace_WordVertical_HasInvalidGaps_ReturnsError(t *testing.T) {
	g := testValidatorGame()
	w1 := word("test", BoardOrigin, true, []int{}, []int{}, []int{})
	g.Board.SetCells(w1.Cells)
	w2 := word("to", BoardOrigin, false, []int{BoardOrigin}, []int{}, []int{})
	g.Board.SetCells(w2.Cells)

	w3 := playerWord(g, []int{BoardOrigin - BoardHeight, BoardOrigin + 3*BoardHeight})
	err := v().ValidatePlace(g, w3)

	if err == nil || err.Error() != "word must be straight" {
		t.Errorf("expected error")
	}
}

// X t e s t X
func TestValidatePlace_WordHorizontal_HasGaps_DoesNotReturnError1(t *testing.T) {
	g := testValidatorGame()
	w1 := word("test", BoardOrigin, true, []int{}, []int{}, []int{})
	g.Board.SetCells(w1.Cells)

	w2 := playerWord(g, []int{BoardOrigin + 4, BoardOrigin - 1})

	if err := v().ValidatePlace(g, w2); err != nil {
		t.Errorf(err.Error())
	}
}

// X t e s t X t o X
func TestValidatePlace_WordHorizontal_HasGaps_DoesNotReturnError2(t *testing.T) {
	g := testValidatorGame()
	w1 := word("test", BoardOrigin, true, []int{}, []int{}, []int{})
	g.Board.SetCells(w1.Cells)
	w2 := word("to", BoardOrigin+5, true, []int{}, []int{}, []int{})
	g.Board.SetCells(w2.Cells)

	w3 := playerWord(g, []int{BoardOrigin + 4, BoardOrigin - 1, BoardOrigin + 7})
	err := v().ValidatePlace(g, w3)

	if err != nil {
		t.Errorf(err.Error())
	}
}

// X t e s t X t o _ X
func TestValidatePlace_WordHorizontal_HasInvalidGaps_ReturnsError(t *testing.T) {
	g := testValidatorGame()
	w1 := word("test", BoardOrigin, true, []int{}, []int{}, []int{})
	g.Board.SetCells(w1.Cells)
	w2 := word("to", BoardOrigin, false, []int{BoardOrigin}, []int{}, []int{})
	g.Board.SetCells(w2.Cells)

	w3 := playerWord(g, []int{BoardOrigin + 4, BoardOrigin - 1, BoardOrigin + 8})
	err := v().ValidatePlace(g, w3)

	if err == nil || err.Error() != "word must be straight" {
		t.Errorf("expected error")
	}
}

// _ _ _ _ _
// t e s t _
// _ _ _ _ _
// X X _ _ _

// _ _ _ _ _ _
// t e s t _ _
// _ _ _ _ X X

// X X _ _
// _ _ _ _
// t e s t

// _ _ _ _ _ _
// X X _ _ _ _
// _ _ t e s t
func TestValidatePlace_WordDoesNotTouchAdjacentTiles_ReturnsError(t *testing.T) {
	testMap := [][]int{
		[]int{BoardOrigin + BoardHeight*2, BoardOrigin + BoardHeight*2 + 1},
		[]int{BoardOrigin + BoardHeight + 4, BoardOrigin + BoardHeight + 5},
		[]int{BoardOrigin - BoardHeight*2, BoardOrigin - BoardHeight*2 - 1},
		[]int{BoardOrigin - BoardHeight - 1, BoardOrigin - BoardHeight - 2},
	}
	for _, tm := range testMap {
		g := testValidatorGame()
		w1 := word("test", BoardOrigin, true, []int{}, []int{}, []int{})
		g.Board.SetCells(w1.Cells)

		w2 := playerWord(g, tm)
		err := v().ValidatePlace(g, w2)

		if err == nil || err.Error() != "adjacent tiles not found" {
			t.Errorf("expected error")
		}
	}
}

func TestValidatePlace_SingleLetterWordFailsDuringWordLookup_ReturnsError(t *testing.T) {
	g := testValidatorGame()
	w := playerWord(g, []int{BoardOrigin})
	validateWordsMock = func(a string, b []string) ([]string, error) {
		return []string{"a"}, errors.New("invalid words")
	}
	err := v().ValidatePlace(g, w)
	if err == nil || err.Error() != "invalid words: [\"a\"]; err: invalid words" {
		t.Errorf("expected error")
	}
}

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

func TestValidatePlace_InvalidWordTilesLength_ReturnsError(t *testing.T) {
	word := NewWord([]*Cell{NewCell(NewTile("test", 1), 0)})
	err := v().ValidatePlace(testValidatorGame(), word)
	if err == nil || err.Error() != "tile letter must be of length 1" {
		t.Errorf("expected an error")
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
	w := playerWord(g, []int{BoardOrigin, BoardOrigin + 1})
	validateWordsMock = nil
	if err := v().ValidatePlace(g, w); err != nil {
		t.Errorf(err.Error())
	}
}

func TestValidatePlace_FailsWhenPlayerHasIncorrectTiles(t *testing.T) {
	g := testValidatorGame()
	tile := NewTile("f", 1)
	cells := []*Cell{NewCell(tile, BoardOrigin)}
	w := NewWord(cells)
	err := v().ValidatePlace(g, w)
	if err == nil || err.Error() != fmt.Sprintf("tile with ID %s not found", tile.ID) {
		t.Errorf("expected error")
	}
}

func TestValidatePlace_FailsWhenWordsNotValid(t *testing.T) {
	g := testValidatorGame()
	w := playerWord(g, []int{BoardOrigin, BoardOrigin + 1})
	validateWordsMock = func(a string, b []string) ([]string, error) {
		return []string{"test"}, errors.New("invalid words")
	}
	err := v().ValidatePlace(g, w)
	if err == nil || err.Error() != "invalid words: [\"test\"]; err: invalid words" {
		t.Errorf("expected error")
	}
}

func TestValidatePlace_PassesWhenWordsValid(t *testing.T) {
	g := testValidatorGame()
	w := playerWord(g, []int{BoardOrigin, BoardOrigin + 1})
	validateWordsMock = func(a string, b []string) ([]string, error) { return []string{}, nil }
	err := v().ValidatePlace(g, w)
	if err != nil {
		t.Errorf("did not expect error")
	}
}

// _ T i _
// _ i t _
func Test_RemovesDuplicateKeys(t *testing.T) {
	g := testValidatorGame()
	g.Board.SetCells([]*Cell{
		NewCell(NewTile("i", 0), BoardOrigin+1),
		NewCell(NewTile("t", 0), BoardOrigin+15+1),
		NewCell(NewTile("i", 0), BoardOrigin+15),
	})
	pt := g.ToMove().Tiles.GetAt(0)
	w := NewWord([]*Cell{NewCell(pt, BoardOrigin)})
	validateWordsMock = func(a string, b []string) ([]string, error) { return b, errors.New("") }
	err := v().ValidatePlace(g, w)
	expected := strings.ToLower(fmt.Sprintf("[\"%si\"]", pt.Letter))
	if !strings.Contains(err.Error(), expected) {
		t.Errorf("expected %s, got %s", expected, err.Error())
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

func v() CanValidate {
	return NewValidator(wordCheckerMock{})
}

func testValidatorGame() Game {
	players := []*Player{testPlayer(), testPlayer()}
	g := *NewClassicGame(English, players, v())
	return g
}

func playerWord(g Game, indices []int) *Word {
	t := NewTiles(g.ToMove().Tiles.Value[0:len(indices)]...)
	cells := []*Cell{}
	for i, ti := range t.Value {
		cells = append(cells, NewCell(ti, indices[i]))
	}
	return NewWord(cells)
}
