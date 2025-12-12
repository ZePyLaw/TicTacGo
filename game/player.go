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

func (p *Player) Opponent(players [2]*Player) *Player {
	if p == players[0] {
		return players[1]
	}
	return players[0]
}
