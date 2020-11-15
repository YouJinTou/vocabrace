package scrabble

import (
	"strings"
	"testing"

	"github.com/YouJinTou/vocabrace/tools"
)

// _ _ _ _ _ _
// _ h o r n _
// _ _ _ _ _ _
func TestExtractOneLoneWord(t *testing.T) {
	b := NewBoard()

	expected := []string{"horn"}
	w := word(expected[0], 46, true, []int{})
	b.SetCells(w.Cells)

	words := Extract(b, w)

	validateExpected(words, expected, b, t)
}

// _ _ _ _ _ _
// _ h o r n s
// _ _ _ _ _ _
func TestExtractOneExtendingAcross(t *testing.T) {
	b := NewBoard()

	w := word("horn", 46, true, []int{})
	b.SetCells(w.Cells)

	expected := []string{"horns"}
	w2 := word(expected[0], 46, true, []int{46, 47, 48, 49})
	b.SetCells(w2.Cells)

	words := Extract(b, w2)

	validateExpected(words, expected, b, t)
}

// _ _ _ _ _ _
// _ _ _ f _ _
// _ _ _ a _ _
// _ h o r n _
// _ _ _ m _ _
func TestExtractOneCrossWord(t *testing.T) {
	b := NewBoard()

	w1 := word("horn", 46, true, []int{})
	b.SetCells(w1.Cells)

	expected := []string{"farm"}
	w2 := word(expected[0], 18, false, []int{48})
	b.SetCells(w2.Cells)

	words := Extract(b, w2)

	validateExpected(words, expected, b, t)
}

// _ _ _ _ _ _
// _ _ _ f _ _
// _ _ _ a _ _
// _ h o r n _
// _ _ _ m _ _
// _ p a s t e
func TestExtractTwoWordsInverseT(t *testing.T) {
	b := NewBoard()

	w1 := word("horn", 46, true, []int{})
	b.SetCells(w1.Cells)
	w2 := word("farm", 18, false, []int{48})
	b.SetCells(w2.Cells)

	expected := []string{"paste", "farms"}
	w3 := word(expected[0], 76, true, []int{})
	b.SetCells(w3.Cells)

	words := Extract(b, w3)

	validateExpected(words, expected, b, t)
}

// _ _ _ _ _ _
// _ _ _ f _ _
// _ _ _ a _ _
// _ h o r n _
// _ _ _ m o b
// _ p a s t e
func TestExtractThreeWordsSandwich(t *testing.T) {
	b := NewBoard()

	w1 := word("horn", 46, true, []int{})
	b.SetCells(w1.Cells)
	w2 := word("farm", 18, false, []int{48})
	b.SetCells(w2.Cells)
	w3 := word("paste", 76, true, []int{})
	b.SetCells(w3.Cells)

	expected := []string{"mob", "not", "be"}
	w4 := word(expected[0], 63, true, []int{63})
	b.SetCells(w4.Cells)

	words := Extract(b, w4)

	validateExpected(words, expected, b, t)
}

// _ _ _ _ _ _
// _ _ _ f _ _
// _ _ _ a _ _
// _ h o r n _
// _ _ _ m o b
// _ p a s t e
// b i t _ _ _
func TestExtractThreeWordsSandwich2(t *testing.T) {
	b := NewBoard()

	w1 := word("horn", 46, true, []int{})
	b.SetCells(w1.Cells)
	w2 := word("farm", 18, false, []int{48})
	b.SetCells(w2.Cells)
	w3 := word("paste", 76, true, []int{})
	b.SetCells(w3.Cells)
	w4 := word("mob", 63, true, []int{63})
	b.SetCells(w4.Cells)

	expected := []string{"bit", "pi", "at"}
	w5 := word(expected[0], 90, true, []int{})
	b.SetCells(w5.Cells)

	words := Extract(b, w5)

	validateExpected(words, expected, b, t)
}

// _ _ _ _ _ _
// _ _ _ f _ _
// _ _ _ a _ _
// _ h o r n _
// _ _ k m o b
// _ p a s t e
func TestExtractFourWordsSandwich(t *testing.T) {
	b := NewBoard()

	w1 := word("horn", 46, true, []int{})
	b.SetCells(w1.Cells)
	w2 := word("farm", 18, false, []int{48})
	b.SetCells(w2.Cells)
	w3 := word("paste", 76, true, []int{})
	b.SetCells(w3.Cells)

	expected := []string{"kmob", "oka", "not", "be"}
	w4 := word(expected[0], 62, true, []int{63})
	b.SetCells(w4.Cells)

	words := Extract(b, w4)

	validateExpected(words, expected, b, t)
}

// _ _ _ _ _
// _ t _ _ _
// _ e _ _ _
// _ s _ _ _
// _ t _ _ _
func TestExtractOneLoneWordDown(t *testing.T) {
	b := NewBoard()

	expected := []string{"test"}
	w := word(expected[0], 16, false, []int{})
	b.SetCells(w.Cells)
	words := Extract(b, w)

	validateExpected(words, expected, b, t)
}

