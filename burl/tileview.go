package burl

import "github.com/bennicholls/burl-E/reximage"

//View object for drawing tiles. (eg. maps). Effectively a buffer for drawing before the console grabs it.
type TileView struct {
	UIElement
}

func NewTileView(w, h, x, y, z int, bord bool) *TileView {
	tv := new(TileView)
	tv.UIElement = NewUIElement(w, h, x, y, z, bord)
	return tv
}

//Draws a glyph to the TileView. if COL_NONE is passed as a parameter, uses the colour that was previously there.
//THINK: should this take a Z value? then you could draw in layers
func (tv *TileView) Draw(x, y, glyph int, f, b uint32) {
	if cell := tv.GetCell(x, y); cell != nil {
		if f == COL_NONE {
			f = cell.ForeColour
		}
		if b == COL_NONE {
			b = cell.BackColour
		}
		tv.ChangeCell(x, y, cell.Z, glyph, f, b)
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
			if x+i >= tv.width {
				break
			}
			tv.ChangeCell(x+i, y, 0, GLYPH_NONE, COL_NONE, c)
		} else {
			if y+i >= tv.height {
				break
			}
			tv.ChangeCell(x, y+i, 0, GLYPH_NONE, COL_NONE, c)
		}
	}
}

func (tv *TileView) LoadImageFromXP(filename string) {
	imageData, err := reximage.Import(filename)
	if err != nil {
		LogError("Error loading image " + filename + ": " + err.Error())
		return
	}

	for y := 0; y < imageData.Height; y++ {
		for x := 0; x < imageData.Width; x++ {
			cell, _ := imageData.GetCell(x, y) //cell from imagedata
			g := int(cell.Glyph)
			fore, back := cell.ARGB()
			tv.ChangeCell(x, y, 0, g, fore, back)
		}
	}

}
