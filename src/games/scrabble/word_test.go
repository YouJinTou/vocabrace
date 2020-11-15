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
	ws := "horn"
	w := word(ws, 46, true, []int{})
	b.SetCells(w.Cells)
	words := Extract(b, w)

	if len(words) != 1 {
		t.Errorf("expected one word")
		return
	}

	validateExpected(words, []string{"horn"}, t)
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

	new := "farm"
	w2 := word(new, 18, false, []int{48})
	b.SetCells(w2.Cells)

	words := Extract(b, w2)

	if len(words) != 1 {
		t.Errorf("expected one word")
		return
	}

	validateExpected(words, []string{"farm"}, t)
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

	new := "paste"
	w3 := word(new, 76, true, []int{})
	b.SetCells(w3.Cells)

	words := Extract(b, w3)

	if len(words) != 2 {
		t.Errorf("expected two words")
		return
	}

	validateExpected(words, []string{"paste", "farms"}, t)
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

	new := "mob"
	w4 := word(new, 63, true, []int{63})
	b.SetCells(w4.Cells)

	words := Extract(b, w4)

	if len(words) != 3 {
		t.Errorf("expected three words")
		return
	}

	validateExpected(words, []string{"mob", "not", "be"}, t)
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

func validateExpected(words []*Word, expected []string, t *testing.T) {
	for _, w := range words {
		if !tools.ContainsString(expected, w.String()) {
			t.Errorf("expected %q, got %s", expected, w.String())
		}
	}
}
