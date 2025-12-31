package assets

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	Logo       *ebiten.Image
)

func init() {
	logoImg, _, err := ebitenutil.NewImageFromFile("assets/static/gonnectmax_logo.png")
	if err != nil {
		log.Fatal(err)
	}
	Logo = logoImg
}
