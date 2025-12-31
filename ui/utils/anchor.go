package utils

type Anchor int

const (
	AnchorTopLeft Anchor = iota
	AnchorTopCenter
	AnchorTopRight
	AnchorCenterLeft
	AnchorCenter
	AnchorCenterRight
	AnchorBottomLeft
	AnchorBottomCenter
	AnchorBottomRight
)

// ComputeAnchoredPosition computes the position of an anchored widget.
// returns the x and y coordinates of the widget's top-left corner.
func ComputeAnchoredPosition(
	anchor Anchor,
	offsetX, offsetY float64,
	w, h float64,
	parentW, parentH float64,
) (float64, float64) {

	x := offsetX
	y := offsetY

	switch anchor {
	case AnchorTopCenter:
		x = parentW/2 - w/2 + offsetX
	case AnchorTopRight:
		x = parentW - w + offsetX

	case AnchorCenterLeft:
		y = parentH/2 - h/2 + offsetY
	case AnchorCenter:
		x = parentW/2 - w/2 + offsetX
		y = parentH/2 - h/2 + offsetY
	case AnchorCenterRight:
		x = parentW - w + offsetX
		y = parentH/2 - h/2 + offsetY

	case AnchorBottomLeft:
		y = parentH - h + offsetY
	case AnchorBottomCenter:
		x = parentW/2 - w/2 + offsetX
		y = parentH - h + offsetY
	case AnchorBottomRight:
		x = parentW - w + offsetX
		y = parentH - h + offsetY
	}

	return x, y
}
