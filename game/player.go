package game

import (
	"GoTicTacToe/assets"
	"image/color"
)

type Player struct {
	Symbol *assets.Symbol
	Points int
	Color  color.Color
}

func NewPlayer(sym *assets.Symbol, color color.Color) *Player {
	return &Player{Symbol: sym, Points: 0, Color: color}
}
