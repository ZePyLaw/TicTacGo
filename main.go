package main

import (
	"bytes"
	"embed"
	"fmt"
	"image"
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
	windowWidth = 480
	windowHeight = 600
	boardWidth = 480
	boardHeight = 480
	fontSize    = 15
	bigFontSize = 100
	boardSize = 3
	numSymbolToWin = 3
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


//go:embed images/*
var imageFS embed.FS

var (
	normalText  text.Face
	bigText     text.Face
	boardImage  *ebiten.Image
	symbolImage *ebiten.Image
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
	screen.DrawImage(boardImage, nil)
	gameImage.Clear()

	// Board Drawing
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
	imageBytes, err := imageFS.ReadFile(fmt.Sprintf("images/%v.png", sym))
	if err != nil {
		log.Fatal(err)
	}
	decoded, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		log.Fatal(err)
	}
	symbolImage = ebiten.NewImageFromImage(decoded)
	opSymbol := &ebiten.DrawImageOptions{}
	opSymbol.GeoM.Translate(float64((160*(x+1)-160)+7), float64((160*(y+1)-160)+7))

	gameImage.DrawImage(symbolImage, opSymbol)
}

func (g *Game) Init() {
	imageBytes, err := imageFS.ReadFile("images/board.png")
	if err != nil {
		log.Fatal(err)
	}
	decoded, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		log.Fatal(err)
	}
	boardImage = ebiten.NewImageFromImage(decoded)
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
