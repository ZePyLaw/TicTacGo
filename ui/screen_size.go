package ui

var (
	screenWidth  int
	screenHeight int
)

// UpdateScreenSize stores the current drawable screen size so widgets can position themselves.
func UpdateScreenSize(w, h int) {
	screenWidth = w
	screenHeight = h
}

func currentScreenSize() (int, int) {
	return screenWidth, screenHeight
}
