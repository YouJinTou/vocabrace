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
func (p *Player) ExchangeTiles(toRemove []string, toReceive []*Tile) ([]*Tile, error) {
	if len(toRemove) != len(toReceive) {
		return []*Tile{}, errors.New("exchange and receive tile counts should match")
	}

	for _, tr := range toRemove {
		foundTile := false
		for _, t := range p.Tiles {
			if t.Letter == tr {
				foundTile = true
				break
			}
		}
		if !foundTile {
			return []*Tile{}, fmt.Errorf("%s tile not found", tr)
		}
	}

	returnTiles := []*Tile{}
	for i, tr := range toRemove {
		for _, t := range p.Tiles {
			if t.Letter == tr {
				returnTiles = append(returnTiles, t.Copy())
				p.Tiles = append(p.Tiles[:i], p.Tiles[i+1:]...)
				break
			}
		}
	}

	for _, tr := range toReceive {
		p.Tiles = append(p.Tiles, tr)
	}

	return returnTiles, nil
}
