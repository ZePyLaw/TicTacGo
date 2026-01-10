// Package screens provides screen-related utilities and helpers used by the UI.
//
// File: window.go
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
//	This file contains small utility helpers related to window and screen
//	management, shared across different screens.
package screens

import "github.com/hajimehoshi/ebiten/v2"

// GetWindowSize returns the current window size in pixels.
//
// The width and height are returned as float64 values to simplify
// layout and scaling computations in the UI layer.
func GetWindowSize() (float64, float64) {
	w, h := ebiten.WindowSize()
	return float64(w), float64(h)
}
