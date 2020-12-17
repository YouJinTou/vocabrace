package wordlines

import (
	"errors"
	"fmt"
	"unicode/utf8"

	"github.com/YouJinTou/vocabrace/tools"
)

// CanValidate is the validator interface.
type CanValidate interface {
	ValidatePlace(g Game, w *Word) error
}

// Validator validates a word placement attempt.
type Validator struct {
	wc WordChecker
}

// NewValidator creates a new validator.
func NewValidator(wc WordChecker) CanValidate {
	return &Validator{wc: wc}
}

// NewDynamoValidator creates a new Dynamo-powered validator.
func NewDynamoValidator() CanValidate {
	return NewValidator(NewDynamoChecker())
}

// ValidatePlace validates a place action.
func (v *Validator) ValidatePlace(g Game, w *Word) error {
	if w.Length() == 0 {
		return errors.New("at least one tile required to form a word")
	}

	if err := v.indicesWithinBounds(w.Cells); err != nil {
		return err
	}

	if err := v.tileLettersOfLength1(w.Cells); err != nil {
		return err
	}

	if err := v.indicesOverlap(&g, w.Cells); err != nil {
		return err
	}

	if err := v.firstPlaceAtOrigin(&g, w.Cells); err != nil {
		return err
	}

	if err := v.checkPlayerTiles(&g, w.Cells); err != nil {
		return err
	}

	if err := v.checkStraightWord(&g, w); err != nil {
		return err
	}

	if err := v.checkAdjacentTilesExist(&g, w); err != nil {
		return err
	}

	g.SetCellTiles(w.Cells)

	g.Board.SetCells(w.Cells)

	words := Extract(&g.Board, w)
	keys := tools.ToLowerStrings(ToStrings(words))
	if notFound, err := v.wc.ValidateWords(g.Language, keys); len(notFound) > 0 || err != nil {
		return fmt.Errorf("invalid words: %q; err: %s", notFound, err)
	}

	return nil
}

func (v *Validator) indicesWithinBounds(cells []*Cell) error {
	for _, c := range cells {
		if c.Index < BoardMinIndex || c.Index >= BoardMaxIndex {
			return fmt.Errorf("valid indices between %d and %d", BoardMinIndex, BoardMaxIndex)
		}
	}
	return nil
}

func (v *Validator) tileLettersOfLength1(cells []*Cell) error {
	for _, c := range cells {
		if !c.Tile.IsBlank() && utf8.RuneCountInString(c.Tile.Letter) != 1 {
			return fmt.Errorf("tile letter must be of length 1")
		}
	}
	return nil
}

func (v *Validator) indicesOverlap(g *Game, cells []*Cell) error {
	for _, c := range cells {
		for _, bc := range g.Board.Cells {
			if c.Index == bc.Index {
				return errors.New("index overlap")
			}
		}
	}
	return nil
}

func (v *Validator) checkPlayerTiles(g *Game, cells []*Cell) error {
	for _, c := range cells {
		tile := g.ToMove().LookupTile(c.Tile.ID)
		if tile == nil {
			return fmt.Errorf("tile with ID %s not found", c.Tile.ID)
		}
	}
	return nil
}

func (v *Validator) firstPlaceAtOrigin(g *Game, cells []*Cell) error {
	if g.FirstMovePlayed() {
		return nil
	}
	for _, c := range cells {
		if c.Index == BoardOrigin {
			return nil
		}
	}
	return errors.New("first place must cross the origin")
}

func (v *Validator) checkStraightWord(g *Game, w *Word) error {
	if !(g.IsHorizontal(w) || g.IsVertical(w)) {
		return errors.New("word must be straight")
	}
	return nil
}

func (v *Validator) checkAdjacentTilesExist(g *Game, w *Word) error {
	if !g.FirstMovePlayed() {
		return nil
	}

	possibilities := []int{}
	for _, c := range w.Cells {
		possibilities = append(
			possibilities, c.Index-1, c.Index+1, c.Index+BoardHeight, c.Index-BoardHeight)
	}
	for _, c := range g.Board.Cells {
		if tools.ContainsInt(possibilities, c.Index) {
			return nil
		}
	}
	return errors.New("adjacent tiles not found")
}
