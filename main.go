package main

import (
	"bytes"
	"fmt"
	"image/color"
	_ "image/png"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	windowWidth = 480 // default 480
	windowHeight = 600 // default 600
	boardWidth = 480 // default 480
	boardHeight = 480 // default 480
	fontSize    = 15
	bigFontSize = 100
	boardSize = 3 // default 3
	numSymbolToWin = 3 // default 3
)

type GameState int

const (
	INIT GameState = iota
	PLAYING
	GAME_END
)

type Player struct {
	Symbol string
	Points int
}

type Game struct {
	state         GameState
	board         *Board
	players       []*Player
	currentPlayer *Player
	winner        *Player
}

func (g *Game) Start() {
	g.state = PLAYING
	g.currentPlayer = g.players[0]
	g.winner = nil
}

func (g *Game) NextPlayer() {
	if len(g.players) == 0 {
		return
	}
	for i, p := range g.players {
		if p == g.currentPlayer {
			g.currentPlayer = g.players[(i+1)%len(g.players)]
			return
		}
	}
	g.currentPlayer = g.players[0]
}

func (g *Game) HandleDraw() {
	g.winner = nil
	g.state = GAME_END
}

func (g *Game) HandleWin(winner *Player) {
	winner.Points++
	g.winner = winner
	g.state = GAME_END
}

func (g *Game) ResetPoints() {
	for _, player := range g.players {
		player.Points = 0
	}
}

func (g *Game) Update() error {
	switch g.state {
	case INIT:
		g.Init()

	case PLAYING:
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			x, y := GetCursorBoardPos(g)

			if g.board.play(g.currentPlayer, x, y) == nil {
				if g.board.CheckDraw() {
					g.HandleDraw()
				}
				if winner := g.board.CheckWin(); winner != nil {
					g.HandleWin(winner)
				}
				g.NextPlayer()
			}
		}

	case GAME_END:
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			g.Init()
		}
	}

	if inpututil.KeyPressDuration(ebiten.KeyR) == 60 {
		g.Init()
		g.ResetPoints()
	}
	if inpututil.KeyPressDuration(ebiten.KeyEscape) == 60 {
		os.Exit(0)
	}
	return nil
}


var (
	normalText  text.Face
	bigText     text.Face
	gameImage   = ebiten.NewImage(windowWidth, windowHeight)
)

func GetCursorBoardPos(g *Game) (int, int) {
	mx, my := ebiten.CursorPosition()
	// Calculate the pixel size of each board cell
	cellSize := boardWidth / boardSize
	// Divide the cursor's pixel coordinates by the cell size to get the board coordinates
	return mx / cellSize, my / cellSize
}

func keyChangeColor(key ebiten.Key, screen *ebiten.Image) {
	if inpututil.KeyPressDuration(key) > 1 {
		var msgText string
		var colorText color.RGBA
		colorChange := 255 - (255 / 60 * uint8(inpututil.KeyPressDuration(key)))
		if key == ebiten.KeyEscape {
			msgText = fmt.Sprintf("CLOSING...")
			colorText = color.RGBA{R: 255, G: colorChange, B: colorChange, A: 255}
		} else if key == ebiten.KeyR {
			msgText = fmt.Sprintf("RESETING...")
			colorText = color.RGBA{R: colorChange, G: 255, B: 255, A: 255}
		}

		msgTextOptions := &text.DrawOptions{}
		msgTextOptions.GeoM.Translate(windowWidth/2, windowHeight-30)
		msgTextOptions.ColorScale.ScaleWithColor(colorText)
		text.Draw(screen, msgText, normalText, msgTextOptions)
	}
}


// DrawBoardLines draws the board lines on the screen.
// It draws the vertical and horizontal lines of the board.
func (g *Game) DrawBoardLines(screen *ebiten.Image) {
	cellSize := float64(boardWidth) / float64(g.board.size)
	lineThickness := 3.0

	// Create a gg context
	dc := gg.NewContext(boardWidth, boardHeight)
	dc.SetRGB(1, 1, 1) // White color
	dc.SetLineWidth(lineThickness)

	// Vertical lines
	for i := 1; i < g.board.size; i++ {
		x := float64(i) * cellSize
		dc.DrawLine(x, 0, x, float64(boardHeight))
		dc.Stroke()
	}

	// Horizontal lines
	for i := 1; i < g.board.size; i++ {
		y := float64(i) * cellSize
		dc.DrawLine(0, y, float64(boardWidth), y)
		dc.Stroke()
	}

	// Border
	dc.DrawRectangle(0, 0, float64(boardWidth), float64(boardHeight))
	dc.Stroke()

	// Convert gg image to ebiten image
	ebitenImg := ebiten.NewImageFromImage(dc.Image())
	screen.DrawImage(ebitenImg, nil)
}


