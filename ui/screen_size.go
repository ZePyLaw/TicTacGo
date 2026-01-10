// Package ui contains reusable UI widgets and views rendered with Ebiten.
//
// File: screen_size.go
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
//	This file stores and exposes the current drawable screen size.
//	It allows UI widgets to position and scale themselves consistently
//	without directly querying Ebiten at draw time.
package ui

var (
	screenWidth  int
	screenHeight int
)

// UpdateScreenSize stores the current drawable screen size.
//
// This function should be called whenever the window size changes
// (typically from Ebiten's Layout method).
func UpdateScreenSize(w, h int) {
	screenWidth = w
	screenHeight = h
}

// currentScreenSize returns the last known drawable screen size.
//
// This function is internal to the ui package and is used by layout helpers
// to compute widget positioning.
func currentScreenSize() (int, int) {
	return screenWidth, screenHeight
}
