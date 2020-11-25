package burl

type drawmode int

const (
	DRAW_GLYPH drawmode = iota
	DRAW_TEXT
)

type Cell struct {
	Glyph      int
	ForeColour uint32
	BackColour uint32
	Z          int
	Dirty      bool
	Border     bool //marks cell as part of a UI Element border.

	//for text rendering mode. TODO:multiple back and fore colours, one for each char
	Mode  drawmode
	Chars [2]int
}

//Sets the properties of a cell all at once for Glyph Mode.
func (c *Cell) SetGlyph(gl int, fore, back uint32, z int) {
	if fore == COL_NONE {
		fore = c.ForeColour
	}
	if back == COL_NONE {
		back = c.BackColour
	}
	if c.Glyph != gl || c.ForeColour != fore || c.BackColour != back || c.Z != z || c.Mode == DRAW_TEXT {
		c.Mode = DRAW_GLYPH
		c.Glyph = gl
		c.ForeColour = fore
		c.BackColour = back
		c.Z = z
		c.Dirty = true
		if gl < GLYPH_BORDER_UD || gl > GLYPH_BORDER_DR {
			c.Border = false
		}
	}
}

//Sets the properties of a cell all at once for Text Mode.
func (c *Cell) SetText(char1, char2 int, fore, back uint32, z int) {
	if fore == COL_NONE {
		fore = c.ForeColour
	}
	if back == COL_NONE {
		back = c.BackColour
	}
	if c.Chars[0] != char1 || c.Chars[1] != char2 || c.ForeColour != fore || c.BackColour != back || c.Z != z || c.Mode == DRAW_GLYPH {
		c.Mode = DRAW_TEXT
		c.Chars[0] = char1
		c.Chars[1] = char2
		c.ForeColour = fore
		c.BackColour = back
		c.Z = z
		c.Dirty = true
	}
}

func (c *Cell) SetBorder(border bool, z int) {
	if c.Z <= z {
		c.Z = z
		c.Dirty = true
		c.Border = border
	}
}

func (c *Cell) CopyCell(c2 *Cell) {
	if c2.Mode == DRAW_GLYPH {
		c.SetGlyph(c2.Glyph, c2.ForeColour, c2.BackColour, c2.Z)
	} else {
		c.SetText(c2.Chars[0], c2.Chars[1], c2.ForeColour, c2.BackColour, c2.Z)
	}
}

//Re-inits a cell back to default blankness.
func (c *Cell) Clear() {
	if c.Mode == DRAW_TEXT {
		c.SetText(32, 32, COL_WHITE, COL_BLACK, 0)
	} else {
		c.SetGlyph(GLYPH_NONE, COL_WHITE, COL_BLACK, 0)
	}
	c.Border = false
}

//Canvas is a z-depthed grid of Cell objects. Cells can have glyph OR text information.
type Canvas struct {
	Cells []Cell

	width, height int
}

func (c *Canvas) Init(w, h int) {
	c.width, c.height = w, h
	c.Cells = make([]Cell, w*h)
}

//Returns the dimensions of the canvas.
func (c *Canvas) Dims() (w, h int) {
	return c.width, c.height
}

//Returns a reference to the cell at (x, y). Returns nil if (x, y) is bad. (outside of canvas)
func (c *Canvas) GetCell(x, y int) *Cell {
	if CheckBounds(x, y, c.width, c.height) {
		return &c.Cells[y*c.width+x]
	}
	return nil
}

//Changes the glyph of a cell in the canvas at position (x, y).
func (c *Canvas) ChangeGlyph(x, y, glyph int) {
	if cell := c.GetCell(x, y); cell != nil {
		cell.SetGlyph(glyph, cell.ForeColour, cell.BackColour, cell.Z)
	}
}

//Changes text of a cell in the canvas at position (x, y).
func (c *Canvas) ChangeText(x, y, z, char1, char2 int) {
	if cell := c.GetCell(x, y); cell != nil && cell.Z <= z {
		cell.SetText(char1, char2, cell.ForeColour, cell.BackColour, z)
	}
}

//Changes a single character on the canvas at position (x,y) in text mode.
//charNum: 0 = Left, 1 = Right (for ease with modulo operations). Throw whatever in here though, it gets modulo 2'd anyways just in case.
func (c *Canvas) ChangeChar(x, y, z, char, charNum int) {
	if cell := c.GetCell(x, y); cell != nil && charNum >= 0 && cell.Z <= z {
		cell.Mode = DRAW_TEXT
		if cell.Chars[charNum%2] != char {
			cell.Chars[charNum%2] = char
			cell.Z = z
			cell.Dirty = true
		}
	}
}

