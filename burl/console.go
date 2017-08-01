package burl

import "github.com/veandco/go-sdl2/sdl"
import "fmt"
import "errors"

type Console struct {
	window       *sdl.Window
	renderer     *sdl.Renderer
	glyphs       *sdl.Texture
	font         *sdl.Texture
	canvasBuffer *sdl.Texture
	format       *sdl.PixelFormat

	width, height, tileSize int

	canvas                []Cell
	forceRedraw           bool
	frameTime, ticks, fps uint32
	frames                int
	showFPS               bool
	showChanges           bool
	Ready                 bool //true when console is ready for drawing and stuff!

	//store render colours so we don't have to set them for every renderer.Copy()
	backDrawColour      uint32
	foreDrawColourText  uint32
	foreDrawColourGlyph uint32
}

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

	//for text rendering mode. TODO:multiple back and fore colours, one for each char
	Mode  drawmode
	Chars [2]int
}

//Sets the properties of a cell all at once for Glyph Mode.
func (c *Cell) SetGlyph(gl int, fore, back uint32, z int) {
	if c.Glyph != gl || c.ForeColour != fore || c.BackColour != back || c.Z != z || c.Mode == DRAW_TEXT {
		c.Mode = DRAW_GLYPH
		c.Glyph = gl
		c.ForeColour = fore
		c.BackColour = back
		c.Z = z
		c.Dirty = true
	}
}

//Sets the properties of a cell all at once for Text Mode.
func (c *Cell) SetText(char1, char2 int, fore, back uint32, z int) {
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

//Re-inits a cell back to default blankness.
func (c *Cell) Clear() {
	if c.Mode == DRAW_TEXT {
		c.SetText(32, 32, 0xFF000000, 0xFF000000, 0)
	} else {
		c.SetGlyph(GLYPH_NONE, 0xFF000000, 0xFF000000, 0)
	}
}

//Setup the game window, renderer, etc
func (c *Console) Setup(w, h int, glyphPath, fontPath, title string) (err error) {
	c.width = w
	c.height = h
	c.tileSize = 24

	c.window, err = sdl.CreateWindow(title, sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, c.width*c.tileSize, c.height*c.tileSize, sdl.WINDOW_OPENGL)
	if err != nil {
		LogError("CONSOLE: Failed to create window. sdl:" + fmt.Sprint(sdl.GetError()))
		return errors.New("Failed to create window.")
	}

	//manually set pixelformat to ARGB (window defaults to RGB for some reason)
	c.format, err = sdl.AllocFormat(uint(sdl.PIXELFORMAT_ARGB8888))
	if err != nil {
		LogError("CONSOLE: Failed to allocate pixelformat. sdl:" + fmt.Sprint(sdl.GetError()))
		return errors.New("No pixelformat.")
	}

	c.renderer, err = sdl.CreateRenderer(c.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		LogError("CONSOLE: Failed to create renderer. sdl:" + fmt.Sprint(sdl.GetError()))
		return errors.New("Failed to create renderer.")
	}
	c.renderer.Clear()

	err = c.CreateCanvasBuffer()
	if err != nil {
		return errors.New("Failed to create canvas buffer.")
	}

	c.canvas = make([]Cell, c.width*c.height)
	c.Clear()

	//init drawing fonts
	err = c.ChangeFonts(glyphPath, fontPath)
	if err != nil {
		return errors.New("Could not load fonts.")
	}

	c.frames = 0
	c.frameTime, c.ticks = 0, 0
	c.fps = 17 //17ms = 60 FPS approx
	c.showFPS = false
	c.showChanges = false
	c.Ready = true

	return nil
}

//Enables fullscreen.
//TODO: the opposite??? do this later when resolution/window mode polish goes in.
func (c *Console) SetFullscreen() {
	c.window.SetFullscreen(sdl.WINDOW_FULLSCREEN)
}

//Loads new fonts to the renderer and changes the tilesize (and by entension, the window size)
func (c *Console) ChangeFonts(glyphPath, fontPath string) (err error) {
	if c.glyphs != nil {
		c.glyphs.Destroy()
	}
	c.glyphs, err = c.LoadTexture(glyphPath)
	if err != nil {
		LogError("CONSOLE: Could not load font at " + glyphPath)
		return
	}
	if c.font != nil {
		c.font.Destroy()
	}
	c.font, err = c.LoadTexture(fontPath)
	if err != nil {
		LogError("CONSOLE: Could not load font at " + fontPath)
		return
	}
	c.Clear()
	LogInfo("CONSOLE: Loaded fonts! Glyph: " + glyphPath + ", Text:" + fontPath)

	_, _, gw, _, _ := c.glyphs.Query()

	//reset window size if fontsize changed
	if int(gw/16) != c.tileSize {
		c.tileSize = int(gw / 16)
		c.window.SetSize(c.tileSize*c.width, c.tileSize*c.height)
		_ = c.CreateCanvasBuffer() //TODO: handle this error?
		LogInfo("CONSOLE: resized window.")
	}

	return
}

func (c *Console) CreateCanvasBuffer() (err error) {
	if c.canvasBuffer != nil {
		c.canvasBuffer.Destroy()
	}
	c.canvasBuffer, err = c.renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_TARGET, c.width*c.tileSize, c.height*c.tileSize)
	if err != nil {
		LogError("CONSOLE: Failed to create buffer texture. sdl:" + fmt.Sprint(sdl.GetError()))
	}
	return
}

