package scrabble

import (
	"errors"
	"fmt"
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

	if err := v.indicesOverlap(&g, w.Cells); err != nil {
		return err
	}

	if err := v.checkPlayerTiles(&g, w.Cells); err != nil {
		return err
	}

	if !v.wc.IsValidWord(g.Language, w.String()) {
		return fmt.Errorf("invalid word %s", w.String())
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