//Changes the foreground drawing colour of a cell in the canvas at position (x, y, z).
func (c *Canvas) ChangeForeColour(x, y, z int, fore uint32) {
	if cell := c.GetCell(x, y); cell != nil && cell.Z <= z {
		if cell.Mode == DRAW_TEXT {
			cell.SetText(cell.Chars[0], cell.Chars[1], fore, cell.BackColour, z)
		} else {
			cell.SetGlyph(cell.Glyph, fore, cell.BackColour, z)
		}
	}
}

//Changes the background colour of a cell in the canvas at position (x, y, z).
func (c *Canvas) ChangeBackColour(x, y, z int, back uint32) {
	if cell := c.GetCell(x, y); cell != nil && cell.Z <= z {
		if cell.Mode == DRAW_TEXT {
			cell.SetText(cell.Chars[0], cell.Chars[1], cell.ForeColour, back, z)
		} else {
			cell.SetGlyph(cell.Glyph, cell.ForeColour, back, z)
		}
	}
}

func (c *Canvas) ChangeColours(x, y, z int, fore, back uint32) {
	if cell := c.GetCell(x, y); cell != nil && cell.Z <= z {
		if cell.Mode == DRAW_TEXT {
			cell.SetText(cell.Chars[0], cell.Chars[1], fore, back, z)
		} else {
			cell.SetGlyph(cell.Glyph, fore, back, z)
		}
	}
}

//Simultaneously changes all characteristics of a glyph cell in the canvas at position (x, y).
//TODO: change name of this to signify it is for changing glyph cells.
func (c *Canvas) ChangeCell(x, y, z, glyph int, fore, back uint32) {
	if cell := c.GetCell(x, y); cell != nil && cell.Z <= z {
		cell.SetGlyph(glyph, fore, back, z)
	}
}

//Draws a string to the console in text mode. CharNum determines which half of the cell we
//start in. See ChageChar() for details.
//TODO: move this to draw.go and include with a suite of primitive drawing functions (circle, square, etc)
func (c *Canvas) DrawText(x, y, z int, txt string, fore, back uint32, charNum int) {
	i := 0 //can't use the index from the range loop since it it counting bytes, not code-points
	for _, char := range txt {
		if CheckBounds(x+(i+charNum)/2, y, c.width, c.height) {
			c.ChangeChar(x+(i+charNum)/2, y, z, int(char), (i+charNum)%2)
			c.ChangeColours(x+(i+charNum)/2, y, z, fore, back)
			if i == len(txt)-1 && (i+charNum)%2 == 0 {
				//if final character is in the left-side of a cell, blank the right side.
				c.ChangeChar(x+(i+charNum)/2, y, z, 32, 1)
			}
		}
		i += 1
	}
}

//TODO: multiple styles.
//Borders work by setting a flag on the cells that need to be borders. At render time, any
//borders with a dirty flag are assigned a border glyph based on the state of their neighbours:
//if the neighbouring cells are on the same z level and also borders, they will connect.
func (c *Canvas) DrawBorder(x, y, z, w, h int, border *Border) {
	if !border.redraw || !border.enabled {
		return
	}

	//Top and bottom.
	for i := -1; i <= w; i++ {
		c.SetCellBorder(x+i, y-1, z, border)
		c.SetCellBorder(x+i, y+h, z, border)
	}
	//Sides
	for i := 0; i < h; i++ {
		c.SetCellBorder(x-1, y+i, z, border)
		c.SetCellBorder(x+w, y+i, z, border)
	}

	//Write centered title.
	if len(border.title) < w && border.title != "" {
		c.DrawText(x+(w/2-len(border.title)/4-1), y-1, z, border.title, border.foreColour, COL_BLACK, 0)
	}

	//Write right-justified hint text
	if border.hint != "" && len(border.hint) < 2*w {
		decoratedHint := TEXT_BORDER_DECO_LEFT + border.hint + TEXT_BORDER_DECO_RIGHT
		offset := w - len(border.hint)/2 - 1
		if len(border.hint)%2 == 1 {
			decoratedHint = TEXT_BORDER_LR + decoratedHint
			offset -= 1
		}

		c.DrawText(x+offset, y+h, z, decoratedHint, border.foreColour, COL_BLACK, 0)
	}

	border.redraw = false
}

