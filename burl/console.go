package burl

//Console is the framebuffer which composites the UI of the current active State. This happens once
//per frame before being sent to the Renderer to be displayed.
type Console struct {
	Canvas

	redraw bool
	defaultBorderStyle Border
	//TODO sync.RWmutex?? i think??
}

//Setup the game window, renderer, etc
func (c *Console) Setup(w, h int) (err error) {
	c.Init(w, h)
	c.Clear()

	c.defaultBorderStyle = Border{
		foreColour: COL_LIGHTGREY,
		backColour: COL_BLACK,
	}

	return nil
}

//Builds the final frame. Renders top-down to eliminate overwriting.
func (c *Console) BuildFrame() {

	if c.redraw {
		c.redraw = false
		c.Clear()
	}

	if debug && debugger.IsVisible() {
		x, y, z := debugger.Pos()
		c.CopyFromCanvas(x, y, z, debugger.GetCanvas())
		c.DrawBorder(x, y, z, debugger.width, debugger.height, debugger.GetBorder())
	}

	if d := gameState.GetDialog(); d != nil {
		x, y, z := d.GetWindow().Pos()
		c.CopyFromCanvas(x, y, z, d.GetWindow().GetCanvas())
	}

	if w := gameState.GetWindow(); w != nil { //is this check stupid????
		x, y, z := w.Pos()
		c.CopyFromCanvas(x, y, z, w.GetCanvas())
	}

	//Post processing. Probably a bunch of stuff here, for now just border linking i guess.
	for i := range c.Cells {
		ix, iy := i%c.width, i/c.width
		cell := c.GetCell(ix, iy)

		if cell.Dirty && cell.Border {
			c.CalcBorderGlyph(ix, iy)
		}
	}
}
