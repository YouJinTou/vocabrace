package scrabble

import (
	"errors"
	"fmt"
)

// Validator validates a word placement attempt.
type Validator struct{}

// NewValidator creates a new validator.
func NewValidator() *Validator {
	return &Validator{}
}

// ValidatePlace validates a place action.
func (v *Validator) ValidatePlace(g *Game, tiles []*Cell) error {
	if len(tiles) == 0 {
		return errors.New("at least one tile required to form a word")
	}

	if err := v.indicesWithinBounds(tiles); err != nil {
		return err
	}

	if err := v.indicesOverlap(g, tiles); err != nil {
		return err
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
