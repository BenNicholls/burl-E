package burl

//Inverts the foreground and background colours
func (c *Console) Invert(x, y, z int) {
	if CheckBounds(x, y, c.width, c.height) {
		s := y*c.width + x
		if c.canvas[s].Z <= z {
			f, b := c.canvas[s].ForeColour, c.canvas[s].BackColour
			c.ChangeColours(x, y, z, b, f)
		}
	}
}
