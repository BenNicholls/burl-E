package burl

import "encoding/csv"
import "os"
import "strconv"
import "strings"

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
			console.ChangeCell(tv.x+offX+i%tv.width, tv.y+offY+i/tv.width, tv.z+offZ, p.Glyph, p.ForeColour, p.BackColour)
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

//Loads an image exported from RexPaint into CSV format.
//TODO: write an actual .xp library for Go! Doesn't seem to be one!
func (tv *TileView) LoadImageFromCSV(filename string) {
	if !strings.HasSuffix(filename, ".csv") {
		LogError("Cannot load image " + filename + "(not csv file!)")
		return
	}

	//open image file
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		LogError("Could not load image data: " + filename)
		return
	}

	//read records/fields from csv into string[][]
	data, err := csv.NewReader(f).ReadAll()
	if err != nil {
		LogError("Could not read csv data: " + filename)
		return
	}

	//parse records (record 0 is header data)
	for i := 1; i < len(data); i++ {
		x, _ := strconv.ParseInt(data[i][0], 10, 0)
		y, _ := strconv.ParseInt(data[i][1], 10, 0)

		if int(x) >= tv.width || int(y) >= tv.height {
			continue
		}

		glyph, _ := strconv.ParseInt(data[i][2], 10, 0)
		f, _ := strconv.ParseInt(data[i][3][1:], 16, 0)
		fore := ChangeAlpha(uint32(f), 0xFF)
		b, _ := strconv.ParseInt(data[i][4][1:], 16, 0)
		back := ChangeAlpha(uint32(b), 0xFF)

		tv.grid[int(x)+tv.width*int(y)].SetGlyph(int(glyph), fore, back, tv.z)
	}
}
