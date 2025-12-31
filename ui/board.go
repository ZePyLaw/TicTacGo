package ui

import (
	"GoTicTacToe/game"
	"GoTicTacToe/ui/utils"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// BoardView is the visual component responsible for rendering the
// tic-tac-toe board and handling user interaction.
type BoardView struct {
	Widget // Embeds Widget: inherits size, position, anchor, AbsPosition(), etc.

	logicBoard  *game.Board      // Reference to the logical board
	OnCellClick func(cx, cy int) // Callback triggered when a cell is clicked

	lastGridW int
	lastGridH int
}

// NewBoardView creates a new visual board component.
func NewBoardView(
	board *game.Board,
	x, y, size float64,
	style utils.WidgetStyle,
	onClick func(cx, cy int),
) *BoardView {

	view := &BoardView{
		Widget: Widget{
			OffsetX: x,
			OffsetY: y,
			Width:   size,
			Height:  size,
			Anchor:  utils.AnchorCenter,
			Style:   style,
		},
		logicBoard:  board,
		OnCellClick: onClick,
	}

	// Pre-generate the grid once (static background)
	view.ensureGridImage(size, size)

	return view
}

// createGridImage renders the static background grid (lines + background)
// and returns the resulting image.
func (v *BoardView) createGridImage(width, height int) *ebiten.Image {
	img := ebiten.NewImage(width, height)

	// Fill board background
	img.Fill(v.Style.BackgroundNormal)

	cellWidth := float64(width) / float64(v.logicBoard.Width)
	cellHeight := float64(height) / float64(v.logicBoard.Height)
	thickness := v.Style.BorderWidth
	lineColor := v.Style.BorderColor

	// Draw vertical grid lines
	for i := 1; i < v.logicBoard.Width; i++ {
		offset := float64(i) * cellWidth

		vert := ebiten.NewImage(int(thickness), height)
		vert.Fill(lineColor)

		opv := &ebiten.DrawImageOptions{}
		opv.GeoM.Translate(offset-thickness/2, 0)
		img.DrawImage(vert, opv)
	}

	// Draw horizontal grid lines
	for i := 1; i < v.logicBoard.Height; i++ {
		offset := float64(i) * cellHeight

		hori := ebiten.NewImage(width, int(thickness))
		hori.Fill(lineColor)

		oph := &ebiten.DrawImageOptions{}
		oph.GeoM.Translate(0, offset-thickness/2)
		img.DrawImage(hori, oph)
	}

	return img
}

func (v *BoardView) ensureGridImage(width, height float64) {
	w := int(width)
	h := int(height)

	if v.Widget.image == nil || v.lastGridW != w || v.lastGridH != h {
		v.Widget.image = v.createGridImage(w, h)
		v.lastGridW = w
		v.lastGridH = h
	}
}

// Update handles mouse click detection and cell coordinate translation.
func (v *BoardView) Update() {
	rect := v.LayoutRect()
	v.ensureGridImage(rect.Width, rect.Height)

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		vx, vy := rect.X, rect.Y

		// Check if click is inside the board boundaries
		if float64(mx) >= vx && float64(mx) <= vx+rect.Width &&
			float64(my) >= vy && float64(my) <= vy+rect.Height {

			cellWidth := rect.Width / float64(v.logicBoard.Width)
			cellHeight := rect.Height / float64(v.logicBoard.Height)

			// Convert pixel coordinates → board grid coordinates
			gridX := int((float64(mx) - vx) / cellWidth)
			gridY := int((float64(my) - vy) / cellHeight)

			// Trigger callback
			if v.OnCellClick != nil {
				v.OnCellClick(gridX, gridY)
			}
		}
	}
}

// Draw renders the grid and the X/O symbols for each cell.
func (v *BoardView) Draw(screen *ebiten.Image) {
	rect := v.LayoutRect()
	v.ensureGridImage(rect.Width, rect.Height)
	vx, vy := rect.X, rect.Y

	// --- Draw static grid background ---
	opGrid := &ebiten.DrawImageOptions{}
	srcW := float64(v.Widget.image.Bounds().Dx())
	srcH := float64(v.Widget.image.Bounds().Dy())
	if srcW != 0 && srcH != 0 {
		opGrid.GeoM.Scale(rect.Width/srcW, rect.Height/srcH)
	}
	opGrid.GeoM.Translate(vx, vy)
	screen.DrawImage(v.Widget.image, opGrid)

	cellWidth := rect.Width / float64(v.logicBoard.Width)
	cellHeight := rect.Height / float64(v.logicBoard.Height)

	// Use the smaller dimension for symbol sizing to maintain aspect ratio
	cellSize := cellWidth
	if cellHeight < cellSize {
		cellSize = cellHeight
	}
	padding := cellSize * 0.1          // 10% padding inside each cell
	usableSize := cellSize - 2*padding // Space available for the symbol

	// --- Draw all symbols (X and O) ---
	for x := 0; x < v.logicBoard.Width; x++ {
		for y := 0; y < v.logicBoard.Height; y++ {
			p := v.logicBoard.Cells[x][y]
			if p == nil || p.Symbol.Image == nil {
				continue
			}

			symbolImg := p.Symbol.Image
			srcW, srcH := symbolImg.Bounds().Dx(), symbolImg.Bounds().Dy()

			// Determine the scaling factor based on largest dimension
			maxDim := float64(srcW)
			if srcH > srcW {
				maxDim = float64(srcH)
			}
			scale := usableSize / maxDim

			// Draw the scaled symbol inside the cell
			opSym := &ebiten.DrawImageOptions{}
			opSym.Filter = ebiten.FilterLinear // Smooth scaling

			// Scale first
			opSym.GeoM.Scale(scale, scale)

			// Position inside the cell with padding, centered
			symbolW := float64(srcW) * scale
			symbolH := float64(srcH) * scale
			drawX := vx + float64(x)*cellWidth + (cellWidth-symbolW)/2
			drawY := vy + float64(y)*cellHeight + (cellHeight-symbolH)/2

			opSym.GeoM.Translate(drawX, drawY)

			// Change the color of the image for the player color
			opSym.ColorScale.ScaleWithColor(p.Color)

			screen.DrawImage(symbolImg, opSym)
		}
	}
}
