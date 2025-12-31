package utils

// SizeMode describes how a widget should resolve its size relative to its parent.
type SizeMode int

const (
	// SizeFixed uses the widget's own Width/Height values.
	SizeFixed SizeMode = iota
	// SizeFill expands the widget to fill the parent's available space on that axis.
	SizeFill
)

// LayoutRect represents a resolved rectangle in absolute coordinates.
type LayoutRect struct {
	X, Y          float64
	Width, Height float64
}

// Insets is a CSS-like padding/margin helper.
type Insets struct {
	Top, Right, Bottom, Left float64
}

// InsetsAll creates uniform insets on all sides.
func InsetsAll(v float64) Insets {
	return Insets{Top: v, Right: v, Bottom: v, Left: v}
}

// Horizontal returns the combined horizontal space.
func (i Insets) Horizontal() float64 {
	return i.Left + i.Right
}

// Vertical returns the combined vertical space.
func (i Insets) Vertical() float64 {
	return i.Top + i.Bottom
}
