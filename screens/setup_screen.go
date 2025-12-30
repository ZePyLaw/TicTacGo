package screens

import (
	"GoTicTacToe/ai_models"
	"GoTicTacToe/assets"
	"GoTicTacToe/ui"
	uiutils "GoTicTacToe/ui/utils"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type playerCardButtons struct {
	role       *ui.Button
	symbolPrev *ui.Button
	symbolNext *ui.Button
	ready      *ui.Button
}

// SetupScreen lets the user configure players and board size before starting.
type SetupScreen struct {
	host          ScreenHost
	config        GameConfig
	buttons       []*ui.Button
	playerButtons []playerCardButtons
	addPlayerBtn  *ui.Button
	root          *ui.Container
	background    *ebiten.Image
	cardBg        *ebiten.Image
	accentStrip   *ebiten.Image
	symbolImages  map[assets.SymbolType]*ebiten.Image
}

const (
	cardWidth    = 280.0
	cardHeight   = 200.0
	cardSpacingX = 30.0
	cardSpacingY = 26.0
	cardStartY   = -10.0
	maxPlayers   = 6
	cardsPerRow  = 4
)

var (
	playerPalette = []color.RGBA{
		{R: 255, G: 99, B: 132, A: 255},
		{R: 54, G: 162, B: 235, A: 255},
		{R: 75, G: 192, B: 192, A: 255},
		{R: 255, G: 206, B: 86, A: 255},
		{R: 153, G: 102, B: 255, A: 255},
		{R: 201, G: 203, B: 207, A: 255},
	}

	playerSymbolOrder = []assets.SymbolType{
		assets.CircleSymbol,
		assets.CrossSymbol,
		assets.TriangleSymbol,
		assets.SquareSymbol,
	}
)

func NewSetupScreen(h ScreenHost, baseCfg GameConfig) *SetupScreen {
	cfg := baseCfg
	if cfg.BoardSize == 0 || len(cfg.Players) == 0 {
		cfg = DefaultGameConfig()
	}
	if cfg.ToWin == 0 {
		cfg.ToWin = cfg.BoardSize
	}

	for i := range cfg.Players {
		if cfg.Players[i].IsAI && !cfg.Players[i].Ready {
			cfg.Players[i].Ready = true
		}
	}

	s := &SetupScreen{
		host:   h,
		config: cfg,
	}
	s.buildButtons()
	s.refreshLabels()
	return s
}

func (s *SetupScreen) buildButtons() {
	s.buttons = nil
	s.addPlayerBtn = nil
	s.playerButtons = make([]playerCardButtons, len(s.config.Players))

	// Root container stretches to the full screen so children can anchor against it.
	s.root = ui.NewContainer(0, 0, 100, 100, uiutils.AnchorTopLeft, 0, uiutils.TransparentWidgetStyle)
	s.root.WidthMode = uiutils.SizeFill
	s.root.HeightMode = uiutils.SizeFill

	// Board size controls
	minus := ui.NewButton("-", -160, -70, uiutils.AnchorTopCenter,
		70, 46, buttonRadius, uiutils.DefaultWidgetStyle,
		func() { s.changeBoardSize(-1) })
	plus := ui.NewButton("+", 160, -70, uiutils.AnchorTopCenter,
		70, 46, buttonRadius, uiutils.DefaultWidgetStyle,
		func() { s.changeBoardSize(+1) })

	s.buttons = append(s.buttons, minus, plus)

	// Player panels (2 columns)
	for i := range s.config.Players {
		cx, cy := s.cardCenter(i)

		roleBtn := ui.NewButton("", cx+cardWidth/2-70, cy-cardHeight/2+26, uiutils.AnchorCenter,
			120, 36, buttonRadius, uiutils.DefaultWidgetStyle,
			func(idx int) func() {
				return func() { s.cycleRole(idx) }
			}(i),
		)

		// AI players are auto-ready; humans start not ready.
		if s.config.Players[i].IsAI && !s.config.Players[i].Ready {
			s.config.Players[i].Ready = true
		}

		symbolPrev := ui.NewButton("<", cx-90, cy+6, uiutils.AnchorCenter,
			60, 50, buttonRadius, uiutils.TransparentWidgetStyle,
			func(idx int) func() {
				return func() { s.cycleSymbol(idx, -1) }
			}(i),
		)

		symbolNext := ui.NewButton(">", cx+90, cy+6, uiutils.AnchorCenter,
			60, 50, buttonRadius, uiutils.TransparentWidgetStyle,
			func(idx int) func() {
				return func() { s.cycleSymbol(idx, +1) }
			}(i),
		)

		readyBtn := ui.NewButton("", cx, cy+cardHeight/2-28, uiutils.AnchorCenter,
			cardWidth-32, 44, buttonRadius, uiutils.DefaultWidgetStyle,
			func(idx int) func() {
				return func() { s.toggleReady(idx) }
			}(i),
		)

		s.playerButtons[i] = playerCardButtons{
			role:       roleBtn,
			symbolPrev: symbolPrev,
			symbolNext: symbolNext,
			ready:      readyBtn,
		}

		s.buttons = append(s.buttons, roleBtn, symbolPrev, symbolNext, readyBtn)
	}

	if len(s.config.Players) < maxPlayers {
		cx, cy := s.cardCenter(len(s.config.Players))
		s.addPlayerBtn = ui.NewButton("+ Add Player", cx, cy, uiutils.AnchorCenter,
			cardWidth-32, cardHeight-32, buttonRadius, uiutils.TransparentWidgetStyle,
			func() { s.addPlayer() },
		)
		s.buttons = append(s.buttons, s.addPlayerBtn)
	}

	startBtn := ui.NewButton("Start Game", -140, 320, uiutils.AnchorCenter,
		220, 60, buttonRadius, uiutils.NormalWidgetStyle,
		func() { s.startGame() })
	backBtn := ui.NewButton("Back", 140, 320, uiutils.AnchorCenter,
		220, 60, buttonRadius, uiutils.TransparentWidgetStyle,
		func() { s.host.SetScreen(NewStartScreen(s.host)) })

	s.buttons = append(s.buttons, startBtn, backBtn)

	for _, b := range s.buttons {
		s.root.AddChild(b)
	}
}

func (s *SetupScreen) Update() error {
	if s.root != nil {
		s.root.Update()
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		s.handleSymbolClick()
	}
	return nil
}

func (s *SetupScreen) Draw(screen *ebiten.Image) {
	w, h := screen.Bounds().Dx(), screen.Bounds().Dy()
	if s.background == nil {
		topColor := color.RGBA{R: 0x11, G: 0x1f, B: 0x39, A: 0xFF}
		bottomColor := color.RGBA{R: 0x0a, G: 0x12, B: 0x24, A: 0xFF}
		s.background = uiutils.CreateGradientBackground(w, h, topColor, bottomColor)
	}
	screen.DrawImage(s.background, nil)

	titleOpts := &text.DrawOptions{}
	titleOpts.PrimaryAlign = text.AlignCenter
	titleOpts.SecondaryAlign = text.AlignCenter
	titleOpts.ColorScale.ScaleWithColor(color.White)
	titleOpts.GeoM.Translate(float64(w)/2, 60)
	text.Draw(screen, "Custom Game Setup", assets.BigFont, titleOpts)

	boardMsg := fmt.Sprintf("Grid: %dx%d (win in %d)", s.config.BoardSize, s.config.BoardSize, s.config.ToWin)
	boardOpts := &text.DrawOptions{}
	boardOpts.PrimaryAlign = text.AlignCenter
	boardOpts.SecondaryAlign = text.AlignCenter
	boardOpts.ColorScale.ScaleWithColor(color.RGBA{R: 200, G: 220, B: 255, A: 255})
	boardOpts.GeoM.Translate(float64(w)/2, 120)
	text.Draw(screen, boardMsg, assets.NormalFont, boardOpts)

	rootRect := s.root.LayoutRect()
	centerX := rootRect.X + rootRect.Width/2
	centerY := rootRect.Y + rootRect.Height/2

	s.drawPlayerCards(screen, centerX, centerY)

	// Buttons
	if s.root != nil {
		s.root.Draw(screen)
	}
}

func (s *SetupScreen) ensureCardAssets() {
	if s.cardBg == nil {
		bgColor := color.RGBA{R: 24, G: 34, B: 58, A: 230}
		s.cardBg = uiutils.CreateRoundedRect(int(cardWidth), int(cardHeight), 14, bgColor)
	}
	if s.accentStrip == nil {
		s.accentStrip = ebiten.NewImage(int(cardWidth), 14)
	}
	if s.symbolImages == nil {
		s.symbolImages = make(map[assets.SymbolType]*ebiten.Image)
	}
}

func (s *SetupScreen) symbolImage(sym assets.SymbolType) *ebiten.Image {
	if img, ok := s.symbolImages[sym]; ok && img != nil {
		return img
	}
	img := assets.NewSymbol(sym).Image
	s.symbolImages[sym] = img
	return img
}

func (s *SetupScreen) drawPlayerCards(screen *ebiten.Image, centerX, centerY float64) {
	s.ensureCardAssets()

	for i, pc := range s.config.Players {
		s.drawPlayerCard(screen, i, pc, centerX, centerY)
	}

	if s.addPlayerBtn != nil {
		s.drawAddCard(screen, len(s.config.Players), centerX, centerY)
	}
}

func (s *SetupScreen) drawPlayerCard(screen *ebiten.Image, idx int, pc PlayerConfig, centerX, centerY float64) {
	cx, cy := s.cardCenter(idx)
	cx += centerX
	cy += centerY

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(cx-cardWidth/2, cy-cardHeight/2)
	screen.DrawImage(s.cardBg, op)

	playerColor := pc.Color
	if playerColor == nil {
		playerColor = playerPalette[idx%len(playerPalette)]
	}

	baseColor := color.RGBAModel.Convert(playerColor).(color.RGBA)
	accentColor := baseColor
	accentColor.A = 210
	s.accentStrip.Fill(accentColor)

	accentOp := &ebiten.DrawImageOptions{}
	accentOp.GeoM.Translate(cx-cardWidth/2, cy-cardHeight/2)
	screen.DrawImage(s.accentStrip, accentOp)

	name := fmt.Sprintf("Player %d", idx+1)
	if pc.Name != "" {
		name = pc.Name
	}

	nameOpts := &text.DrawOptions{}
	nameOpts.PrimaryAlign = text.AlignStart
	nameOpts.SecondaryAlign = text.AlignCenter
	nameOpts.ColorScale.ScaleWithColor(color.White)
	nameOpts.GeoM.Translate(cx-cardWidth/2+14, cy-cardHeight/2+28)
	text.Draw(screen, name, assets.NormalFont, nameOpts)

	stateLabel := s.roleLabel(pc)
	stateOpts := &text.DrawOptions{}
	stateOpts.PrimaryAlign = text.AlignStart
	stateOpts.SecondaryAlign = text.AlignCenter
	stateOpts.ColorScale.ScaleWithColor(color.RGBA{R: 180, G: 200, B: 230, A: 255})
	stateOpts.GeoM.Translate(cx-cardWidth/2+14, cy-cardHeight/2+58)
	text.Draw(screen, stateLabel, assets.NormalFont, stateOpts)

	symbolImg := s.symbolImage(pc.Symbol)
	if symbolImg != nil {
		iconSize := cardHeight * 0.35
		srcW := float64(symbolImg.Bounds().Dx())
		srcH := float64(symbolImg.Bounds().Dy())
		scale := iconSize / srcH
		if srcW > srcH {
			scale = iconSize / srcW
		}

		symOp := &ebiten.DrawImageOptions{}
		symOp.Filter = ebiten.FilterLinear
		symOp.GeoM.Scale(scale, scale)
		symOp.GeoM.Translate(cx-iconSize/2, cy-cardHeight/2+80)
		symOp.ColorScale.ScaleWithColor(baseColor)
		screen.DrawImage(symbolImg, symOp)
	}

}

func (s *SetupScreen) drawAddCard(screen *ebiten.Image, idx int, centerX, centerY float64) {
	cx, cy := s.cardCenter(idx)
	cx += centerX
	cy += centerY

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(cx-cardWidth/2, cy-cardHeight/2)
	screen.DrawImage(s.cardBg, op)

	labelOpts := &text.DrawOptions{}
	labelOpts.PrimaryAlign = text.AlignCenter
	labelOpts.SecondaryAlign = text.AlignCenter
	labelOpts.ColorScale.ScaleWithColor(color.RGBA{R: 210, G: 230, B: 255, A: 255})
	labelOpts.GeoM.Translate(cx, cy)
	text.Draw(screen, "Add Player", assets.NormalFont, labelOpts)
}

func (s *SetupScreen) changeBoardSize(delta int) {
	newSize := s.config.BoardSize + delta
	if newSize < 3 {
		newSize = 3
	}
	if newSize > 8 {
		newSize = 8
	}
	s.config.BoardSize = newSize
	s.config.ToWin = newSize
}

func (s *SetupScreen) cardCenter(idx int) (float64, float64) {
	col := idx % cardsPerRow
	row := idx / cardsPerRow

	offsetX := (float64(col) * (cardWidth + cardSpacingX)) - ((float64(cardsPerRow-1))*(cardWidth+cardSpacingX))/2
	offsetY := cardStartY + float64(row)*(cardHeight+cardSpacingY)
	return offsetX, offsetY
}

func (s *SetupScreen) handleSymbolClick() {
	if s.root == nil {
		return
	}
	mx, my := ebiten.CursorPosition()

	rootRect := s.root.LayoutRect()
	centerX := rootRect.X + rootRect.Width/2
	centerY := rootRect.Y + rootRect.Height/2

	iconSize := cardHeight * 0.35
	for i := range s.config.Players {
		cx, cy := s.cardCenter(i)
		cx += centerX
		cy += centerY

		iconX := cx - iconSize/2
		iconY := cy - cardHeight/2 + 80

		if float64(mx) >= iconX && float64(mx) <= iconX+iconSize &&
			float64(my) >= iconY && float64(my) <= iconY+iconSize {
			s.cycleColor(i)
			break
		}
	}
}

func (s *SetupScreen) cycleRole(idx int) {
	pc := &s.config.Players[idx]

	state := "human"
	if pc.IsAI {
		switch pc.AIModel.(type) {
		case ai_models.MinimaxAI:
			state = "ai-hard"
		default:
			state = "ai-easy"
		}
	}

	switch state {
	case "human":
		pc.IsAI = true
		pc.Ready = true // AI auto ready
		if pc.AIModel == nil {
			pc.AIModel = ai_models.RandomAI{}
		} else {
			switch pc.AIModel.(type) {
			case ai_models.MinimaxAI:
				pc.AIModel = ai_models.RandomAI{}
			}
		}
	case "ai-easy":
		pc.IsAI = true
		pc.Ready = true
		pc.AIModel = ai_models.MinimaxAI{}
	case "ai-hard":
		s.removePlayer(idx)
		return
	default: // removed
		pc.IsAI = false
		pc.Ready = false
	}
	s.refreshLabels()
}

func (s *SetupScreen) removePlayer(idx int) {
	if idx < 0 || idx >= len(s.config.Players) {
		return
	}
	if len(s.config.Players) <= 1 {
		return
	}

	// Remove player at idx
	s.config.Players = append(s.config.Players[:idx], s.config.Players[idx+1:]...)
	s.buildButtons()
	s.refreshLabels()
}

func (s *SetupScreen) toggleReady(idx int) {
	pc := &s.config.Players[idx]
	pc.Ready = !pc.Ready
	if pc.IsAI && pc.AIModel == nil {
		pc.AIModel = ai_models.RandomAI{}
	}
	s.refreshLabels()
}

func (s *SetupScreen) cycleSymbol(idx int, delta int) {
	pc := &s.config.Players[idx]

	current := 0
	for i, sym := range playerSymbolOrder {
		if pc.Symbol == sym {
			current = i
			break
		}
	}

	next := (current + delta) % len(playerSymbolOrder)
	if next < 0 {
		next += len(playerSymbolOrder)
	}
	pc.Symbol = playerSymbolOrder[next]
	s.refreshLabels()
}

func (s *SetupScreen) cycleColor(idx int) {
	pc := &s.config.Players[idx]

	if pc.Color == nil {
		pc.Color = playerPalette[0]
	}

	current := 0
	for i, c := range playerPalette {
		if colorsEqual(pc.Color, c) {
			current = i
			break
		}
	}
	pc.Color = playerPalette[(current+1)%len(playerPalette)]
	s.refreshLabels()
}

func (s *SetupScreen) addPlayer() {
	if len(s.config.Players) >= maxPlayers {
		return
	}

	idx := len(s.config.Players)
	s.config.Players = append(s.config.Players, PlayerConfig{
		Name:   fmt.Sprintf("Player %d", idx+1),
		Color:  playerPalette[idx%len(playerPalette)],
		Symbol: playerSymbolOrder[idx%len(playerSymbolOrder)],
	})

	s.buildButtons()
	s.refreshLabels()
}

func (s *SetupScreen) refreshLabels() {
	for i, pc := range s.config.Players {
		if i >= len(s.playerButtons) {
			continue
		}
		pb := s.playerButtons[i]

		if pb.role != nil {
			pb.role.Label = s.roleLabel(pc)
		}
		if pb.symbolPrev != nil {
			pb.symbolPrev.Label = "<"
		}
		if pb.symbolNext != nil {
			pb.symbolNext.Label = ">"
		}
		if pb.ready != nil {
			if pc.Ready {
				pb.ready.Label = "Ready"
				setReadyButtonStyle(pb.ready, true)
			} else {
				pb.ready.Label = "Not Ready"
				setReadyButtonStyle(pb.ready, false)
			}
		}
	}

	if s.addPlayerBtn != nil {
		s.addPlayerBtn.Label = fmt.Sprintf("+ Add Player (%d/%d)", len(s.config.Players), maxPlayers)
	}
}

func setReadyButtonStyle(btn *ui.Button, ready bool) {
	if btn == nil {
		return
	}

	style := uiutils.DefaultWidgetStyle
	if ready {
		style = uiutils.SuccessWidgetStyle
	}
	btn.Style = style
	//btn.Hoverable = uiutils.NewHoverable(style.HoverMode)
}

func (s *SetupScreen) roleLabel(pc PlayerConfig) string {
	if pc.IsAI {
		switch pc.AIModel.(type) {
		case ai_models.MinimaxAI:
			return "AI (Hard)"
		default:
			return "AI (Easy)"
		}
	}

	return "Human"
}

func (s *SetupScreen) startGame() {
	// Ensure at least two player slots exist
	for len(s.config.Players) < 2 && len(s.config.Players) < maxPlayers {
		idx := len(s.config.Players)
		s.config.Players = append(s.config.Players, PlayerConfig{
			Name:   fmt.Sprintf("Player %d", idx+1),
			Color:  playerPalette[idx%len(playerPalette)],
			Symbol: playerSymbolOrder[idx%len(playerSymbolOrder)],
			Ready:  false,
		})
	}

	// Guarantee at least two ready players before starting.
	readyCount := 0
	for i := range s.config.Players {
		if s.config.Players[i].Ready {
			readyCount++
		}
	}
	if readyCount < 2 {
		for i := range s.config.Players {
			if !s.config.Players[i].Ready {
				s.config.Players[i].Ready = true
				readyCount++
				if readyCount >= 2 {
					break
				}
			}
		}
	}

	s.refreshLabels()
	s.host.SetScreen(NewGameScreen(s.host, s.config))
}

func symbolName(sym assets.SymbolType) string {
	switch sym {
	case assets.CircleSymbol:
		return "Circle"
	case assets.CrossSymbol:
		return "Cross"
	case assets.TriangleSymbol:
		return "Triangle"
	case assets.SquareSymbol:
		return "Square"
	default:
		return "Circle"
	}
}

func colorsEqual(a color.Color, b color.Color) bool {
	if a == nil || b == nil {
		return a == b
	}
	ar, ag, ab, aa := a.RGBA()
	br, bg, bb, ba := b.RGBA()
	return ar == br && ag == bg && ab == bb && aa == ba
}

func hexColor(c color.Color) string {
	if c == nil {
		c = playerPalette[0]
	}
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("%02X%02X%02X", uint8(r>>8), uint8(g>>8), uint8(b>>8))
}
