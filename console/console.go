package console

import "github.com/veandco/go-sdl2/sdl"
import "github.com/bennicholls/delveengine/util"
import "fmt"
import "math/rand"
import "errors"

var window *sdl.Window
var renderer *sdl.Renderer
var glyphs *sdl.Texture
var font *sdl.Texture
var format *sdl.PixelFormat

var width, height, tileSize int

var canvas []Cell
var masterDirty bool //is this necessary?

var frameTime, ticks, fps uint32
var frames int
var showFPS bool

//Border colours are defined here so we can change them program-wide,
//for reasons that I hope will come in handy later.
var BorderColour1 uint32 //focused element colour
var BorderColour2 uint32 //unfocused element colour

type Cell struct {
	Glyph      int
	ForeColour uint32
	BackColour uint32
	Z          int
	Dirty      bool

	//for text rendering mode. TODO:multiple back and fore colours, one for each char
	TextMode bool
	Chars [2]int
}

//Sets the properties of a cell all at once for Glyph Mode.
func (c *Cell) SetGlyph(gl int, fore, back uint32, z int) {
	if c.Glyph != gl || c.ForeColour != fore || c.BackColour != back || c.Z != z || c.TextMode {
		c.TextMode = false
		c.Glyph = gl
		c.ForeColour = fore
		c.BackColour = back
		c.Z = z
		c.Dirty = true
	}
}

//Sets the properties of a cell all at once for Text Mode.
func (c *Cell) SetText(char1, char2 int, fore, back uint32, z int) {
	if c.Chars[0] != char1 || c.Chars[1] != char2 || c.ForeColour != fore || c.BackColour != back || c.Z != z || c.TextMode == false {
		c.TextMode = true
		c.Chars[0] = char1
		c.Chars[1] = char2 
		c.ForeColour = fore
		c.BackColour = back
		c.Z = z
		c.Dirty = true
	}
}

//Re-inits a cell back to default. Defaults to Glyph Mode.
func (c *Cell) Clear() {
	c.SetGlyph(0, 0, 0, 0)
}

//Setup the game window, renderer, etc
//TODO: extraact image loading to its own function, resizable window.
func Setup(w, h int, glyphPath, fontPath, title string) error {
	width = w
	height = h
	var err error

	tileSize = 16

	window, err = sdl.CreateWindow(title, sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, width*tileSize, height*tileSize, sdl.WINDOW_OPENGL)
	if err != nil {
		return errors.New("Failed to create window: " + fmt.Sprint(sdl.GetError()))
	}

	//manually set pixelformat to ARGB (window defaults to RGB for some reason)
	format, err = sdl.AllocFormat(uint(sdl.PIXELFORMAT_ARGB8888))
	if err != nil {
		return errors.New("No pixelformat: " + fmt.Sprint(sdl.GetError()))
	}

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_SOFTWARE) //Software renderer because ACCELERATED borks my laptop for some reason.
	if err != nil {
		return errors.New("Failed to create renderer: " + fmt.Sprint(sdl.GetError()))
	}
	renderer.Clear()

	canvas = make([]Cell, width*height)
	masterDirty = true

	//init drawing fonts
	err = ChangeFonts(glyphPath, fontPath)
	if err != nil {
		return nil
	}

	frames = 0
	frameTime, ticks = 0, 0
	fps = 17 //17ms = 60 FPS approx
	showFPS = false
	BorderColour1 = 0xFFE28F00
	BorderColour2 = 0xFF555555

	return nil
}

func ChangeFonts(glyphPath, fontPath string) error {
	var err error
	if glyphs != nil {
		glyphs.Destroy()
	}
	glyphs, err = LoadTexture(glyphPath)
	if err != nil {
		return err
	}
	if font != nil {
		font.Destroy()
	}
	font, err = LoadTexture(fontPath)
	if err != nil {
		return err
	}
	Clear()

	_, _, gw, _, _ := glyphs.Query()
	
	//reset window size if fontsize changed
	if int(gw/16) != tileSize {
		tileSize = int(gw/16)
		window.SetSize(tileSize*width, tileSize*height)
	}

	return nil
}

