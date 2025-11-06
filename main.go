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
	windowWidth      = 480 // default 480
	windowHeight     = 600 // default 600
	boardWidth       = 480 // default 480
	boardHeight      = 480 // default 480
	fontSize         = 15
	bigFontSize      = 100
	boardSize        = 3 // default 3
	numSymbolToWin   = 3 // default 3
	lineThickness    = 3.0
	keyDurationExit  = 60
	keyDurationReset = 60
	cursorOffset     = 15.0
	messageYOffset   = 30.0
	fpsYOffset       = 60.0
	colorMaxValue    = 255
	keyColorFadeDiv  = 60
	paddingRatio     = 0.2
	thicknessRatio   = 0.08
	aiMoveDelay      = 500 * time.Millisecond // delay before AI move
)

// --- Game logic types and constants ---

type GameState int

const (
	INIT GameState = iota
	PLAYING
	GAME_END
)

type Player struct {
	Symbol string
	Points int
	IsAI   bool // define if this player is an AI
}

type Game struct {
	state         GameState
	board         *Board
	players       []*Player
	currentPlayer *Player
	winner        *Player
	aiEnabled     bool
	lastAIMove    time.Time // NEW: track AI move timing
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

// Helper: safely get cursor position
func safeCursorPosition() (int, int, error) {
	mx, my := ebiten.CursorPosition()
	if mx < 0 || my < 0 {
		return 0, 0, fmt.Errorf("invalid cursor position: %d,%d", mx, my)
	}
	return mx, my, nil
}

// --- MINIMAX ALGORITHM SECTION ---

// evaluateBoard analyse l'état actuel du plateau et retourne un score numérique.
// +1 si l'IA gagne, -1 si le joueur humain gagne, 0 sinon.
func (g *Game) evaluateBoard() int {
	// On demande au plateau s'il y a un vainqueur.
	if winner := g.board.CheckWin(); winner != nil {
		// Si le gagnant est le joueur IA (O)
		if winner.Symbol == g.players[1].Symbol {
			return +1 // Victoire de l'IA, score positif
		} else {
			return -1 // Victoire du joueur humain, score négatif
		}
	}
	// Aucun gagnant, score neutre
	return 0
}

// isMovesLeft vérifie s'il reste encore des cases libres à jouer sur le plateau.
// Renvoie true s'il reste au moins une case vide, false sinon.
func (g *Game) isMovesLeft() bool {
	for y := 0; y < g.board.size; y++ { // Parcourt chaque ligne
		for x := 0; x < g.board.size; x++ { // Parcourt chaque colonne
			if g.board.cells[x][y] == nil { // Si une case est vide
				return true // Il reste donc des coups possibles
			}
		}
	}
	return false // Aucun mouvement disponible, le plateau est plein
}

// minimax explore récursivement tous les coups possibles pour déterminer
// le meilleur score que l'IA peut obtenir contre un joueur parfait.
func (g *Game) minimax(isMaximizing bool) int {
	// 1. On évalue immédiatement la position actuelle
	// pour savoir si quelqu'un a déjà gagné ou perdu.
	score := g.evaluateBoard()

	// Si un score non nul est trouvé (victoire ou défaite),
	// on renvoie ce score et on arrête la recherche ici.
	if score != 0 {
		return score
	}

	// 2. Si plus aucun mouvement n'est possible (plateau plein),
	// on considère que c'est un match nul et on renvoie 0.
	if !g.isMovesLeft() {
		return 0
	}

	// 3. Sinon, on continue l'exploration récursive :
	// - Si c'est le tour de l'IA (isMaximizing = true),
	//   on cherche à maximiser le score.
	// - Si c'est le tour du joueur (isMaximizing = false),
	//   on cherche à minimiser le score.

	if isMaximizing {
		// Cas 1 : Tour de l'IA, on cherche le meilleur score possible
		best := -9999 // On initialise avec une valeur très basse

		// Parcourt toutes les cases du plateau
		for y := 0; y < g.board.size; y++ {
			for x := 0; x < g.board.size; x++ {
				// Si la case est vide, l'IA peut y jouer virtuellement
				if g.board.cells[x][y] == nil {
					// L'IA joue virtuellement sur cette case
					g.board.cells[x][y] = g.players[1]

					// On appelle minimax en passant false,
					// car c'est maintenant au joueur humain de jouer
					val := g.minimax(false)

					// On annule le coup (backtracking)
					// pour ne pas interférer avec les autres simulations
					g.board.cells[x][y] = nil

					// Si le score obtenu est meilleur que le meilleur trouvé,
					// on met à jour la variable best
					if val > best {
						best = val
					}
				}
			}
		}
		// On renvoie le meilleur score trouvé
		return best

	} else {
		// Cas 2 : Tour du joueur humain, on cherche le pire score pour l'IA
		best := 9999 // On initialise avec une valeur très haute

		// Parcourt toutes les cases du plateau
		for y := 0; y < g.board.size; y++ {
			for x := 0; x < g.board.size; x++ {
				// Si la case est vide, le joueur peut y jouer virtuellement
				if g.board.cells[x][y] == nil {
					// Le joueur joue virtuellement sur cette case
					g.board.cells[x][y] = g.players[0]

					// On appelle minimax en passant true,
					// car c'est maintenant au tour de l'IA
					val := g.minimax(true)

					// On annule le coup (backtracking)
					g.board.cells[x][y] = nil

					// Si le score obtenu est plus petit,
					// on met à jour la variable best
					if val < best {
						best = val
					}
				}
			}
		}
		// On renvoie le score le plus petit (le moins favorable à l'IA)
		return best
	}
}

// chooseBestMove détermine le meilleur coup possible pour l'IA
// en utilisant la fonction minimax pour évaluer chaque position possible.
func (g *Game) chooseBestMove() (int, int) {
	// On initialise une valeur très basse comme "pire" score possible.
	// Cela permet de comparer plus facilement les futurs résultats.
	bestVal := -9999

	// On initialise les coordonnées du meilleur coup à des valeurs invalides.
	// Si aucun coup n'est trouvé, cela évitera de jouer sur une case incorrecte.
	bestMove := [2]int{-1, -1}

	// On parcourt toutes les cases du plateau ligne par ligne.
	for y := 0; y < g.board.size; y++ {
		for x := 0; x < g.board.size; x++ {

			// Si la case est vide, elle peut être testée comme un coup potentiel.
			if g.board.cells[x][y] == nil {

				// L'IA joue virtuellement sur cette case.
				// On place temporairement son symbole.
				g.board.cells[x][y] = g.players[1]

				// On appelle la fonction minimax pour évaluer la conséquence
				// de ce coup, en supposant que le joueur jouera ensuite.
				moveVal := g.minimax(false)

				// Une fois la simulation terminée, on annule le coup.
				// Cela permet de garder le plateau propre pour les autres tests.
				g.board.cells[x][y] = nil

				// Si la valeur obtenue est meilleure que la meilleure connue,
				// cela signifie que ce coup est plus favorable pour l'IA.
				if moveVal > bestVal {
					// On mémorise ce coup comme le meilleur trouvé.
					bestMove = [2]int{x, y}
					// Et on met à jour la meilleure valeur.
					bestVal = moveVal
				}
			}
		}
	}

	// Après avoir testé toutes les cases, on renvoie les coordonnées
	// du meilleur coup trouvé par l'IA.
	return bestMove[0], bestMove[1]
}

// --- END MINIMAX ---

// AI move (with delay)
func (g *Game) playAI() {
	if g.currentPlayer == nil || !g.currentPlayer.IsAI {
		return
	}

	// Wait some time before AI plays
	if time.Since(g.lastAIMove) < aiMoveDelay {
		return
	}

	x, y := g.chooseBestMove()
	if x < 0 || y < 0 {
		return
	}

	_ = g.board.play(g.currentPlayer, x, y)
	if g.board.CheckDraw() {
		g.HandleDraw()
	}
	if winner := g.board.CheckWin(); winner != nil {
		g.HandleWin(winner)
	}
	g.NextPlayer()
	g.lastAIMove = time.Now()
}

func (g *Game) Update() error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in Update: %v", r)
		}
	}()

	switch g.state {
	case INIT:
		g.Init()

	case PLAYING:
		if g.currentPlayer.IsAI {
			g.playAI()
			return nil
		}

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			x, y := GetCursorBoardPos(g)
			if err := g.board.play(g.currentPlayer, x, y); err != nil {
				log.Printf("invalid move: %v", err)
				return nil
			}
			if g.board.CheckDraw() {
				g.HandleDraw()
			}
			if winner := g.board.CheckWin(); winner != nil {
				g.HandleWin(winner)
			}
			g.NextPlayer()
			g.lastAIMove = time.Now() // start AI timer after human move
		}

	case GAME_END:
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			g.Init()
		}
	}

	if inpututil.KeyPressDuration(ebiten.KeyR) == keyDurationReset {
		g.Init()
		g.ResetPoints()
	}
	if inpututil.KeyPressDuration(ebiten.KeyEscape) == keyDurationExit {
		os.Exit(0)
	}
	return nil
}

