package burl

import (
	"errors"
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type Renderer interface {
	Setup(string, string, string) error
	Ready() bool
	Cleanup()
	ChangeFonts(string, string) error
	SetFullscreen(bool)
	ToggleFullscreen()
	SetFramerate(int)
	Render()
	ForceRedraw()
	ToggleDebugMode(string)
}

type SDLRenderer struct {
	window       *sdl.Window
	renderer     *sdl.Renderer
	glyphs       *sdl.Texture
	font         *sdl.Texture
	canvasBuffer *sdl.Texture

	tileSize int

	forceRedraw bool
	showFPS     bool
	showChanges bool

	frameTime               time.Time
	frameTargetDur, elapsed time.Duration
	frames                  int

	//store render colours so we don't have to set them for every renderer.Copy()
	backDrawColour      uint32
	foreDrawColourText  uint32
	foreDrawColourGlyph uint32

	ready bool
}

func (sdlr *SDLRenderer) Setup(glyphPath, fontPath, title string) (err error) {
	//renderer defaults to 800x600, once fonts are loaded it figured out the resolution to use and resizes accordingly
	sdlr.window, err = sdl.CreateWindow(title, sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, 800, 600, sdl.WINDOW_OPENGL)
	if err != nil {
		LogError("SDL Renderer: Failed to create window. sdl:" + fmt.Sprint(sdl.GetError()))
		return errors.New("Failed to create window.")
	}

	sdlr.renderer, err = sdl.CreateRenderer(sdlr.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		LogError("SDL Renderer: Failed to create renderer. sdl:" + fmt.Sprint(sdl.GetError()))
		return errors.New("Failed to create renderer.")
	}
	sdlr.renderer.Clear()

	err = sdlr.ChangeFonts(glyphPath, fontPath)
	if err != nil {
		return err
	}

	sdlr.SetFramerate(60)

	sdlr.ready = true
	return
}

func (sdlr *SDLRenderer) Ready() bool {
	return sdlr.ready
}

//Deletes special graphics structures, closes files, etc. Defer this function!
func (sdlr *SDLRenderer) Cleanup() {
	sdlr.glyphs.Destroy()
	sdlr.font.Destroy()
	sdlr.canvasBuffer.Destroy()
	sdlr.renderer.Destroy()
	sdlr.window.Destroy()
}

//Loads new fonts to the renderer and changes the tilesize (and by extension, the window size)
func (sdlr *SDLRenderer) ChangeFonts(glyphPath, fontPath string) (err error) {
	if sdlr.glyphs != nil {
		sdlr.glyphs.Destroy()
	}
	sdlr.glyphs, err = sdlr.loadTexture(glyphPath)
	if err != nil {
		LogError("SDL Renderer: Could not load font at " + glyphPath)
		return
	}
	if sdlr.font != nil {
		sdlr.font.Destroy()
	}
	sdlr.font, err = sdlr.loadTexture(fontPath)
	if err != nil {
		LogError("SDL Renderer: Could not load font at " + fontPath)
		return
	}
	LogInfo("SDL Renderer: Loaded fonts! Glyph: " + glyphPath + ", Text: " + fontPath)

	_, _, gw, _, _ := sdlr.glyphs.Query()

	//reset window size if fontsize changed
	if int(gw/16) != sdlr.tileSize {
		sdlr.tileSize = int(gw / 16)
		if console == nil {
			LogError("SDL Renderer: Console not initialized, cannot determine screen size.")
			err = errors.New("Console not intialized")
			return
		}
		sdlr.window.SetSize(int32(sdlr.tileSize*console.width), int32(sdlr.tileSize*console.height))
		_ = sdlr.createCanvasBuffer() //TODO: handle this error?
		LogInfo("RENDERER: resized window.")
	}

	return
}

//Loads a bmp font into the GPU using the current window renderer.
//TODO: support more than bmps?
func (sdlr *SDLRenderer) loadTexture(path string) (*sdl.Texture, error) {
	image, err := sdl.LoadBMP(path)
	defer image.Free()
	if err != nil {
		return nil, errors.New("Failed to load image: " + fmt.Sprint(sdl.GetError()))
	}
	image.SetColorKey(true, COL_FUSCHIA)
	texture, err := sdlr.renderer.CreateTextureFromSurface(image)
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

func (sdlr *SDLRenderer) createCanvasBuffer() (err error) {
	if sdlr.canvasBuffer != nil {
		sdlr.canvasBuffer.Destroy()
	}
	sdlr.canvasBuffer, err = sdlr.renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_TARGET, int32(console.width*sdlr.tileSize), int32(console.height*sdlr.tileSize))
	if err != nil {
		LogError("CONSOLE: Failed to create buffer texture. sdl:" + fmt.Sprint(sdl.GetError()))
	}
	return
}

//Enables or disables fullscreen. All burl consoles use borderless fullscreen instead of native
//and the output is scaled to the monitor size.
func (sdlr *SDLRenderer) SetFullscreen(enable bool) {
	if enable {
		sdlr.window.SetFullscreen(sdl.WINDOW_FULLSCREEN_DESKTOP)
		sdlr.window.SetBordered(false)
	} else {
		sdlr.window.SetFullscreen(0)
		sdlr.window.SetBordered(true)
	}
}

//Toggles between fullscreen modes.
func (sdlr *SDLRenderer) ToggleFullscreen() {
	if (sdlr.window.GetFlags() & sdl.WINDOW_FULLSCREEN_DESKTOP) != 0 {
		sdlr.SetFullscreen(false)
	} else {
		sdlr.SetFullscreen(true)
	}
}

//Renders the console to the GPU and flips the buffer.
func (sdlr *SDLRenderer) Render() {
	//render fps counter
	if sdlr.showFPS && sdlr.frames%(30) == 0 {
		fpsString := fmt.Sprintf("%d fps", sdlr.frames*1000/int(sdl.GetTicks()))
		console.DrawText(0, 0, 100, fpsString, COL_WHITE, COL_BLACK, 0)
	}

	w := console.width

	//render the scene!
	var src, dst sdl.Rect
	t := sdlr.renderer.GetRenderTarget()             //store window texture, we'll switch back to it once we're done with the buffer.
	sdlr.renderer.SetRenderTarget(sdlr.canvasBuffer) //point renderer at buffer texture, we'll draw there
	for i, cell := range console.Cells {
		if cell.Dirty || sdlr.forceRedraw {
			if cell.Mode == DRAW_TEXT {
				for c_i, char := range cell.Chars {
					dst = makeRect((i%w)*sdlr.tileSize+c_i*sdlr.tileSize/2, (i/w)*sdlr.tileSize, sdlr.tileSize/2, sdlr.tileSize)
					src = makeRect((char%32)*sdlr.tileSize/2, (char/32)*sdlr.tileSize, sdlr.tileSize/2, sdlr.tileSize)
					sdlr.copyToRenderer(DRAW_TEXT, src, dst, cell.ForeColour, cell.BackColour, char)
				}
			} else {
				if cell.Border {
					console.CalcBorderGlyph(i%w, i/w)
				}
				g := console.Cells[i].Glyph
				dst = makeRect((i%w)*sdlr.tileSize, (i/w)*sdlr.tileSize, sdlr.tileSize, sdlr.tileSize)
				src = makeRect((g%16)*sdlr.tileSize, (g/16)*sdlr.tileSize, sdlr.tileSize, sdlr.tileSize)
				sdlr.copyToRenderer(DRAW_GLYPH, src, dst, cell.ForeColour, cell.BackColour, g)
			}

			console.Cells[i].Dirty = false
		}
	}

	sdlr.renderer.SetRenderTarget(t) //point renderer at window again
	sdlr.renderer.Copy(sdlr.canvasBuffer, nil, nil)
	sdlr.renderer.Present()
	sdlr.renderer.Clear()
	sdlr.forceRedraw = false

	//framerate limiter, so the cpu doesn't implode
	sdlr.elapsed = time.Since(sdlr.frameTime)
	if sdlr.elapsed < sdlr.frameTargetDur {
		time.Sleep(sdlr.frameTargetDur - sdlr.elapsed)
	}
	sdlr.frameTime = time.Now()
	sdlr.frames++
}

//Copies a rect of pixeldata from a source texture to a rect on the renderer's target.
func (sdlr *SDLRenderer) copyToRenderer(mode drawmode, src, dst sdl.Rect, fore, back uint32, g int) {
	//change backcolour if it is different from previous draw
	if back != sdlr.backDrawColour {
		sdlr.backDrawColour = back
		sdlr.renderer.SetDrawColor(GetRGBA(back))
	}

	if sdlr.showChanges {
		sdlr.renderer.SetDrawColor(uint8((sdlr.frames*10)%255), uint8(((sdlr.frames+100)*10)%255), uint8(((sdlr.frames+200)*10)%255), 0xFF) //Test Function
	}

	sdlr.renderer.FillRect(&dst)

	//if we're drawing a nothing character (space, whatever), skip next part.
	if mode == DRAW_GLYPH && (g == GLYPH_NONE || g == GLYPH_SPACE) {
		return
	} else if mode == DRAW_TEXT && g == 32 {
		return
	}

	//change texture color mod if it is different from previous draw, then draw glyph/text
	if mode == DRAW_GLYPH {
		if fore != sdlr.foreDrawColourGlyph {
			sdlr.foreDrawColourGlyph = fore
			sdlr.setTextureColour(sdlr.glyphs, sdlr.foreDrawColourGlyph)
		}
		sdlr.renderer.Copy(sdlr.glyphs, &src, &dst)
	} else {
		if fore != sdlr.foreDrawColourText {
			sdlr.foreDrawColourText = fore
			sdlr.setTextureColour(sdlr.font, sdlr.foreDrawColourText)
		}
		sdlr.renderer.Copy(sdlr.font, &src, &dst)
	}
}

func (sdlr *SDLRenderer) setTextureColour(tex *sdl.Texture, colour uint32) {
	r, g, b, a := GetRGBA(colour)
	tex.SetColorMod(r, g, b)
	tex.SetAlphaMod(a)
}

//Sets maximum framerate as enforced by the framerate limiter. NOTE: cannot go higher than 1000 fps.
func (sdlr *SDLRenderer) SetFramerate(f int) {
	f = Min(f, 1000)
	sdlr.frameTargetDur = time.Duration(1000/float64(f+1)) * time.Millisecond
}

func (sdlr *SDLRenderer) ForceRedraw() {
	sdlr.forceRedraw = true
}

func (sdlr *SDLRenderer) ToggleDebugMode(m string) {
	switch m {
	case "fps":
		sdlr.showFPS = !sdlr.showFPS
	case "changes":
		sdlr.showChanges = !sdlr.showChanges
	default:
		LogError("SDL Renderer: no debug mode called " + m)
	}
}

//int32 for rect arguments. what a world.
func makeRect(x, y, w, h int) sdl.Rect {
	return sdl.Rect{X: int32(x), Y: int32(y), W: int32(w), H: int32(h)}
}
