package wordlines

import (
	"fmt"
)

// Player encapsulates player data.
type Player struct {
	ID     string `json:"i"`
	Name   string `json:"n"`
	Points int    `json:"p"`
	Tiles  *Tiles `json:"t"`
}

// ExchangeTiles removes a set of tiles from the player's set of tiles and replaces them.
func (p *Player) ExchangeTiles(ids []string, toReceive *Tiles) (*Tiles, error) {
	returnTiles := NewTiles()

	for _, tr := range ids {
		match := p.Tiles.RemoveByID(tr)
		if match == nil {
			return NewTiles(), fmt.Errorf("%s tile not found", tr)
		}

		if returnTiles.Count() < toReceive.Count() {
			returnTiles.Append(match.Copy(true))
		}
	}

	for _, tr := range toReceive.Value {
		p.Tiles.Append(tr)
	}

	return returnTiles, nil
}

// HasTiles checks if the player has any tiles.
func (p *Player) HasTiles() bool {
	return p.Tiles.Count() > 0
}

// AwardPoints awards points to the player.
func (p *Player) AwardPoints(points int) {
	p.Points += points
}

// LookupTile finds a tile in the player's stack given an ID.
func (p *Player) LookupTile(ID string) *Tile {
	return p.Tiles.FindByID(ID)
}
