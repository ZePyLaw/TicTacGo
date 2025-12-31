package main

import (
	"GoTicTacToe/screens"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	windowWidth  = 1920 * 0.7
	windowHeight = 1080 * 0.7
)

func main() {
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Go Tic Tac Toe")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	host := screens.NewScreenHost()
	host.SetScreen(screens.NewStartScreen(host)) // start screen

	if err := ebiten.RunGame(host); err != nil {
		log.Fatal(err)
	}
}
