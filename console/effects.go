package console

import "github.com/bennicholls/delvetown/util"

//Inverts the foreground and background colours
func Invert(x, y, z int) {
	if util.CheckBounds(x, y, width, height) {
		s := y*width + x
		if grid[s].Z > z {
			return
		}
		f, b := grid[s].ForeColour, grid[s].BackColour
		ChangeBackColour(x, y, f)
		ChangeForeColour(x, y, b)
	}
}
