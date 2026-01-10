package assets

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Path to the application logo image.
const logoImagePath = "assets/static/gonnectmax_logo.png"

// Logo is the application logo displayed in the user interface.
var Logo *ebiten.Image

// init loads the static image assets.
//
// The logo image is loaded at startup and stored in a shared variable
// to avoid reloading it multiple times during the game.
func init() {
	logoImg, _, err := ebitenutil.NewImageFromFile(logoImagePath)
	if err != nil {
		log.Fatal(err)
	}
	Logo = logoImg
}
