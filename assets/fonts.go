package assets

import (
	"bytes"
	"log"

	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	NormalFont text.Face
	BigFont    text.Face
)

func init() {
	tt, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}

	NormalFont = &text.GoTextFace{
		Source: tt,
		Size:   20,
	}

	BigFont = &text.GoTextFace{
		Source: tt,
		Size:   80,
	}
}
