package screens

import (
	"GoTicTacToe/ui"
	"fmt"
	"image"

	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/v2"
)

// Screen is an interface that all screens must implement
// It defines the methods that a screen must implement to be used in the game.
type Screen interface {
	Update() error
	Draw(screen *ebiten.Image)
}

type ScreenHost interface {
	SetScreen(Screen)
}

type screenHost struct {
	width, height int
	current       Screen
	debugui       debugui.DebugUI
}

func NewScreenHost(width, height int) *screenHost {
	ui.UpdateScreenSize(width, height)
	return &screenHost{
		width:  width,
		height: height,
	}
}

func (h *screenHost) SetScreen(s Screen) {
	h.current = s
}

func (h *screenHost) Update() error {
	if _, err := h.debugui.Update(func(ctx *debugui.Context) error {
		ctx.Window("Debug", image.Rect(0, 0, 100, 100), func(layout debugui.ContainerLayout) {
			msgTps := fmt.Sprintf("TPS: %.2f", ebiten.ActualTPS())
			msgFps := fmt.Sprintf("FPS: %.2f", ebiten.ActualFPS())
			ctx.Text(msgTps)
			ctx.Text(msgFps)
		})
		return nil
	}); err != nil {
		return err
	}
	if h.current != nil {
		return h.current.Update()
	}
	return nil
}

func (h *screenHost) Draw(screen *ebiten.Image) {
	w, hgt := screen.Bounds().Dx(), screen.Bounds().Dy()
	if w != 0 && hgt != 0 {
		ui.UpdateScreenSize(w, hgt)
	}
	h.debugui.Draw(screen)
	if h.current != nil {
		h.current.Draw(screen)
	}
}

func (h *screenHost) Layout(outsideWidth, outsideHeight int) (int, int) {
	if outsideWidth > 0 && outsideHeight > 0 {
		h.width = outsideWidth
		h.height = outsideHeight
	}
	ui.UpdateScreenSize(h.width, h.height)
	return h.width, h.height
}