func (c *Canvas) SetCellBorder(x, y, z int, border *Border) {
	if cell := c.GetCell(x, y); cell != nil && z >= cell.Z {
		cell.SetBorder(border.enabled, z)
		if border.enabled {
			cell.ForeColour = border.foreColour
			cell.BackColour = border.backColour
			cell.Mode = DRAW_GLYPH
		}
	}
}

//Chooses a border glyph for a cell at (x,y) based on the border state of it's neighbours.
func (c *Canvas) CalcBorderGlyph(x, y int) {
	cell := c.GetCell(x, y)
	if cell == nil {
		return
	}

	var g int
	var u, d, l, r bool

	if uCell := c.GetCell(x, y-1); uCell != nil && uCell.Z == cell.Z {
		u = uCell.Border
	}
	if dCell := c.GetCell(x, y+1); dCell != nil && dCell.Z == cell.Z {
		d = dCell.Border
	}
	if lCell := c.GetCell(x-1, y); lCell != nil && lCell.Z == cell.Z {
		l = lCell.Border
	}
	if rCell := c.GetCell(x+1, y); rCell != nil && rCell.Z == cell.Z {
		r = rCell.Border
	}

	switch {
	case u && d && l && r:
		g = GLYPH_BORDER_UDLR
	case u && d && l:
		g = GLYPH_BORDER_UDL
	case u && d && !l && r:
		g = GLYPH_BORDER_UDR
	case u && !d && l && r:
		g = GLYPH_BORDER_ULR
	case u && !d && !l && r:
		g = GLYPH_BORDER_UR
	case u && !d && l && !r:
		g = GLYPH_BORDER_UL
	case !u && d && l && r:
		g = GLYPH_BORDER_DLR
	case !u && d && l && !r:
		g = GLYPH_BORDER_DL
	case !u && d && !l && r:
		g = GLYPH_BORDER_DR
	case (u || d) && (!l && !r):
		g = GLYPH_BORDER_UD
	case (!u && !d) && (l || r):
		g = GLYPH_BORDER_LR
	}

	if g != 0 {
		cell.Glyph = g
	}
}

//Clears the canvas. Optionally takes args in form (w, h, x, y, z), with z optional, so you can clear
//specific areas of the canvas. If z if provided, it will leave the area at the specified z-level.
func (c *Canvas) Clear(area ...int) {
	w, h, x, y, z := c.width, c.height, 0, 0, 0

	if len(area) >= 4 {
		w, h, x, y = area[0], area[1], area[2], area[3]
	}

	if len(area) == 5 {
		z = area[4]
	}

	for i := 0; i < w*h; i++ {
		ix := x + i%w
		iy := y + i/w
		if cell := c.GetCell(ix, iy); cell != nil {
			cell.Clear()
			cell.Z = z
		}
	}
}

//Fill fills a rect of the console with the provided glyph visuals, at the provided z level.
//TODO: redo this as a drawing function in draw.go (it's basically DrawFilledRect() or something)
func (c *Canvas) Fill(x, y, z, w, h, g int, fore, back uint32) {
	for i := 0; i < w*h; i++ {
		ix := x + i%w
		iy := y + i/w
		c.ChangeCell(ix, iy, z, g, fore, back)
	}
}

//Copies the contents of canvas c2 to canvas c at specified position x,y,z. Any dirty cells copied
//are marked clean.
func (c *Canvas) CopyFromCanvas(x, y, z int, c2 *Canvas) {
	w, h := c2.Dims()
	for i := 0; i < w*h; i++ {
		ix := x + i%w
		iy := y + i/w

		cell2 := c2.GetCell(i%w, i/w)
		if cell2.Border && cell2.Dirty {
			c2.CalcBorderGlyph(i%w, i/w)
		}

		if cell2.Mode == DRAW_GLYPH {
			c.ChangeCell(ix, iy, z+cell2.Z, cell2.Glyph, cell2.ForeColour, cell2.BackColour)
		} else {
			c.ChangeText(ix, iy, z+cell2.Z, cell2.Chars[0], cell2.Chars[1])
			c.ChangeColours(ix, iy, z+cell2.Z, cell2.ForeColour, cell2.BackColour)
		}

		//copy border state
		if cell := c.GetCell(ix, iy); cell != nil {
			cell.SetBorder(cell2.Border, z+cell2.Z)
		}

		cell2.Dirty = false
	}
}