//Loads a bmp font into the GPU using the current window renderer.
//TODO: support more than bmps?
func (c *Console) LoadTexture(path string) (*sdl.Texture, error) {
	image, err := sdl.LoadBMP(path)
	defer image.Free()
	if err != nil {
		return nil, errors.New("Failed to load image: " + fmt.Sprint(sdl.GetError()))
	}
	image.SetColorKey(1, 0xFFFF00FF)
	texture, err := c.renderer.CreateTextureFromSurface(image)
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
func (c *Console) Render() {
	//render fps counter
	if c.showFPS && c.frames%(30) == 0 {
		fpsString := fmt.Sprintf("%d fps", c.frames*1000/int(sdl.GetTicks()))
		c.DrawText(0, 0, 10, fpsString, 0xFFFFFFFF, 0xFF000000)
	}

	//render the scene!
	var src, dst sdl.Rect
	t := c.renderer.GetRenderTarget()          //store window texture, we'll switch back to it once we're done with the buffer.
	c.renderer.SetRenderTarget(c.canvasBuffer) //point renderer at buffer texture, we'll draw there
	for i, cell := range c.canvas {
		if cell.Dirty || c.forceRedraw {
			if cell.Mode == DRAW_TEXT {
				for c_i, char := range cell.Chars {
					dst = makeRect((i%c.width)*c.tileSize+c_i*c.tileSize/2, (i/c.width)*c.tileSize, c.tileSize/2, c.tileSize)
					src = makeRect((char%32)*c.tileSize/2, (char/32)*c.tileSize, c.tileSize/2, c.tileSize)
					c.CopyToRenderer(DRAW_TEXT, src, dst, cell.ForeColour, cell.BackColour, char)
				}
			} else {
				dst = makeRect((i%c.width)*c.tileSize, (i/c.width)*c.tileSize, c.tileSize, c.tileSize)
				src = makeRect((cell.Glyph%16)*c.tileSize, (cell.Glyph/16)*c.tileSize, c.tileSize, c.tileSize)
				c.CopyToRenderer(DRAW_GLYPH, src, dst, cell.ForeColour, cell.BackColour, cell.Glyph)
			}

			c.canvas[i].Dirty = false
		}
	}

	c.renderer.SetRenderTarget(t) //point renderer at window again
	r := makeRect(0, 0, c.width*c.tileSize, c.height*c.tileSize)
	c.renderer.Copy(c.canvasBuffer, &r, &r)
	c.renderer.Present()
	c.renderer.Clear()
	c.forceRedraw = false

	//framerate limiter, so the cpu doesn't implode
	c.ticks = sdl.GetTicks() - c.frameTime
	if c.ticks < c.fps {
		sdl.Delay(c.fps - c.ticks)
	}
	c.frameTime = sdl.GetTicks()
	c.frames++
}

//Copies a rect of pixeldata from a source texture to a rect on the renderer's target.
func (c *Console) CopyToRenderer(mode drawmode, src, dst sdl.Rect, fore, back uint32, g int) {
	//change backcolour if it is different from previous draw
	if back != c.backDrawColour {
		c.backDrawColour = back
		c.renderer.SetDrawColor(sdl.GetRGBA(back, c.format))
	}

	if c.showChanges {
		c.renderer.SetDrawColor(sdl.GetRGBA(c.MakeColour((c.frames*10)%255, ((c.frames+100)*10)%255, ((c.frames+200)*10)%255), c.format)) //Test Function
	}

	c.renderer.FillRect(&dst)

	//if we're drawing a nothing character (space, whatever), skip next part.
	if mode == DRAW_GLYPH && (g == GLYPH_NONE || g == GLYPH_SPACE) {
		return
	} else if mode == DRAW_TEXT && g == 32 {
		return
	}

	//change texture color mod if it is different from previous draw, then draw glyph/text
	if mode == DRAW_GLYPH {
		if fore != c.foreDrawColourGlyph {
			c.foreDrawColourGlyph = fore
			c.SetTextureColour(c.glyphs, c.foreDrawColourGlyph)
		}
		c.renderer.Copy(c.glyphs, &src, &dst)
	} else {
		if fore != c.foreDrawColourText {
			c.foreDrawColourText = fore
			c.SetTextureColour(c.font, c.foreDrawColourText)
		}
		c.renderer.Copy(c.font, &src, &dst)
	}
}

func (c *Console) SetTextureColour(tex *sdl.Texture, colour uint32) {
	r, g, b, a := sdl.GetRGBA(colour, c.format)
	tex.SetColorMod(r, g, b)
	tex.SetAlphaMod(a)
}

//Sets maximum framerate as enforced by the framerate limiter. NOTE: cannot go higher than 1000 fps.
func (c *Console) SetFramerate(f int) {
	c.fps = uint32(1000/f) + 1
}

//Toggles rendering of the FPS meter.
func (c *Console) ToggleFPS() {
	c.showFPS = !c.showFPS
}

func (c *Console) ToggleChanges() {
	c.showChanges = !c.showChanges
}

func (c *Console) ForceRedraw() {
	c.forceRedraw = true
}

//Deletes special graphics structures, closes files, etc. Defer this function!
func (c *Console) Cleanup() {
	c.format.Free()
	c.glyphs.Destroy()
	c.font.Destroy()
	c.canvasBuffer.Destroy()
	c.renderer.Destroy()
	c.window.Destroy()
}

//Changes the glyph of a cell in the canvas at position (x, y).
func (c *Console) ChangeGlyph(x, y, glyph int) {
	s := y*c.width + x
	if CheckBounds(x, y, c.width, c.height) {
		c.canvas[s].SetGlyph(glyph, c.canvas[s].ForeColour, c.canvas[s].BackColour, c.canvas[s].Z)
	}
}

//Changes text of a cell in the canvas at position (x, y).
func (c *Console) ChangeText(x, y, z, char1, char2 int) {
	s := y*c.width + x
	if CheckBounds(x, y, c.width, c.height) && c.canvas[s].Z <= z {
		c.canvas[s].Mode = DRAW_TEXT
		if c.canvas[s].Chars[0] != char1 || c.canvas[s].Chars[1] != char2 {
			c.canvas[s].Chars[0] = char1
			c.canvas[s].Chars[1] = char2
			c.canvas[s].Z = z
			c.canvas[s].Dirty = true
		}
	}
}

//Changes a single character on the canvas at position (x,y) in text mode.
//charNum: 0 = Left, 1 = Right (for ease with modulo operations). Throw whatever in here though, it gets modulo 2'd anyways just in case.
func (c *Console) ChangeChar(x, y, z, char, charNum int) {
	s := y*c.width + x
	if CheckBounds(x, y, c.width, c.height) && charNum >= 0 && c.canvas[s].Z <= z {
		c.canvas[s].Mode = DRAW_TEXT
		if c.canvas[s].Chars[charNum%2] != char {
			c.canvas[s].Chars[charNum%2] = char
			c.canvas[s].Z = z
			c.canvas[s].Dirty = true
		}
	}
}

//Changes the foreground drawing colour of a cell in the canvas at position (x, y).
func (c *Console) ChangeForeColour(x, y, z int, fore uint32) {
	s := y*c.width + x
	if CheckBounds(x, y, c.width, c.height) && c.canvas[s].Z <= z {
		if c.canvas[s].Mode == DRAW_TEXT {
			c.canvas[s].SetText(c.canvas[s].Chars[0], c.canvas[s].Chars[1], fore, c.canvas[s].BackColour, z)
		} else {
			c.canvas[s].SetGlyph(c.canvas[s].Glyph, fore, c.canvas[s].BackColour, z)
		}
	}
}

//Changes the background colour of a cell in the canvas at position (x, y).
func (c *Console) ChangeBackColour(x, y, z int, back uint32) {
	s := y*c.width + x
	if CheckBounds(x, y, c.width, c.height) && c.canvas[s].Z <= z {
		if c.canvas[s].Mode == DRAW_TEXT {
			c.canvas[s].SetText(c.canvas[s].Chars[0], c.canvas[s].Chars[1], c.canvas[s].ForeColour, back, z)
		} else {
			c.canvas[s].SetGlyph(c.canvas[s].Glyph, c.canvas[s].ForeColour, back, z)
		}
	}
}

func (c *Console) ChangeColours(x, y, z int, fore, back uint32) {
	s := y*c.width + x
	if CheckBounds(x, y, c.width, c.height) && c.canvas[s].Z <= z {
		if c.canvas[s].Mode == DRAW_TEXT {
			c.canvas[s].SetText(c.canvas[s].Chars[0], c.canvas[s].Chars[1], fore, back, z)
		} else {
			c.canvas[s].SetGlyph(c.canvas[s].Glyph, fore, back, z)
		}
	}
}

//Simultaneously changes all characteristics of a glyph cell in the canvas at position (x, y).
//TODO: change name of this to signify it is for changing glyph cells.
func (c *Console) ChangeCell(x, y, z, glyph int, fore, back uint32) {
	s := y*c.width + x
	if CheckBounds(x, y, c.width, c.height) && c.canvas[s].Z <= z {
		c.canvas[s].SetGlyph(glyph, fore, back, z)
	}
}

//Draws a string to the console in text mode.
func (c *Console) DrawText(x, y, z int, txt string, fore, back uint32) {
	for i, char := range txt {
		if CheckBounds(x+i/2, y, c.width, c.height) {
			c.ChangeChar(x+i/2, y, z, int(char), i%2)
			if i%2 == 0 {
				//only need to change colour each cell, not each character
				c.ChangeForeColour(x+i/2, y, z, fore)
				c.ChangeBackColour(x+i/2, y, z, back)
				if i == len(txt)-1 {
					//if final character is in the left-side of a cell, blank the right side.
					c.ChangeChar(x+i/2, y, z, 32, 1)
				}
			}
		}
	}
}

//TODO: custom colouring, multiple styles.
//NOTE: current border colouring thing is a bit of a hack. Need to add actual support for
//border and ui styling.
func (c *Console) DrawBorder(x, y, z, w, h int, title string, focused bool) {
	//set border colour.
	bc := uint32(0xFFE28F00)
	if !focused {
		bc = 0xFF555555
	}
	//Top and bottom.
	for i := 0; i < w; i++ {
		c.ChangeCell(x+i, y-1, z, GLYPH_BORDER_LR, bc, 0xFF000000)
		c.ChangeCell(x+i, y+h, z, GLYPH_BORDER_LR, bc, 0xFF000000)
	}
	//Sides
	for i := 0; i < h; i++ {
		c.ChangeCell(x-1, y+i, z, GLYPH_BORDER_UD, bc, 0xFF000000)
		c.ChangeCell(x+w, y+i, z, GLYPH_BORDER_UD, bc, 0xFF000000)
	}
	//corners
	c.ChangeCell(x-1, y-1, z, GLYPH_BORDER_DR, bc, 0xFF000000)
	c.ChangeCell(x-1, y+h, z, GLYPH_BORDER_UR, bc, 0xFF000000)
	c.ChangeCell(x+w, y+h, z, GLYPH_BORDER_UL, bc, 0xFF000000)
	c.ChangeCell(x+w, y-1, z, GLYPH_BORDER_DL, bc, 0xFF000000)

	//Write centered title.
	if len(title) < w && title != "" {
		c.DrawText(x+(w/2-len(title)/4-1), y-1, z+1, title, 0xFFFFFFFF, 0xFF000000)
	}
}

//Clears an area of the canvas. Optionally takes a rect (defined by 4 ints) so you can clear specific areas of the console
func (c *Console) Clear(rect ...int) {
	offX, offY, w, h := 0, 0, c.width, c.height

	if len(rect) == 4 {
		offX, offY, w, h = rect[0], rect[1], rect[2], rect[3]
	}

	for i := 0; i < w*h; i++ {
		x := offX + i%w
		y := offY + i/w
		if CheckBounds(x, y, c.width, c.height) {
			c.canvas[y*c.width+x].Clear()
		}
	}
}

//Returns the dimensions of the canvas.
func (c Console) Dims() (w, h int) {
	return c.width, c.height
}

//Takes r,g,b ints and creates a colour as defined by the pixelformat with alpha 255.
//TODO: rgba version of this function? variatic function that can optionally take an alpha? Hmm.
func (c Console) MakeColour(r, g, b int) uint32 {
	return sdl.MapRGBA(c.format, uint8(r), uint8(g), uint8(b), 255)
}

//Changes alpha of a colour.
func (c Console) ChangeColourAlpha(colour uint32, alpha uint8) uint32 {
	r, g, b := sdl.GetRGB(colour, c.format)
	return sdl.MapRGBA(c.format, r, g, b, alpha)
}

//int32 for rect arguments. what a world.
func makeRect(x, y, w, h int) sdl.Rect {
	return sdl.Rect{X: int32(x), Y: int32(y), W: int32(w), H: int32(h)}
}
