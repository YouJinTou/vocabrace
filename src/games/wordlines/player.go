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

// PlaceTiles removes tiles from the player's stack.
func (p *Player) PlaceTiles(IDs []string, toReceive *Tiles) error {
	for _, ID := range IDs {
		match := p.Tiles.RemoveByID(ID)
		if match == nil {
			return fmt.Errorf("%s tile not found", ID)
		}
	}

	for _, tr := range toReceive.Value {
		p.Tiles.Append(tr)
	}
	return nil
}

// ExchangeTiles removes a set of tiles from the player's set of tiles and replaces them.
// The first set of tiles are to be returned to the bag. The second set of tiles are the tiles
// the player will receive back.
func (p *Player) ExchangeTiles(ids []string, toReceive *Tiles) (*Tiles, *Tiles, error) {
	returnTiles := NewTiles()
	toReceiveTiles := NewTiles(toReceive.Value...)

	for _, tr := range ids {
		match := p.Tiles.RemoveByID(tr)
		if match == nil {
			return NewTiles(), NewTiles(), fmt.Errorf("%s tile not found", tr)
		}

		if returnTiles.Count() < toReceive.Count() {
			returnTiles.Append(match.Copy(true))
		} else {
			toReceiveTiles.Append(match.Copy(true))
		}
	}

	for _, tr := range toReceiveTiles.Value {
		p.Tiles.Append(tr)
	}

	return returnTiles, toReceiveTiles, nil
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
