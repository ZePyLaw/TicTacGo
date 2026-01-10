/**
 ******************************************************************************
 * @file            : fonts.go
 * @brief           : GoTicTacToe - Font assets management
 * @author          : Alexandre Schmid <alexandre.schmid@edu.heia-fr.ch>
 * @author          : Jeremy Prin <jeremy.prin@edu.heia-fr.ch>
 * @date            : 09. January 2026
 ******************************************************************************
 * @copyright   : Copyright (c) 2026 HEIA-FR / ISC
 *                Haute école d'ingénierie et d'architecture de Fribourg
 *                Informatique et Systèmes de Communication
 * @attention   : SPDX-License-Identifier: MIT OR Apache-2.0
 ******************************************************************************
 * @details
 * This file defines and initializes the font assets used throughout
 * the graphical user interface of the game.
 *
 * Fonts are loaded once at startup and shared globally via the assets package.
 ******************************************************************************
 */

package assets

import (
	"bytes"
	"log"

	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// Font size constants used across the user interface.
const (
	normalFontSize = 20
	bigFontSize    = 80
)

// NormalFont is the default font used for standard UI text
// such as labels, scores, and informational messages.
var NormalFont text.Face

// BigFont is a larger font used for prominent UI elements
// such as titles or end-of-game messages.
var BigFont text.Face

// init loads the font resources and initializes the shared font faces.
//
// The fonts are embedded via the Ebitengine example resources and
// converted into Go text faces usable by the rendering system.
func init() {
	tt, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}

	NormalFont = &text.GoTextFace{
		Source: tt,
		Size:   normalFontSize,
	}

	BigFont = &text.GoTextFace{
		Source: tt,
		Size:   bigFontSize,
	}
}
