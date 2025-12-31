package screens

import (
	"GoTicTacToe/ai_models"
	"GoTicTacToe/assets"
	"GoTicTacToe/ui"
	uiutils "GoTicTacToe/ui/utils"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// playerCardButtons groups all the interactive buttons associated with a single player card.
type playerCardButtons struct {
	role       *ui.Button // Button to cycle through player roles (Human/AI Easy/AI Hard)
	symbolPrev *ui.Button // Button to select the previous symbol
	symbolNext *ui.Button // Button to select the next symbol
	ready      *ui.Button // Button to toggle the player's ready state
}

// SetupScreen lets the user configure players and board size before starting.
// It provides controls for grid dimensions, win condition, and player configuration.
type SetupScreen struct {
	host          ScreenHost           // Reference to the screen manager for navigation
	config        GameConfig           // Current game configuration being edited
	buttons       []*ui.Button         // All interactive buttons on the screen
	playerCards   []*ui.PlayerCardView // Visual cards displaying player info
	playerButtons []playerCardButtons  // Button groups for each player
	addPlayerBtn  *ui.Button           // Button to add a new player
	startBtn      *ui.Button           // Button to start the game
	root          *ui.Container        // Root UI container for layout
	background    *ebiten.Image        // Cached gradient background
}

// Layout constants for player cards.
const (
	cardWidth    = 280.0 // Width of each player card in pixels
	cardHeight   = 200.0 // Height of each player card in pixels
	cardSpacingX = 30.0  // Horizontal spacing between cards
	cardSpacingY = 26.0  // Vertical spacing between card rows
	cardStartY   = 30.0  // Vertical offset from center for first row of cards
	maxPlayers   = 4     // Maximum number of players allowed
	cardsPerRow  = 4     // Number of player cards per row
)

// Grid size constraints.
const (
	minGridSize = 3 // Minimum grid dimension
	maxGridSize = 8 // Maximum grid dimension
	minToWin    = 3 // Minimum symbols needed to win
)

// playerPalette defines the available colors for players.
var playerPalette = []color.RGBA{
	{R: 255, G: 99, B: 132, A: 255},  // Pink/Red
	{R: 54, G: 162, B: 235, A: 255},  // Blue
	{R: 75, G: 192, B: 192, A: 255},  // Teal
	{R: 255, G: 206, B: 86, A: 255},  // Yellow
	{R: 153, G: 102, B: 255, A: 255}, // Purple
	{R: 201, G: 203, B: 207, A: 255}, // Gray
}

// playerSymbolOrder defines the order in which symbols can be cycled.
var playerSymbolOrder = []assets.SymbolType{
	assets.CircleSymbol,
	assets.CrossSymbol,
	assets.TriangleSymbol,
	assets.SquareSymbol,
}

// NewSetupScreen creates a new setup screen with the given base configuration.
// If the base configuration is empty or invalid, defaults are applied.
func NewSetupScreen(h ScreenHost, baseCfg GameConfig) *SetupScreen {
	cfg := baseCfg

	// Apply defaults if configuration is empty
	if cfg.BoardWidth == 0 || cfg.BoardHeight == 0 || len(cfg.Players) == 0 {
		cfg = DefaultGameConfig()
	}

	// Ensure ToWin is valid
	if cfg.ToWin == 0 {
		cfg.ToWin = minToWin
	}
	cfg.ToWin = clampToWin(cfg.ToWin, cfg.BoardWidth, cfg.BoardHeight)

	// Auto-ready AI players
	for i := range cfg.Players {
		if cfg.Players[i].IsAI && !cfg.Players[i].Ready {
			cfg.Players[i].Ready = true
		}
	}

	s := &SetupScreen{
		host:   h,
		config: cfg,
	}
	s.init()
	s.refreshLabels()
	return s
}

// clampToWin ensures ToWin is within valid bounds based on grid dimensions.
// ToWin must be between minToWin and the smallest grid dimension.
func clampToWin(toWin, width, height int) int {
	maxToWin := width
	if height < maxToWin {
		maxToWin = height
	}
	if toWin < minToWin {
		return minToWin
	}
	if toWin > maxToWin {
		return maxToWin
	}
	return toWin
}

// init builds all UI elements for the setup screen.
func (s *SetupScreen) init() {
	s.buttons = nil
	s.addPlayerBtn = nil
	s.playerCards = make([]*ui.PlayerCardView, len(s.config.Players))
	s.playerButtons = make([]playerCardButtons, len(s.config.Players))

	// Root container stretches to the full screen so children can anchor against it.
	s.root = ui.NewContainer(0, 0, 100, 100, uiutils.AnchorTopLeft, 0, uiutils.TransparentWidgetStyle)
	s.root.WidthMode = uiutils.SizeFill
	s.root.HeightMode = uiutils.SizeFill

	// Grid configuration controls (positioned below title)
	s.buildGridControls()

	// Player cards and their associated buttons
	s.buildPlayerCards()

	// Add player button (shown if below max players)
	if len(s.config.Players) < maxPlayers {
		cx, cy := s.cardCenter(len(s.config.Players))
		s.addPlayerBtn = ui.NewButton("+ Add Player", cx, cy, uiutils.AnchorCenter,
			cardWidth-32, cardHeight-32, buttonRadius, uiutils.TransparentWidgetStyle,
			func() { s.addPlayer() },
		)
		s.buttons = append(s.buttons, s.addPlayerBtn)
	}

	// Bottom action buttons
	s.startBtn = ui.NewButton("Start Game", -140, 320, uiutils.AnchorCenter,
		220, 60, buttonRadius, uiutils.NormalWidgetStyle,
		func() { s.startGame() })
	backBtn := ui.NewButton("Back", 140, 320, uiutils.AnchorCenter,
		220, 60, buttonRadius, uiutils.TransparentWidgetStyle,
		func() { s.host.SetScreen(NewStartScreen(s.host)) })

	s.buttons = append(s.buttons, s.startBtn, backBtn)

	// Add all buttons to the root container
	for _, b := range s.buttons {
		s.root.AddChild(b)
	}
}

// buildGridControls creates the buttons for adjusting grid width, height, and win condition.
func (s *SetupScreen) buildGridControls() {
	controlY := -230.0 // Y position relative to center

	// Width controls: [-] Width: X [+]
	s.buttons = append(s.buttons,
		ui.NewButton("-", -280, controlY, uiutils.AnchorCenter,
			50, 40, buttonRadius, uiutils.DefaultWidgetStyle,
			func() { s.changeGridWidth(-1) }),
		ui.NewButton("+", -130, controlY, uiutils.AnchorCenter,
			50, 40, buttonRadius, uiutils.DefaultWidgetStyle,
			func() { s.changeGridWidth(+1) }),
	)

	// Height controls: [-] Height: X [+]
	s.buttons = append(s.buttons,
		ui.NewButton("-", -30, controlY, uiutils.AnchorCenter,
			50, 40, buttonRadius, uiutils.DefaultWidgetStyle,
			func() { s.changeGridHeight(-1) }),
		ui.NewButton("+", 120, controlY, uiutils.AnchorCenter,
			50, 40, buttonRadius, uiutils.DefaultWidgetStyle,
			func() { s.changeGridHeight(+1) }),
	)

	// ToWin controls: [-] Win: X [+]
	s.buttons = append(s.buttons,
		ui.NewButton("-", 220, controlY, uiutils.AnchorCenter,
			50, 40, buttonRadius, uiutils.DefaultWidgetStyle,
			func() { s.changeToWin(-1) }),
		ui.NewButton("+", 370, controlY, uiutils.AnchorCenter,
			50, 40, buttonRadius, uiutils.DefaultWidgetStyle,
			func() { s.changeToWin(+1) }),
	)
}

// buildPlayerCards creates the player cards and their associated control buttons.
func (s *SetupScreen) buildPlayerCards() {
	for i := range s.config.Players {
		cx, cy := s.cardCenter(i)

		// Create the player card widget
		card := ui.NewPlayerCard(cx, cy, cardWidth, cardHeight, uiutils.AnchorCenter)
		card.OnSymbolClick = func(idx int) func() {
			return func() { s.cycleColor(idx) }
		}(i)
		s.playerCards[i] = card
		s.root.AddChild(card)

		// Role button (top-right of card)
		roleBtn := ui.NewButton("", cx+cardWidth/2-70, cy-cardHeight/2+26, uiutils.AnchorCenter,
			120, 36, buttonRadius, uiutils.DefaultWidgetStyle,
			func(idx int) func() {
				return func() { s.cycleRole(idx) }
			}(i),
		)

		// Auto-ready AI players
		if s.config.Players[i].IsAI && !s.config.Players[i].Ready {
			s.config.Players[i].Ready = true
		}

		// Symbol navigation buttons (left and right of symbol)
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

		// Ready button (bottom of card)
		readyBtn := ui.NewButton("", cx, cy+cardHeight/2-28, uiutils.AnchorCenter,
			cardWidth-32, 44, buttonRadius, uiutils.DefaultWidgetStyle,
			func(idx int) func() {
				return func() { s.toggleReady(idx) }
			}(i),
		)

		// Link ready button to card for automatic style updates
		card.ReadyButton = readyBtn

		s.playerButtons[i] = playerCardButtons{
			role:       roleBtn,
			symbolPrev: symbolPrev,
			symbolNext: symbolNext,
			ready:      readyBtn,
		}

		s.buttons = append(s.buttons, roleBtn, symbolPrev, symbolNext, readyBtn)
	}
}

// Update processes input and updates UI state each frame.
func (s *SetupScreen) Update() error {
	if s.root != nil {
		s.root.Update()
	}
	return nil
}

// Draw renders the setup screen to the provided image.
func (s *SetupScreen) Draw(screen *ebiten.Image) {
	w, h := screen.Bounds().Dx(), screen.Bounds().Dy()

	// Create and cache gradient background
	if s.background == nil {
		topColor := color.RGBA{R: 0x11, G: 0x1f, B: 0x39, A: 0xFF}
		bottomColor := color.RGBA{R: 0x0a, G: 0x12, B: 0x24, A: 0xFF}
		s.background = uiutils.CreateGradientBackground(w, h, topColor, bottomColor)
	}
	screen.DrawImage(s.background, nil)

	// Draw title
	titleOpts := &text.DrawOptions{}
	titleOpts.PrimaryAlign = text.AlignCenter
	titleOpts.SecondaryAlign = text.AlignCenter
	titleOpts.ColorScale.ScaleWithColor(color.White)
	titleOpts.GeoM.Translate(float64(w)/2, 60)
	text.Draw(screen, "Custom Game Setup", assets.BigFont, titleOpts)

	// Draw grid configuration info
	s.drawGridInfo(screen, w)

	// Draw all UI elements
	if s.root != nil {
		s.root.Draw(screen)
	}
}

// drawGridInfo renders the grid dimension and win condition labels.
func (s *SetupScreen) drawGridInfo(screen *ebiten.Image, screenWidth int) {
	screenHeight := screen.Bounds().Dy()
	centerX := float64(screenWidth) / 2
	centerY := float64(screenHeight) / 2

	// Position labels at the same Y as the control buttons
	infoY := centerY - 230
	textColor := color.RGBA{R: 200, G: 220, B: 255, A: 255}

	// Width label (centered between the - and + buttons at -280 and -130)
	widthOpts := &text.DrawOptions{}
	widthOpts.PrimaryAlign = text.AlignCenter
	widthOpts.SecondaryAlign = text.AlignCenter
	widthOpts.ColorScale.ScaleWithColor(textColor)
	widthOpts.GeoM.Translate(centerX-205, infoY)
	text.Draw(screen, fmt.Sprintf("Width: %d", s.config.BoardWidth), assets.NormalFont, widthOpts)

	// Height label (centered between the - and + buttons at -30 and 120)
	heightOpts := &text.DrawOptions{}
	heightOpts.PrimaryAlign = text.AlignCenter
	heightOpts.SecondaryAlign = text.AlignCenter
	heightOpts.ColorScale.ScaleWithColor(textColor)
	heightOpts.GeoM.Translate(centerX+45, infoY)
	text.Draw(screen, fmt.Sprintf("Height: %d", s.config.BoardHeight), assets.NormalFont, heightOpts)

	// Win condition label (centered between the - and + buttons at 220 and 370)
	winOpts := &text.DrawOptions{}
	winOpts.PrimaryAlign = text.AlignCenter
	winOpts.SecondaryAlign = text.AlignCenter
	winOpts.ColorScale.ScaleWithColor(textColor)
	winOpts.GeoM.Translate(centerX+295, infoY)
	text.Draw(screen, fmt.Sprintf("Win: %d", s.config.ToWin), assets.NormalFont, winOpts)
}

// changeGridWidth adjusts the grid width by delta, clamping to valid bounds.
func (s *SetupScreen) changeGridWidth(delta int) {
	newWidth := s.config.BoardWidth + delta
	if newWidth < minGridSize {
		newWidth = minGridSize
	}
	if newWidth > maxGridSize {
		newWidth = maxGridSize
	}
	s.config.BoardWidth = newWidth

	// Adjust ToWin if it exceeds the new minimum dimension
	s.config.ToWin = clampToWin(s.config.ToWin, s.config.BoardWidth, s.config.BoardHeight)
}

// changeGridHeight adjusts the grid height by delta, clamping to valid bounds.
func (s *SetupScreen) changeGridHeight(delta int) {
	newHeight := s.config.BoardHeight + delta
	if newHeight < minGridSize {
		newHeight = minGridSize
	}
	if newHeight > maxGridSize {
		newHeight = maxGridSize
	}
	s.config.BoardHeight = newHeight

	// Adjust ToWin if it exceeds the new minimum dimension
	s.config.ToWin = clampToWin(s.config.ToWin, s.config.BoardWidth, s.config.BoardHeight)
}

// changeToWin adjusts the win condition by delta, clamping to valid bounds.
func (s *SetupScreen) changeToWin(delta int) {
	s.config.ToWin = clampToWin(s.config.ToWin+delta, s.config.BoardWidth, s.config.BoardHeight)
}

// cardCenter calculates the center position for a player card at the given index.
func (s *SetupScreen) cardCenter(idx int) (float64, float64) {
	col := idx % cardsPerRow
	row := idx / cardsPerRow

	// Center cards horizontally based on total width
	offsetX := (float64(col) * (cardWidth + cardSpacingX)) - ((float64(cardsPerRow-1))*(cardWidth+cardSpacingX))/2
	offsetY := cardStartY + float64(row)*(cardHeight+cardSpacingY)
	return offsetX, offsetY
}

// cycleRole cycles through player roles: Human -> AI Easy -> AI Hard -> Remove (or back to Human if last player).
func (s *SetupScreen) cycleRole(idx int) {
	pc := &s.config.Players[idx]

	// Determine current role state
	state := "human"
	if pc.IsAI {
		switch pc.AIModel.(type) {
		case ai_models.MinimaxAI:
			state = "ai-hard"
		default:
			state = "ai-easy"
		}
	}

	// Cycle to next role
	switch state {
	case "human":
		pc.IsAI = true
		pc.Ready = true
		pc.AIModel = ai_models.RandomAI{}
	case "ai-easy":
		pc.IsAI = true
		pc.Ready = true
		pc.AIModel = ai_models.MinimaxAI{}
	case "ai-hard":
		if len(s.config.Players) <= 1 {
			// Can't remove the last player, cycle back to human
			pc.IsAI = false
			pc.Ready = false
			pc.AIModel = nil
		} else {
			s.removePlayer(idx)
			return
		}
	}
	s.refreshLabels()
}

// removePlayer removes the player at the given index from the configuration.
func (s *SetupScreen) removePlayer(idx int) {
	if idx < 0 || idx >= len(s.config.Players) {
		return
	}
	if len(s.config.Players) <= 1 {
		return
	}

	s.config.Players = append(s.config.Players[:idx], s.config.Players[idx+1:]...)
	s.init()
	s.refreshLabels()
}

// toggleReady toggles the ready state of the player at the given index.
func (s *SetupScreen) toggleReady(idx int) {
	pc := &s.config.Players[idx]
	pc.Ready = !pc.Ready

	// Ensure AI players have a model assigned
	if pc.IsAI && pc.AIModel == nil {
		pc.AIModel = ai_models.RandomAI{}
	}
	s.refreshLabels()
}

// cycleSymbol cycles the player's symbol forward or backward in the symbol order.
func (s *SetupScreen) cycleSymbol(idx int, delta int) {
	pc := &s.config.Players[idx]

	// Find current symbol index
	current := 0
	for i, sym := range playerSymbolOrder {
		if pc.Symbol == sym {
			current = i
			break
		}
	}

	// Calculate next index with wrapping
	next := (current + delta) % len(playerSymbolOrder)
	if next < 0 {
		next += len(playerSymbolOrder)
	}
	pc.Symbol = playerSymbolOrder[next]
	s.refreshLabels()
}

// cycleColor cycles through the available player colors.
func (s *SetupScreen) cycleColor(idx int) {
	pc := &s.config.Players[idx]

	if pc.Color == nil {
		pc.Color = playerPalette[0]
	}

	// Find current color index
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

// addPlayer adds a new player to the configuration with default settings.
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

	s.init()
	s.refreshLabels()
}

// refreshLabels updates all dynamic UI labels and styles based on current configuration.
func (s *SetupScreen) refreshLabels() {
	for i, pc := range s.config.Players {
		if i >= len(s.playerButtons) {
			continue
		}
		pb := s.playerButtons[i]

		// Determine player display name
		playerName := fmt.Sprintf("Player %d", i+1)
		if pc.Name != "" {
			playerName = pc.Name
		}

		// Determine player color with fallback
		playerColor := pc.Color
		if playerColor == nil {
			playerColor = playerPalette[i%len(playerPalette)]
		}

		// Update card display
		if i < len(s.playerCards) && s.playerCards[i] != nil {
			s.playerCards[i].UpdateFromConfig(ui.PlayerCardConfig{
				Name:     playerName,
				Subtitle: s.roleLabel(pc),
				Symbol:   pc.Symbol,
				Color:    playerColor,
				Ready:    pc.Ready,
			})
		}

		// Update role button label
		if pb.role != nil {
			pb.role.Label = s.roleLabel(pc)
		}
	}

	// Update add player button label
	if s.addPlayerBtn != nil {
		s.addPlayerBtn.Label = fmt.Sprintf("+ Add Player (%d/%d)", len(s.config.Players), maxPlayers)
	}

	// Update start button style based on whether game can start
	if s.startBtn != nil {
		if s.canStartGame() {
			s.startBtn.Style = uiutils.NormalWidgetStyle
		} else {
			s.startBtn.Style = uiutils.DisabledWidgetStyle
		}
	}
}

// roleLabel returns a human-readable label for the player's current role.
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

// canStartGame returns true if all conditions are met to start a game:
// - At least 2 players
// - All players are ready
func (s *SetupScreen) canStartGame() bool {
	if len(s.config.Players) < 2 {
		return false
	}
	for _, pc := range s.config.Players {
		if !pc.Ready {
			return false
		}
	}
	return true
}

// startGame transitions to the game screen if all start conditions are met.
func (s *SetupScreen) startGame() {
	if !s.canStartGame() {
		return
	}
	s.host.SetScreen(NewGameScreen(s.host, s.config))
}

// colorsEqual compares two colors for equality by their RGBA components.
func colorsEqual(a color.Color, b color.Color) bool {
	if a == nil || b == nil {
		return a == b
	}
	ar, ag, ab, aa := a.RGBA()
	br, bg, bb, ba := b.RGBA()
	return ar == br && ag == bg && ab == bb && aa == ba
}
