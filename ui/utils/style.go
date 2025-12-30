package utils

import "image/color"

// WidgetStyle represents the style of a widget.
type WidgetStyle struct {

	// Background
	BackgroundNormal color.Color
	BackgroundHover  color.Color

	// Text color
	TextColor color.Color

	// Border colo
	BorderColor color.Color
	BorderWidth float64

	// Hover behaviour
	HoverMode HoverMode
}

// DefaultWidgetStyle is the default style for widgets.
var DefaultWidgetStyle = WidgetStyle{
	BackgroundNormal: color.RGBA{40, 40, 40, 255},
	BackgroundHover:  color.RGBA{70, 70, 70, 255},

	TextColor: color.White,

	BorderColor: color.RGBA{255, 255, 255, 50},
	BorderWidth: 2,

	HoverMode: HoverColorLerp,
}

// TransparentWidgetStyle is the style for transparent widgets.
var TransparentWidgetStyle = WidgetStyle{
	BackgroundNormal: color.RGBA{0, 0, 0, 0},
	BackgroundHover:  color.RGBA{255, 255, 255, 20},

	TextColor: color.White,

	BorderColor: color.RGBA{255, 255, 255, 100},
	BorderWidth: 1,

	HoverMode: HoverFade,
}

var NormalWidgetStyle = WidgetStyle{
	BackgroundNormal: color.RGBA{64, 92, 245, 100},
	BackgroundHover:  color.RGBA{43, 61, 163, 255},

	TextColor: color.White,

	BorderColor: color.RGBA{60, 0, 0, 255},
	BorderWidth: 2,

	HoverMode: HoverColorLerp,
}

// DangerWidgetStyle is the style for danger widgets.
var DangerWidgetStyle = WidgetStyle{
	BackgroundNormal: color.RGBA{200, 30, 30, 255},
	BackgroundHover:  color.RGBA{200, 30, 30, 255},

	TextColor: color.White,

	BorderColor: color.RGBA{60, 0, 0, 255},
	BorderWidth: 2,

	HoverMode: HoverSolid,
}

// SuccessWidgetStyle is the style for positive/confirm actions (e.g., ready).
var SuccessWidgetStyle = WidgetStyle{
	BackgroundNormal: color.RGBA{40, 170, 80, 255},
	BackgroundHover:  color.RGBA{30, 140, 60, 255},

	TextColor: color.White,

	BorderColor: color.RGBA{0, 80, 0, 255},
	BorderWidth: 2,

	HoverMode: HoverColorLerp,
}
