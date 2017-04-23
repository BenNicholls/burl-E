package console

import "github.com/veandco/go-sdl2/sdl"
import "github.com/bennicholls/delveengine/util"
import "fmt"
import "math/rand"
import "errors"

var window *sdl.Window
var renderer *sdl.Renderer
var sprites *sdl.Texture
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
}

func (g *Cell) Set(gl int, fore, back uint32, z int) {
	if g.Glyph != gl || g.ForeColour != fore || g.BackColour != back || g.Z != z {
		g.Glyph = gl
		g.ForeColour = fore
		g.BackColour = back
		g.Z = z
		g.Dirty = true
	}
}

func (g *Cell) Clear() {
	g.Set(0, 0, 0, 0)
}

//Setup the game window, renderer, etc
func Setup(w, h int, spritesheet, title string) error {

	width = w
	height = h
	var err error

	//load spritesheet first so we can infer tileSize
	image, err := sdl.LoadBMP(spritesheet)
	if err != nil {
		return errors.New("Failed to load image: " + fmt.Sprint(sdl.GetError()))
	}
	defer image.Free()
	image.SetColorKey(1, 0xFF00FF)
	tileSize = int(image.W / 16)

	window, err = sdl.CreateWindow(title, sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, width*tileSize, height*tileSize, sdl.WINDOW_OPENGL)
	if err != nil {
		return errors.New("Failed to create window: " + fmt.Sprint(sdl.GetError()))
	}

	//manually set pixelformat to ARGB (window defaults to RGB for some reason)
	format, err = sdl.AllocFormat(uint(sdl.PIXELFORMAT_ARGB8888))
	if err != nil {
		return errors.New("No pixelformat: " + fmt.Sprint(sdl.GetError()))
	}

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_SOFTWARE)
	if err != nil {
		return errors.New("Failed to create renderer: " + fmt.Sprint(sdl.GetError()))
	}
	renderer.Clear()

	sprites, err = renderer.CreateTextureFromSurface(image)
	if err != nil {
		return errors.New("Failed to create sprite texture: " + fmt.Sprint(sdl.GetError()))
	}
	err = sprites.SetBlendMode(sdl.BLENDMODE_BLEND)
	if err != nil {
		return errors.New("Failed to set blendmode: " + fmt.Sprint(sdl.GetError()))
	}

	canvas = make([]Cell, width*height)
	masterDirty = true

	frames = 0
	frameTime, ticks = 0, 0
	fps = 35 //17ms = 60 FPS approx
	showFPS = false
	BorderColour1 = 0xFFE28F00
	BorderColour2 = 0xFF555555

	return nil
}

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
				dst = makeRect((i%width)*tileSize, (i/width)*tileSize, tileSize, tileSize)
				src = makeRect((s.Glyph%16)*tileSize, (s.Glyph/16)*tileSize, tileSize, tileSize)

				renderer.SetDrawColor(sdl.GetRGBA(s.BackColour, format))
				renderer.FillRect(&dst)

				r, g, b, a := sdl.GetRGBA(s.ForeColour, format)
				sprites.SetColorMod(r, g, b)
				sprites.SetAlphaMod(a)
				renderer.Copy(sprites, &src, &dst)

				canvas[i].Dirty = false
			}
		}

		renderer.Present()
		masterDirty = false
	}

	//framerate limiter, so my cpu doesn't implode
	ticks = sdl.GetTicks() - frameTime
	if ticks < fps {
		sdl.Delay(fps - ticks)
	}
	frameTime = sdl.GetTicks()
	frames++
}

func SetFramerate(f uint32) {
	fps = f
}

//int32 for rect arguments. what a world.
func makeRect(x, y, w, h int) sdl.Rect {
	return sdl.Rect{int32(x), int32(y), int32(w), int32(h)}
}

func Cleanup() {
	format.Free()
	sprites.Destroy()
	renderer.Destroy()
	window.Destroy()
}

func ChangeGlyph(x, y, glyph int) {
	if util.CheckBounds(x, y, width, height) && canvas[y*width+x].Glyph != glyph {
		canvas[y*width+x].Glyph = glyph
		canvas[y*width+x].Dirty = true
		masterDirty = true
	}
}

func ChangeForeColour(x, y int, fore uint32) {
	if util.CheckBounds(x, y, width, height) && canvas[y*width+x].ForeColour != fore {
		canvas[y*width+x].ForeColour = fore
		canvas[y*width+x].Dirty = true
		masterDirty = true
	}
}

func ChangeBackColour(x, y int, back uint32) {
	if util.CheckBounds(x, y, width, height) && canvas[y*width+x].BackColour != back {
		canvas[y*width+x].BackColour = back
		canvas[y*width+x].Dirty = true
		masterDirty = true
	}
}

func ToggleFPS() {
	showFPS = !showFPS
}

func ChangeCell(x, y, z, glyph int, fore, back uint32) {
	s := y*width + x
	if util.CheckBounds(x, y, width, height) && canvas[s].Z <= z {
		canvas[s].Set(glyph, fore, back, z)
		masterDirty = true
	}
}

//TODO: custom colouring, multiple styles
func DrawBorder(x, y, z, w, h int, title string, focused bool) {
	bc := BorderColour1
	if !focused {
		bc = BorderColour2
	}
	for i := 0; i < w; i++ {
		ChangeCell(x+i, y-1, z, 0xc4, bc, 0xFF000000)
		ChangeCell(x+i, y+h, z, 0xc4, bc, 0xFF000000)
	}
	for i := 0; i < h; i++ {
		ChangeCell(x-1, y+i, z, 0xb3, bc, 0xFF000000)
		ChangeCell(x+w, y+i, z, 0xb3, bc, 0xFF000000)
	}
	ChangeCell(x-1, y-1, z, 0xda, bc, 0xFF000000)
	ChangeCell(x-1, y+h, z, 0xc0, bc, 0xFF000000)
	ChangeCell(x+w, y+h, z, 0xd9, bc, 0xFF000000)
	ChangeCell(x+w, y-1, z, 0xbf, bc, 0xFF000000)

	if len(title) < w && title != "" {
		for i, r := range title {
			ChangeCell(x+(w/2-len(title)/2)+i, y-1, z, int(r), 0xFFFFFFFF, 0xFF000000)
		}
	}
}

//Optionally takes a rect so you can clear specific areas of the console
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

func Dims() (w, h int) {
	return width, height
}

//Test function.
func SpamGlyphs() {
	for n := 0; n < 100; n++ {
		x := rand.Intn(width)
		y := rand.Intn(height)
		ChangeCell(x, y, 0, rand.Intn(255), sdl.MapRGBA(format, 0, 255, 0, 50), sdl.MapRGBA(format, 255, 0, 0, 255))
	}
}

func MakeColour(r, g, b int) uint32 {
	return sdl.MapRGBA(format, uint8(r), uint8(g), uint8(b), 255)
}

func ChangeColourAlpha(c uint32, a uint8) uint32 {
	r, g, b := sdl.GetRGB(c, format)
	return sdl.MapRGBA(format, r, g, b, a)
}
