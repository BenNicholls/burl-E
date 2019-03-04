package burl

import "github.com/bennicholls/burl-E/reximage"

//View object for drawing tiles. (eg. maps). Effectively a buffer for drawing before the console grabs it.
type TileView struct {
	UIElement
	grid []Cell
}

func NewTileView(w, h, x, y, z int, bord bool) *TileView {
	tv := new(TileView)
	tv.UIElement = NewUIElement(w, h, x, y, z, bord)
	tv.grid = make([]Cell, w*h)
	tv.Reset()
	return tv
}

//Draws a glyph to the TileView. if COL_NONE is passed as a parameter, uses the colour that was previously there.
func (tv *TileView) Draw(x, y, glyph int, f, b uint32) {
	if CheckBounds(x, y, tv.width, tv.height) {
		if f == COL_NONE {
			f = tv.grid[y*tv.width+x].ForeColour
		}
		if b == COL_NONE {
			b = tv.grid[y*tv.width+x].BackColour
		}
		tv.grid[y*tv.width+x].SetGlyph(glyph, f, b, tv.z)
	}
}

//Draws a drawable object on the tileview at coord (x, y). If (x, y) not in bounds, does nothing.
func (tv *TileView) DrawObject(x, y int, d Drawable) {
	tv.Draw(x, y, d.GetVisuals().Glyph, d.GetVisuals().ForeColour, d.GetVisuals().BackColour)
}

func (tv *TileView) DrawCircle(x, y, r, glyph int, f, b uint32) {
	DrawCircle(Coord{x, y}, r,
		func(x, y int) {
			tv.Draw(x, y, glyph, f, b)
		})
}

//draws a palette to the tileview, one colour per tile. stops when it hits the edge of the view object. dir is HORIZONTAL or VERTICAL.
func (tv *TileView) DrawPalette(x, y int, p Palette, dir int) {
	for i, c := range p {
		if dir == HORIZONTAL {
			if x + i >= tv.width {
				break
			}
			tv.grid[y*tv.width+x+i].SetGlyph(GLYPH_FILL, c, COL_BLACK, 0)
		} else {
			if y + i >= tv.height {
				break
			}
			tv.grid[(y+i)*tv.width+x].SetGlyph(GLYPH_FILL, c, COL_BLACK, 0)
		}
	}
}

func (tv TileView) Render() {
	if tv.visible {
		for i, p := range tv.grid {
			console.ChangeCell(tv.x+i%tv.width, tv.y+i/tv.width, tv.z, p.Glyph, p.ForeColour, p.BackColour)
		}
		tv.UIElement.Render()
	}
}

//Resets the TileView
func (tv *TileView) Reset() {
	for i := range tv.grid {
		tv.grid[i].Clear()
	}
}

func (tv *TileView) LoadImageFromXP(filename string) {
	imageData, err := reximage.Import(filename)
	if err != nil {
		LogError("Error loading image " + filename + ": " + err.Error())
	}

	for j := 0; j < imageData.Height; j++ {
		for i := 0; i < imageData.Width; i++ {
			cell, _ := imageData.GetCell(i, j) //cell from imagedata
			g := int(cell.Glyph)
			fore, back := cell.ARGB()
			tv.grid[i+j*tv.width].SetGlyph(g, fore, back, tv.z)
		}
	}

}