var (
	normalText text.Face
	bigText    text.Face
	gameImage  = ebiten.NewImage(windowWidth, windowHeight)
)

func GetCursorBoardPos(g *Game) (int, int) {
	mx, my, err := safeCursorPosition()
	if err != nil {
		log.Println(err)
		return 0, 0
	}
	cellSize := boardWidth / boardSize
	return mx / cellSize, my / cellSize
}

func drawFadingText(screen *ebiten.Image, msg string, clr color.RGBA) {
	msgTextOptions := &text.DrawOptions{}
	msgTextOptions.GeoM.Translate(windowWidth/2, windowHeight-messageYOffset)
	msgTextOptions.ColorScale.ScaleWithColor(clr)
	text.Draw(screen, msg, normalText, msgTextOptions)
}

func keyChangeColor(key ebiten.Key, screen *ebiten.Image) {
	if inpututil.KeyPressDuration(key) <= 1 {
		return
	}
	colorChange := colorMaxValue - (colorMaxValue / keyColorFadeDiv * uint8(inpututil.KeyPressDuration(key)))
	var msg string
	var clr color.RGBA

	switch key {
	case ebiten.KeyEscape:
		msg = "CLOSING..."
		clr = color.RGBA{R: colorMaxValue, G: colorChange, B: colorChange, A: colorMaxValue}
	case ebiten.KeyR:
		msg = "RESETING..."
		clr = color.RGBA{R: colorChange, G: colorMaxValue, B: colorMaxValue, A: colorMaxValue}
	default:
		return
	}

	drawFadingText(screen, msg, clr)
}

