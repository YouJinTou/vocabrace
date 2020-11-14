package scrabble

import (
	"errors"
	"fmt"
)

// Player encapsulates player data.
type Player struct {
	ID     string  `json:"id"`
	Name   string  `json:"n"`
	Points int     `json:"p"`
	Tiles  []*Tile `json:"t"`
}

// ExchangeTiles removes a set of tiles from the player's set of tiles and replaces them.
func (p *Player) ExchangeTiles(ids []string, toReceive []*Tile) ([]*Tile, error) {
	if len(ids) != len(toReceive) {
		return []*Tile{}, errors.New("exchange and receive tile counts should match")
	}

	returnTiles := []*Tile{}
	for _, tr := range ids {
		foundTile := false
		for i, t := range p.Tiles {
			if t.Index == tr {
				foundTile = true
				p.Tiles = append(p.Tiles[:i], p.Tiles[i+1:]...)
				returnTiles = append(returnTiles, t.Copy(true))
				break
			}
		}
		if !foundTile {
			return []*Tile{}, fmt.Errorf("%s tile not found", tr)
		}
	}

	for _, tr := range toReceive {
		p.Tiles = append(p.Tiles, tr)
	}

	return returnTiles, nil
}

// AwardPoints awards points to the player.
func (p *Player) AwardPoints(points int) {
	p.Points += points
}