//Loads a bmp font into the GPU using the current window renderer.
//TODO: support more than bmps?
func LoadTexture(path string) (*sdl.Texture, error) {
	image, err := sdl.LoadBMP(path)
	defer image.Free()
	if err != nil {
		return nil, errors.New("Failed to load image: " + fmt.Sprint(sdl.GetError()))
	}
	image.SetColorKey(1, 0xFF00FF)
	texture, err := renderer.CreateTextureFromSurface(image)
	if err != nil {
		return nil, errors.New("Failed to create texture: " + fmt.Sprint(sdl.GetError()))
	}
	err = texture.SetBlendMode(sdl.BLENDMODE_BLEND)
	if err != nil {
		texture.Destroy()
		return nil, errors.New("Failed to set blendmode: " + fmt.Sprint(sdl.GetError()))
	}

	return texture, nil
}

//Renders the canvas to the GPU and flips the buffer.
func Render() {
	//render fps counter
	if showFPS {
		fpsString := fmt.Sprintf("%d fps\n", frames*1000/int(sdl.GetTicks()))
		for i, r := range fpsString {
			ChangeCell(i, 0, 10, int(r), 0xFF00FF00, 0xFFFF0000)
		}
	}

	//render the scene!
	if masterDirty {
		var src, dst sdl.Rect

		for i, s := range canvas {
			if s.Dirty {
				if s.TextMode {
					for c_i, c := range s.Chars {
						dst = makeRect((i%width)*tileSize + c_i*tileSize/2, (i/width)*tileSize, tileSize/2, tileSize)
						src = makeRect((c%32)*tileSize/2, (c/32)*tileSize, tileSize/2, tileSize)
						CopyToRenderer(s, font, src, dst)
					}
				} else {
					dst = makeRect((i%width)*tileSize, (i/width)*tileSize, tileSize, tileSize)
					src = makeRect((s.Glyph%16)*tileSize, (s.Glyph/16)*tileSize, tileSize, tileSize)
					CopyToRenderer(s, glyphs, src, dst)
				}

				canvas[i].Dirty = false
			}
		}

		renderer.Present()
		masterDirty = false
	}

	//framerate limiter, so the cpu doesn't implode
	//TODO: option to turn this off? I guess you can set the fps arbitrarily high...
	//NOTE: should this be the responsibility of the main game loop? 
	ticks = sdl.GetTicks() - frameTime
	if ticks < fps {
		sdl.Delay(fps - ticks)
	}
	frameTime = sdl.GetTicks()
	frames++
}

//Copies a rect of pixeldata from a source texture to a rect on the renderer texture for the console.
func CopyToRenderer(c Cell, tex *sdl.Texture, src, dst sdl.Rect) {
	renderer.SetDrawColor(sdl.GetRGBA(c.BackColour, format)) //should NOT be doing this every cell.
	renderer.FillRect(&dst)
	r, g, b, a := sdl.GetRGBA(c.ForeColour, format) //should NOT be doing this every cell.
	tex.SetColorMod(r, g, b)
	tex.SetAlphaMod(a)
	renderer.Copy(tex, &src, &dst)
}

//Sets maximum framerate as enforced by the framerate limiter. NOTE: cannot go higher than 1000 fps.
func SetFramerate(f int) {
	fps = uint32(1000/f) + 1
}


//Toggles rendering of the FPS meter.
func ToggleFPS() {
	showFPS = !showFPS
}

//int32 for rect arguments. what a world.
func makeRect(x, y, w, h int) sdl.Rect {
	return sdl.Rect{int32(x), int32(y), int32(w), int32(h)}
}

//Deletes special graphics structures, closes files, etc. Defer this function!
func Cleanup() {
	format.Free()
	glyphs.Destroy()
	font.Destroy()
	renderer.Destroy()
	window.Destroy()
}

//Changes the glyph of a cell in the canvas at position (x, y).
func ChangeGlyph(x, y, glyph int) {
	if util.CheckBounds(x, y, width, height) && canvas[y*width+x].Glyph != glyph {
		canvas[y*width+x].TextMode = false
		canvas[y*width+x].Glyph = glyph
		canvas[y*width+x].Dirty = true
		masterDirty = true
	}
}

//Changes text of a cell in the canvas at position (x, y).
func ChangeText(x, y, char1, char2 int) {
	if util.CheckBounds(x, y, width, height) {
		canvas[y*width+x].TextMode = true
		canvas[y*width+x].Chars[0] = char1
		canvas[y*width+x].Chars[1] = char2
		canvas[y*width+x].Dirty = true

	}
}

