package scrabble

import "errors"

// Validator validates a word placement attempt.
type Validator struct{}

// NewValidator creates a new validator.
func NewValidator() *Validator {
	return Validator{}
}

// ValidatePlace validates a place action.
func (v *Validator) ValidatePlace(g *Game, tiles []*Cell) error {
	if len(tiles) == 0 {
		return errors.New("at least one tile required to form a word")
	}

	return nil
}
