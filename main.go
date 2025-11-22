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
	boardLineThickness = 3.0
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
	symbol *Symbol
	points int
	color color.Color
}

type GameRessources struct {
	boardImage *ebiten.Image
}

type Game struct {
	state         GameState
	board         *Board
	players       []*Player
	currentPlayer *Player
	winner        *Player
	resources     *GameRessources
}

func (g *Game) Start() {
	g.state = PLAYING
	g.currentPlayer = g.players[0]
	g.winner = nil
	g.resources = &GameRessources{}
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
	winner.points++
	g.winner = winner
	g.state = GAME_END
}

func (g *Game) ResetPoints() {
	for _, player := range g.players {
		player.points = 0
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

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw board
	screen.DrawImage(g.resources.boardImage, nil)
	gameImage.Clear()

	// Board symbols drawing
	for y := 0; y < g.board.size; y++ {
		for x := 0; x < g.board.size; x++ {
			player := g.board.cells[x][y]
			if player != nil {
				g.DrawSymbol(x, y, player)
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
	msgOX := fmt.Sprintf("O: %v | X: %v", g.players[1].points, g.players[0].points)
	msgOXOptions := &text.DrawOptions{}
	msgOXOptions.GeoM.Translate(windowWidth/2, windowHeight-30)
	msgOXOptions.PrimaryAlign = text.AlignCenter
	msgOXOptions.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, msgOX, normalText, msgOXOptions)

	// Show Cursor
	msg := fmt.Sprintf("%v", g.currentPlayer)
	msgOptions := &text.DrawOptions{}
	msgOptions.GeoM.Translate(float64(mx-15), float64(my-15))
	msgOptions.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, msg, normalText, msgOptions)

	// End message
	if g.state == GAME_END {
		var msgWin string
		if g.winner != nil {
			msgWin = fmt.Sprintf("%v wins!", g.winner)
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

func (g *Game) DrawSymbol(gridX, gridY int, p *Player) {
    cellSize := float64(boardWidth) / float64(g.board.size)

    // Padding (15% de la taille de la cellule)
    padding := cellSize * 0.1
    usableSize := cellSize - 2*padding

    // Facteur d'échelle pour que l'image tienne dans la zone utile
    scale := usableSize / float64(imageSize)

    // Position centrée et padding intégré
    px := float64(gridX)*cellSize + padding
    py := float64(gridY)*cellSize + padding

    op := &ebiten.DrawImageOptions{}
    op.Filter = ebiten.FilterLinear
    op.GeoM.Scale(scale, scale)
    op.GeoM.Translate(px, py)

    gameImage.DrawImage(p.symbol.image, op)
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
		{symbol: newSymbol(CROSS), points: 0},
		{symbol: newSymbol(CIRCLE), points: 0},
	}
	g.resources = &GameRessources{}
	g.resources.boardImage = g.board.GenerateImage()
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