//Changes the foreground drawing colour of a cell in the canvas at position (x, y).
func ChangeForeColour(x, y int, fore uint32) {
	if util.CheckBounds(x, y, width, height) && canvas[y*width+x].ForeColour != fore {
		canvas[y*width+x].ForeColour = fore
		canvas[y*width+x].Dirty = true
		masterDirty = true
	}
}

//Changes the background colour of a cell in the canvas at position (x, y).
func ChangeBackColour(x, y int, back uint32) {
	if util.CheckBounds(x, y, width, height) && canvas[y*width+x].BackColour != back {
		canvas[y*width+x].BackColour = back
		canvas[y*width+x].Dirty = true
		masterDirty = true
	}
}

//Simultaneously changes all characteristics of a glyph cell in the canvas at position (x, y).
//TODO: change name of this to signify it is for changing glyph cells.
func ChangeCell(x, y, z, glyph int, fore, back uint32) {
	s := y*width + x
	if util.CheckBounds(x, y, width, height) && canvas[s].Z <= z {
		canvas[s].SetGlyph(glyph, fore, back, z)
		masterDirty = true
	}
}

//Draws a string to the console in text mode.
func DrawText(x, y, z int, txt string, fore, back uint32) {
	for i, c := range txt {
		cell := canvas[y*width + x + i/2]
		if i % 2 == 0 {
			ChangeText(x + i/2, y, int(c), cell.Chars[1])
		} else {
			ChangeText(x + i/2, y, cell.Chars[0], int(c))
		}
		ChangeForeColour(x + i/2, y, fore)
		ChangeBackColour(x + i/2, y, back)
	}
}

//TODO: custom colouring, multiple styles. 
//NOTE: current border colouring thing is a bit of a hack. Need to add actual support for
//border and ui styling. (Should this be in delveengine/ui??? hmmm.)
func DrawBorder(x, y, z, w, h int, title string, focused bool) {
	//set border colour.
	bc := BorderColour1
	if !focused {
		bc = BorderColour2
	}
	//Top and bottom.
	for i := 0; i < w; i++ {
		ChangeCell(x+i, y-1, z, 0xc4, bc, 0xFF000000)
		ChangeCell(x+i, y+h, z, 0xc4, bc, 0xFF000000)
	}
	//Sides
	for i := 0; i < h; i++ {
		ChangeCell(x-1, y+i, z, 0xb3, bc, 0xFF000000)
		ChangeCell(x+w, y+i, z, 0xb3, bc, 0xFF000000)
	}
	//corners
	ChangeCell(x-1, y-1, z, 0xda, bc, 0xFF000000)
	ChangeCell(x-1, y+h, z, 0xc0, bc, 0xFF000000)
	ChangeCell(x+w, y+h, z, 0xd9, bc, 0xFF000000)
	ChangeCell(x+w, y-1, z, 0xbf, bc, 0xFF000000)

	//Write centered title.
	if len(title) < w && title != "" {
		DrawText(x+(w/2 - len(title)/4 - 1), y-1, z, title, 0xFFFFFFFF, 0xFF000000)
	}
}

//Clears an area of the canvas. Optionally takes a rect (defined by 4 ints) so you can clear specific areas of the console
func Clear(rect ...int) {

	offX, offY, w, h := 0, 0, width, height

	if len(rect) == 4 {
		offX, offY, w, h = rect[0], rect[1], rect[2], rect[3]
	}

	for i := 0; i < w*h; i++ {
		x := offX + i%w
		y := offY + i/w
		canvas[y*width+x].Clear()
	}
}

//Returns the dimensions of the canvas.
func Dims() (w, h int) {
	return width, height
}

//Test function. Changes 100 glyphs randomly each frame.
func SpamGlyphs() {
	for n := 0; n < 100; n++ {
		x := rand.Intn(width)
		y := rand.Intn(height)
		ChangeCell(x, y, 0, rand.Intn(255), sdl.MapRGBA(format, 0, 255, 0, 50), sdl.MapRGBA(format, 255, 0, 0, 255))
	}
}

//Takes r,g,b ints and creates a colour as defined by the pixelformat with alpha 255. 
//TODO: rgba version of this function? variatic function that can optionally take an alpha? Hmm.
func MakeColour(r, g, b int) uint32 {
	return sdl.MapRGBA(format, uint8(r), uint8(g), uint8(b), 255)
}

//Changes alpha of a colour.
func ChangeColourAlpha(c uint32, a uint8) uint32 {
	r, g, b := sdl.GetRGB(c, format)
	return sdl.MapRGBA(format, r, g, b, a)
}