// _ _ _ _ _
// _ t _ _ _
// _ e a s y
// _ s _ _ _
// _ t _ _ _
func TestExtractOneCrosswordAcross(t *testing.T) {
	b := NewBoard()

	w := word("test", 16, false, []int{})
	b.SetCells(w.Cells)

	expected := []string{"easy"}
	w2 := word(expected[0], 31, true, []int{31})
	b.SetCells(w2.Cells)

	words := Extract(b, w2)

	validateExpected(words, expected, b, t)
}

// _ _ _ _ _
// _ t _ _ _
// _ e _ _ _
// _ s _ _ _
// _ t _ _ _
// _ s _ _ _
func TestExtractOneExtendingDown(t *testing.T) {
	b := NewBoard()

	w := word("test", 16, false, []int{})
	b.SetCells(w.Cells)

	expected := []string{"tests"}
	w2 := word(expected[0], 16, false, []int{16, 31, 46, 61})
	b.SetCells(w2.Cells)

	words := Extract(b, w2)

	validateExpected(words, expected, b, t)
}

// o _ _ _ _
// g a s _ _
// l _ _ _ _
// e _ _ _ _
func TestExtractTwoWordsInverseTDown(t *testing.T) {
	b := NewBoard()

	w := word("as", 16, true, []int{})
	b.SetCells(w.Cells)

	expected := []string{"ogle", "gas"}
	w2 := word(expected[0], 0, false, []int{})
	b.SetCells(w2.Cells)

	words := Extract(b, w2)

	validateExpected(words, expected, b, t)
}

// _ _ _ _ _ _ _
// p a s s a g e
// _ t o o t _ _
func TestExtractSandwichAcross(t *testing.T) {
	b := NewBoard()

	w := word("passage", 15, true, []int{})
	b.SetCells(w.Cells)

	expected := []string{"toot", "at", "so", "so", "at"}
	w2 := word(expected[0], 31, true, []int{})
	b.SetCells(w2.Cells)

	words := Extract(b, w2)

	validateExpected(words, expected, b, t)
}

// _ d o n
// _ o l a
// s o d _
// _ d i m
// _ l e e
// b e s t
func TestExtractMassiveSandwichDown(t *testing.T) {
	b := NewBoard()

	w := word("doodle", 1, false, []int{})
	b.SetCells(w.Cells)

	w2 := word("so", 30, true, []int{31})
	b.SetCells(w2.Cells)

	w3 := word("be", 75, true, []int{76})
	b.SetCells(w3.Cells)

	w4 := word("met", 48, false, []int{})
	b.SetCells(w4.Cells)

	w5 := word("na", 3, false, []int{})
	b.SetCells(w5.Cells)

	expected := []string{"oldies", "sod", "dim", "lee", "best", "don", "ola"}
	w6 := word(expected[0], 2, false, []int{})
	b.SetCells(w6.Cells)

	words := Extract(b, w6)

	validateExpected(words, expected, b, t)
}

func word(word string, startIndex int, isAcross bool, skipAt []int) *Word {
	tokens := strings.Split(word, "")
	cells := []*Cell{}
	idx := startIndex
	for _, t := range tokens {
		if tools.ContainsInt(skipAt, idx) {
			idx = incrementWordIndex(idx, isAcross)
			continue
		}

		cells = append(cells, NewCell(NewTile(t, 1), idx))
		idx = incrementWordIndex(idx, isAcross)
		if idx > BoardMaxIndex || idx < BoardMinIndex {
			panic("invalid testing index")
		}
	}
	return NewWord(cells)
}

func incrementWordIndex(current int, isAcross bool) int {
	if isAcross {
		return current + 1
	}
	return current + BoardHeight
}

func validateExpected(words []*Word, expected []string, b *Board, t *testing.T) {
	if len(words) != len(expected) {
		t.Errorf("expected %d words, got %d", len(expected), len(words))

		for _, e := range expected {
			t.Errorf("Expected: %s", e)
		}

		for _, w := range words {
			t.Errorf("Got: %s", w.String())
		}

		b.Print()

		return
	}

	expectedCounts := make(map[string]int)
	for _, e := range expected {
		if _, ok := expectedCounts[e]; ok {
			expectedCounts[e]++
		} else {
			expectedCounts[e] = 1
		}
	}

	wordCounts := make(map[string]int)
	for _, w := range words {
		if _, ok := wordCounts[w.String()]; ok {
			wordCounts[w.String()]++
		} else {
			wordCounts[w.String()] = 1
		}
	}

	for _, e := range expected {
		expectedCount, _ := expectedCounts[e]
		if receivedCount, ok := wordCounts[e]; ok {
			if expectedCount != receivedCount {
				t.Errorf("expected %d, received %d", expectedCount, receivedCount)
			}
		} else {
			t.Errorf("expected to find %s", e)
		}
	}
}
