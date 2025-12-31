package game

import (
	"GoTicTacToe/assets"
	"image/color"
)

type Player struct {
	Symbol *assets.Symbol
	Points int
	Color  color.Color
	Name   string
	IsAI   bool
}


func NewPlayer(sym *assets.Symbol, color color.Color) *Player {
	return &Player{Symbol: sym, Points: 0, Color: color}
}


// Opponent returns all the players except the current player
func (p *Player) Opponent(players []*Player) *Player {
	for _, candidate := range players {
		if candidate != p {
			return candidate
		}
	}
	return nil
}
