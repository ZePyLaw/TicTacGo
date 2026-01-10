// Package main is the entry point of the GoTicTacToe application.
//
// File: main.go
//
// Project: GoTicTacToe
// Authors:
//   - Alexandre Schmid <alexandre.schmid@edu.heia-fr.ch>
//   - Jeremy Prin <jeremy.prin@edu.heia-fr.ch>
//
// Date: 09 January 2026
//
// Copyright:
//
//	Copyright (c) 2026 HEIA-FR / ISC
//	Haute école d'ingénierie et d'architecture de Fribourg
//	Informatique et Systèmes de Communication
//
// License:
//
//	SPDX-License-Identifier: MIT OR Apache-2.0
//
// Description:
//
//	This file contains the application entry point. It initializes the Ebiten
//	window, sets up the screen host, and launches the main game loop.
package main

import (
	"GoTicTacToe/screens"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

// Window configuration constants.
const (
	baseWindowWidth  = 1920
	baseWindowHeight = 1080

	// windowScaleFactor defines the initial window size as a ratio
	// of a Full HD resolution.
	windowScaleFactor = 0.7
)

// main initializes the window and starts the Ebiten game loop.
func main() {
	windowWidth := int(float64(baseWindowWidth) * windowScaleFactor)
	windowHeight := int(float64(baseWindowHeight) * windowScaleFactor)

	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Go Tic Tac Toe")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// Create the screen host and set the initial screen.
	host := screens.NewScreenHost()
	host.SetScreen(screens.NewStartScreen(host))

	// Run the game.
	if err := ebiten.RunGame(host); err != nil {
		log.Fatal(err)
	}
}
