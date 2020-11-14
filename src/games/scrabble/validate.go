package scrabble

import "errors"

// PlaceValidator validates a word placement attempt.
type PlaceValidator struct{}

// ValidatePlaceAction validates a place action.
func (v *PlaceValidator) ValidatePlaceAction(g *Game, tiles []*Cell) error {
	if len(tiles) == 0 {
		return errors.New("at least one tile required to form a word")
	}

	return nil
}
