package console

import "github.com/bennicholls/delveengine/util"

//Inverts the foreground and background colours
func Invert(x, y, z int) {
	if util.CheckBounds(x, y, width, height) {
		s := y*width + x
		if canvas[s].Z > z {
			return
		}
		f, b := canvas[s].ForeColour, canvas[s].BackColour
		ChangeBackColour(x, y, f)
		ChangeForeColour(x, y, b)
	}
}