// DrawBoardLines draws the board lines on the screen.
func (g *Game) DrawBoardLines(screen *ebiten.Image) {
	cellSize := float64(boardWidth) / float64(g.board.size)
	dc := gg.NewContext(boardWidth, boardHeight)
	dc.SetRGB(1, 1, 1)
	dc.SetLineWidth(lineThickness)
	for i := 1; i < g.board.size; i++ {
		pos := float64(i) * cellSize
		dc.DrawLine(pos, 0, pos, float64(boardHeight))
		dc.DrawLine(0, pos, float64(boardWidth), pos)
		dc.Stroke()
	}
	dc.DrawRectangle(0, 0, float64(boardWidth), float64(boardHeight))
	dc.Stroke()
	ebitenImg := ebiten.NewImageFromImage(dc.Image())
	screen.DrawImage(ebitenImg, nil)
}

func (g *Game) Draw(screen *ebiten.Image) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in Draw: %v", r)
		}
	}()

	g.DrawBoardLines(screen)
	gameImage.Clear()

	for y := 0; y < g.board.size; y++ {
		for x := 0; x < g.board.size; x++ {
			player := g.board.cells[x][y]
			if player != nil {
				g.DrawSymbol(x, y, player.Symbol)
			}
		}
	}

	screen.DrawImage(gameImage, nil)

	mx, my, _ := safeCursorPosition()
	msgFPS := fmt.Sprintf("TPS: %.2f | FPS: %.2f | Cursor: %v,%v", ebiten.ActualTPS(), ebiten.ActualFPS(), mx, my)
	op := &text.DrawOptions{}
	op.GeoM.Translate(0, windowHeight-fpsYOffset)
	op.ColorScale.ScaleWithColor(color.White)
	op.LayoutOptions.LineSpacing = fontSize
	text.Draw(screen, msgFPS, normalText, op)

	keyChangeColor(ebiten.KeyEscape, screen)
	keyChangeColor(ebiten.KeyR, screen)

	msgOX := fmt.Sprintf("O: %v | X: %v", g.players[1].Points, g.players[0].Points)
	msgOXOptions := &text.DrawOptions{}
	msgOXOptions.GeoM.Translate(windowWidth/2, windowHeight-messageYOffset)
	msgOXOptions.PrimaryAlign = text.AlignCenter
	msgOXOptions.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, msgOX, normalText, msgOXOptions)

	msg := fmt.Sprintf("%v", g.currentPlayer.Symbol)
	msgOptions := &text.DrawOptions{}
	msgOptions.GeoM.Translate(float64(mx)-cursorOffset, float64(my)-cursorOffset)
	msgOptions.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, msg, normalText, msgOptions)

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
	padding := cellSize * paddingRatio
	thickness := cellSize * thicknessRatio
	dc := gg.NewContext(boardWidth, boardHeight)
	dc.SetRGB(1, 1, 1)
	dc.SetLineWidth(thickness)
	switch sym {
	case "O":
		radius := (cellSize / 2) - padding
		dc.DrawCircle(centerX, centerY, radius)
	case "X":
		offset := (cellSize / 2) - padding
		dc.DrawLine(centerX-offset, centerY-offset, centerX+offset, centerY+offset)
		dc.DrawLine(centerX-offset, centerY+offset, centerX+offset, centerY-offset)
	}
	dc.Stroke()
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
		{Symbol: "X", Points: 0, IsAI: false},
		{Symbol: "O", Points: 0, IsAI: true},
	}
	g.currentPlayer = g.players[0]
	g.winner = nil
	g.aiEnabled = true
	g.lastAIMove = time.Now()
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
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("Critical error: %v", r)
		}
	}()
	game := &Game{}
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("TicTacToe")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatalf("runtime error: %v", err)
	}
}
