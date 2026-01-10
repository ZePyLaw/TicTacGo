package screens

import (
	"GoTicTacToe/ui"
	"fmt"
	"image"

	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/v2"
)

// Screen is implemented by every application screen (menu, game, setup, etc.).
//
// The screenHost will forward Ebiten's Update and Draw calls to the currently
// active Screen.
type Screen interface {
	Update() error
	Draw(screen *ebiten.Image)
}

// ScreenHost allows screens to request a screen change (navigation).
type ScreenHost interface {
	SetScreen(Screen)
}

// Debug window layout constants.
const (
	debugWindowX      = 0
	debugWindowY      = 0
	debugWindowWidth  = 100
	debugWindowHeight = 100
)

// screenHost owns the currently active screen and the debug UI instance.
//
// It is also used as the Ebiten game object (implements Update/Draw/Layout).
type screenHost struct {
	current Screen
	debugui debugui.DebugUI
}

// NewScreenHost creates a new screen host with no active screen.
// Call SetScreen to display the first screen.
func NewScreenHost() *screenHost {
	return &screenHost{}
}

// SetScreen changes the currently active screen.
func (h *screenHost) SetScreen(s Screen) {
	h.current = s
}

// Update updates the debug UI and forwards the update call to the current screen.
func (h *screenHost) Update() error {
	if _, err := h.debugui.Update(func(ctx *debugui.Context) error {
		ctx.Window("Debug", image.Rect(debugWindowX, debugWindowY, debugWindowWidth, debugWindowHeight), func(layout debugui.ContainerLayout) {
			// Display real-time Ebiten performance metrics.
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

// Draw draws the debug overlay and forwards the draw call to the current screen.
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

// Layout handles window resizing and updates global UI sizing.
//
// Ebiten calls Layout to determine the logical screen size used for rendering.
func (h *screenHost) Layout(outsideWidth, outsideHeight int) (int, int) {
	ui.UpdateScreenSize(outsideWidth, outsideHeight)
	return outsideWidth, outsideHeight
}
