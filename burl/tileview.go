package burl

//View object for drawing tiles. (eg. maps). Effectively a buffer for drawing before the console grabs it.
type TileView struct {
	UIElement
	grid []Cell
}

func NewTileView(w, h, x, y, z int, bord bool) *TileView {
	tv := new(TileView)
	tv.UIElement = NewUIElement(x, y, z, w, h, bord)
	tv.grid = make([]Cell, w*h)
	tv.Reset()
	return tv
}

func (tv *TileView) SetTitle(s string) {
	tv.title = s
}

//Draws a glyph to the TileView.
func (tv *TileView) Draw(x, y, glyph int, f, b uint32) {
	if CheckBounds(x, y, tv.width, tv.height) {
		tv.grid[y*tv.width+x].SetGlyph(glyph, f, b, tv.z)
	}
}

func (tv *TileView) DrawCircle(x, y, r, glyph int, f, b uint32) {
	DrawCircle(Coord{x, y}, r,
		func(x, y int) {
			tv.Draw(x, y, glyph, f, b)
		})
}

func (tv TileView) Render(offset ...int) {
	if tv.visible {
		offX, offY, offZ := processOffset(offset)
		for i, p := range tv.grid {
			if p.Dirty {
				console.ChangeCell(tv.x+offX+i%tv.width, tv.y+offY+i/tv.width, tv.z+offZ, p.Glyph, p.ForeColour, p.BackColour)
				tv.grid[i].Dirty = false
			}
		}
		tv.UIElement.Render(offX, offY, offZ)
	}
}

//Resets the TileView
func (tv *TileView) Reset() {
	for i, _ := range tv.grid {
		tv.grid[i].Clear()
	}
}
