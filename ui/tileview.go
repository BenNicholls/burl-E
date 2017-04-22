package ui

import "github.com/bennicholls/delvetown/console"
import "github.com/bennicholls/delvetown/util"

//View object for drawing tiles. (eg. maps). Effectively a buffer for drawing before the console grabs it.
type TileView struct {
	Width, Height int
	x, y, z       int
	bordered      bool
	title         string
	visible       bool
	focused       bool

	grid []console.GridCell
}

func NewTileView(w, h, x, y, z int, bord bool) *TileView {
	return &TileView{w, h, x, y, z, bord, "", true, false, make([]console.GridCell, w*h)}
}

func (tv *TileView) SetTitle(s string) {
	tv.title = s
}

//Draws a glyph to the TileView.
func (tv *TileView) Draw(x, y, glyph int, f, b uint32) {
	if util.CheckBounds(x, y, tv.Width, tv.Height) {
		tv.grid[y*tv.Width+x].Set(glyph, f, b, tv.z)
	}
}

//Apply light level. 0-255. TODO: add colour mask (soft orange glow??)
//TODO: This should not be here at all... move this to delvetown.render(), and add a tv.DrawAlpha() that can take an alpha value.
func (tv *TileView) ApplyLight(x, y, b int) {
	if util.CheckBounds(x, y, tv.Width, tv.Height) {
		s := y*tv.Width + x
		if b > 255 {
			b = 255
		} else if b < 80 && b > 0 {
			b = 80 //Baseline brightness for memory... TODO: implement this less magically.
		}
		tv.grid[s].ForeColour = console.ChangeColourAlpha(tv.grid[s].ForeColour, uint8(b))
	}
}

func (tv TileView) Render(offset ...int) {
	if tv.visible {
		offX, offY, offZ := processOffset(offset)
		for i, p := range tv.grid {
			if p.Dirty {
				console.ChangeGridPoint(tv.x+offX+i%tv.Width, tv.y+offY+i/tv.Width, tv.z+offZ, p.Glyph, p.ForeColour, p.BackColour)
				p.Dirty = false
			}
		}
		if tv.bordered {
			console.DrawBorder(tv.x+offX, tv.y+offY, tv.z+offZ, tv.Width, tv.Height, tv.title, tv.focused)
		}
	}
}

func (tv TileView) Dims() (int, int) {
	return tv.Width, tv.Height
}

func (tv TileView) Pos() (int, int, int) {
	return tv.x, tv.y, tv.z
}

//Resets the TileView
func (tv *TileView) Clear() {
	for i, _ := range tv.grid {
		tv.grid[i].Set(0, 0, 0, 0)
		tv.grid[i].Dirty = true
	}
}

func (tv *TileView) ToggleVisible() {
	tv.visible = !tv.visible
	console.Clear()
}

func (tv *TileView) SetVisibility(v bool) {
	tv.visible = v
	console.Clear()
}

func (tv *TileView) ToggleFocus() {
	tv.focused = !tv.focused
}

func (tv *TileView) MoveTo(x, y, z int) {
	tv.x = x
	tv.y = y
	tv.z = z
}