func (g *Game) Draw(screen *ebiten.Image) {
	g.DrawBoardLines(screen)
	gameImage.Clear()

	// Board symbols drawing
	for y := 0; y < g.board.size; y++ {
		for x := 0; x < g.board.size; x++ {
			player := g.board.cells[x][y]
			if player != nil {
				g.DrawSymbol(x, y, player.Symbol)
			}
		}
	}

	screen.DrawImage(gameImage, nil)

	// Debug
	mx, my := ebiten.CursorPosition()
	msgFPS := fmt.Sprintf("TPS: %.2f | FPS: %.2f | Cursor: %v,%v", ebiten.ActualTPS(), ebiten.ActualFPS(), mx, my)
	op := &text.DrawOptions{}
	op.GeoM.Translate(0, windowHeight-60)
	op.ColorScale.ScaleWithColor(color.White)
	op.LayoutOptions.LineSpacing = 15
	text.Draw(screen, msgFPS, normalText, op)

	keyChangeColor(ebiten.KeyEscape, screen)
	keyChangeColor(ebiten.KeyR, screen)

	// Score
	msgOX := fmt.Sprintf("O: %v | X: %v", g.players[1].Points, g.players[0].Points)
	msgOXOptions := &text.DrawOptions{}
	msgOXOptions.GeoM.Translate(windowWidth/2, windowHeight-30)
	msgOXOptions.PrimaryAlign = text.AlignCenter
	msgOXOptions.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, msgOX, normalText, msgOXOptions)

	// Show Cursor
	msg := fmt.Sprintf("%v", g.currentPlayer.Symbol)
	msgOptions := &text.DrawOptions{}
	msgOptions.GeoM.Translate(float64(mx-15), float64(my-15))
	msgOptions.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, msg, normalText, msgOptions)

	// End message
	if g.state == GAME_END {
		var msgWin string
		if g.winner != nil {
			msgWin = fmt.Sprintf("%v wins!", g.winner.Symbol)
		} else {
			msgWin = "It's a draw!"
		}
		msgWinOptions := &text.DrawOptions{}
		msgWinOptions.GeoM.Translate(windowWidth/2, windowHeight/2)
		msgWinOptions.PrimaryAlign = text.AlignCenter
		msgWinOptions.SecondaryAlign = text.AlignCenter
		msgWinOptions.ColorScale.ScaleWithColor(color.White)
		text.Draw(screen, msgWin, bigText, msgWinOptions)
	}
}

func (g *Game) DrawSymbol(x, y int, sym string) {
	cellSize := float64(boardWidth) / float64(g.board.size)
	centerX := float64(x)*cellSize + cellSize/2
	centerY := float64(y)*cellSize + cellSize/2
	padding := cellSize * 0.2
	thickness := cellSize * 0.08

	// Create a gg context for the entire board
	dc := gg.NewContext(boardWidth, boardHeight)
	dc.SetRGB(1, 1, 1) // White color
	dc.SetLineWidth(thickness)

	switch sym {
	case "O":
		// Circle
		radius := (cellSize / 2) - padding
		dc.DrawCircle(centerX, centerY, radius)
		dc.Stroke()

	case "X":
		// Cross
		offset := (cellSize / 2) - padding
		dc.DrawLine(centerX-offset, centerY-offset, centerX+offset, centerY+offset)
		dc.Stroke()
		dc.DrawLine(centerX-offset, centerY+offset, centerX+offset, centerY-offset)
		dc.Stroke()
	}

	// Convert gg image to ebiten image and draw it
	symbolImg := ebiten.NewImageFromImage(dc.Image())
	op := &ebiten.DrawImageOptions{}
	gameImage.DrawImage(symbolImg, op)
}

func (g *Game) Init() {
	g.Load()
	g.ResetPoints()
	g.state = PLAYING
}

func (g *Game) Load() {
	gameImage.Clear()
	g.state = INIT
	g.board = NewBoard(boardSize, numSymbolToWin)
	g.players = []*Player{
		{Symbol: "X", Points: 0},
		{Symbol: "O", Points: 0},
	}
	g.currentPlayer = g.players[0]
	g.winner = nil
}


func init() {
	tt, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	normalText = &text.GoTextFace{
		Source: tt,
		Size:   fontSize,
	}

	bigText = &text.GoTextFace{
		Source: tt,
		Size:   bigFontSize,
	}
}

func newRandom() *rand.Rand {
	s1 := rand.NewSource(time.Now().UnixNano())
	return rand.New(s1)
}

func (g *Game) Layout(int, int) (int, int) {
	return windowWidth, windowHeight
}

func main() {
	game := &Game{}
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("TicTacToe")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
