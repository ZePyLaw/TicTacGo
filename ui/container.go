package ui

import (
	"GoTicTacToe/ui/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

// Element is the minimal interface required by the container to position and render children.
type Element interface {
	SetParentBounds(utils.LayoutRect)
	Update()
	Draw(screen *ebiten.Image)
}

// Container is a layout helper that anchors itself like any other widget
// and lets you place child elements relative to its local space.
type Container struct {
	Widget
	Padding  utils.Insets
	children []Element
}

// NewContainer creates a rectangular container with a pre-rendered background.
func NewContainer(
	offsetX, offsetY float64,
	width, height float64,
	anchor utils.Anchor,
	radius float64,
	style utils.WidgetStyle,
) *Container {
	bg := utils.CreateRoundedRect(int(width), int(height), radius, style.BackgroundNormal)

	return &Container{
		Widget: Widget{
			OffsetX: offsetX,
			OffsetY: offsetY,
			Width:   width,
			Height:  height,
			Anchor:  anchor,
			image:   bg,
			Style:   style,
		},
		Padding: utils.Insets{},
	}
}

// AddChild appends a child element to the container.
func (c *Container) AddChild(child Element) {
	c.children = append(c.children, child)
}

// Update forwards input to children using the container's anchored position.
func (c *Container) Update() {
	_, inner := c.resolveRects()
	c.updateChildren(inner)
}

// Draw renders the container background and all children relative to it.
func (c *Container) Draw(screen *ebiten.Image) {
	outer, inner := c.resolveRects()

	if c.image != nil {
		op := &ebiten.DrawImageOptions{}
		scaleX := outer.Width / float64(c.image.Bounds().Dx())
		scaleY := outer.Height / float64(c.image.Bounds().Dy())
		op.GeoM.Scale(scaleX, scaleY)
		op.GeoM.Translate(outer.X, outer.Y)
		screen.DrawImage(c.image, op)
	}

	for _, child := range c.children {
		child.SetParentBounds(inner)
		child.Draw(screen)
	}
}

func (c *Container) resolveRects() (utils.LayoutRect, utils.LayoutRect) {
	outer := c.LayoutRect()

	inner := utils.LayoutRect{
		X:      outer.X + c.Padding.Left,
		Y:      outer.Y + c.Padding.Top,
		Width:  outer.Width - c.Padding.Horizontal(),
		Height: outer.Height - c.Padding.Vertical(),
	}

	if inner.Width < 0 {
		inner.Width = 0
	}
	if inner.Height < 0 {
		inner.Height = 0
	}

	return outer, inner
}

func (c *Container) updateChildren(inner utils.LayoutRect) {
	for _, child := range c.children {
		child.SetParentBounds(inner)
		child.Update()
	}
}
