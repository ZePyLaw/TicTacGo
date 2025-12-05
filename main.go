package main

import (
	"GoTicTacToe/screens"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	windowWidth  = 1920 / 2
	windowHeight = 1080 / 2
)

func main() {
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Go Tic Tac Toe")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	host := screens.NewScreenHost(windowWidth, windowHeight)
	host.SetScreen(screens.NewStartScreen(host)) // start screen

	if err := ebiten.RunGame(host); err != nil {
		log.Fatal(err)
	}
}
