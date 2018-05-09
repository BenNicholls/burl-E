package burl

import "github.com/veandco/go-sdl2/sdl"
import "fmt"
import "errors"
import "time"

type Console struct {
	window       *sdl.Window
	renderer     *sdl.Renderer
	glyphs       *sdl.Texture
	font         *sdl.Texture
	canvasBuffer *sdl.Texture

	width, height, tileSize int

	canvas       []Cell
	forceRedraw  bool
	frameTime    time.Time
	fps, elapsed time.Duration
	frames       int
	showFPS      bool
	showChanges  bool
	Ready        bool //true when console is ready for drawing and stuff!

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
	Border     bool //marks cell as part of a UI Element border.

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
		if gl < GLYPH_BORDER_UD || gl > GLYPH_BORDER_DR {
			c.Border = false
		}
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
		c.SetText(32, 32, COL_BLACK, COL_BLACK, 0)
	} else {
		c.SetGlyph(GLYPH_NONE, COL_BLACK, COL_BLACK, 0)
	}
}

//Setup the game window, renderer, etc
func (c *Console) Setup(w, h int, glyphPath, fontPath, title string) (err error) {
	c.width = w
	c.height = h
	c.tileSize = 24

	c.window, err = sdl.CreateWindow(title, sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, int32(c.width*c.tileSize), int32(c.height*c.tileSize), sdl.WINDOW_OPENGL)
	if err != nil {
		LogError("CONSOLE: Failed to create window. sdl:" + fmt.Sprint(sdl.GetError()))
		return errors.New("Failed to create window.")
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
	c.SetFramerate(60)
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
		c.window.SetSize(int32(c.tileSize*c.width), int32(c.tileSize*c.height))
		_ = c.CreateCanvasBuffer() //TODO: handle this error?
		LogInfo("CONSOLE: resized window.")
	}

	return
}

func (c *Console) CreateCanvasBuffer() (err error) {
	if c.canvasBuffer != nil {
		c.canvasBuffer.Destroy()
	}
	c.canvasBuffer, err = c.renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_TARGET, int32(c.width*c.tileSize), int32(c.height*c.tileSize))
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
	image.SetColorKey(true, COL_FUSCHIA)
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
		c.DrawText(0, 0, 100, fpsString, COL_WHITE, COL_BLACK, 0)
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
				if cell.Border {
					c.CalcBorderGlyph(i%c.width, i/c.width)
				}
				dst = makeRect((i%c.width)*c.tileSize, (i/c.width)*c.tileSize, c.tileSize, c.tileSize)
				src = makeRect((c.canvas[i].Glyph%16)*c.tileSize, (c.canvas[i].Glyph/16)*c.tileSize, c.tileSize, c.tileSize)
				c.CopyToRenderer(DRAW_GLYPH, src, dst, cell.ForeColour, cell.BackColour, c.canvas[i].Glyph)
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
	c.elapsed = time.Since(c.frameTime)
	if c.elapsed < c.fps {
		time.Sleep(c.fps - c.elapsed)
	}
	c.frameTime = time.Now()
	c.frames++
}

//Copies a rect of pixeldata from a source texture to a rect on the renderer's target.
func (c *Console) CopyToRenderer(mode drawmode, src, dst sdl.Rect, fore, back uint32, g int) {
	//change backcolour if it is different from previous draw
	if back != c.backDrawColour {
		c.backDrawColour = back
		c.renderer.SetDrawColor(GetRGBA(back))
	}

	if c.showChanges {
		c.renderer.SetDrawColor(uint8((c.frames*10)%255), uint8(((c.frames+100)*10)%255), uint8(((c.frames+200)*10)%255), 0xFF) //Test Function
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
	r, g, b, a := GetRGBA(colour)
	tex.SetColorMod(r, g, b)
	tex.SetAlphaMod(a)
}

//Sets maximum framerate as enforced by the framerate limiter. NOTE: cannot go higher than 1000 fps.
func (c *Console) SetFramerate(f int) {
	c.fps = time.Duration(1000/float64(f+1)) * time.Millisecond
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
	c.glyphs.Destroy()
	c.font.Destroy()
	c.canvasBuffer.Destroy()
	c.renderer.Destroy()
	c.window.Destroy()
}

//Returns a reference to the cell at (x, y). Returns nil if (x, y) is bad.
func (c *Console) GetCell(x, y int) *Cell {
	if CheckBounds(x, y, c.width, c.height) {
		return &c.canvas[y*c.width+x]
	}
	return nil
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

//Draws a string to the console in text mode. CharNum determines which half of the cell we
//start in. See ChageChar() for details.
func (c *Console) DrawText(x, y, z int, txt string, fore, back uint32, charNum int) {
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

//TODO: custom colouring, multiple styles.
//NOTE: current border colouring thing is a bit of a hack. Need to add actual support for
//border and ui styling.
//Borders work by setting a flag on the cells that need to be borders. At render time, any
//borders with a dirty flag are assigned a border glyph based on the state of their neighbours:
//if the neighnouring cells are on the same z level and also borders, they will connect.
func (c *Console) DrawBorder(x, y, z, w, h int, title, hint string, focused bool) {
	//set border colour.
	bc := COL_PURPLE
	if !focused {
		bc = COL_LIGHTGREY
	}
	//Top and bottom.
	for i := -1; i <= w; i++ {
		c.SetCellBorder(x+i, y-1, z, bc)
		c.SetCellBorder(x+i, y+h, z, bc)
	}
	//Sides
	for i := 0; i < h; i++ {
		c.SetCellBorder(x-1, y+i, z, bc)
		c.SetCellBorder(x+w, y+i, z, bc)
	}

	//Write centered title.
	if len(title) < w && title != "" {
		c.DrawText(x+(w/2-len(title)/4-1), y-1, z+1, title, COL_WHITE, COL_BLACK, 0)
	}

	//Write right-justified hint text
	if len(hint) < 2*w && hint != "" {
		decoratedHint := TEXT_BORDER_DECO_LEFT + hint + TEXT_BORDER_DECO_RIGHT
		offset := w - len(hint)/2 - 1
		if len(hint)%2 == 1 {
			decoratedHint = TEXT_BORDER_LR + decoratedHint
			offset -= 1
		}

		c.DrawText(x+offset, y+h, z, decoratedHint, bc, COL_BLACK, 0)
	}
}

//Sets the cell at (x, y) as a border cell, to be drawn out later when the frame is rendered.
func (c *Console) SetCellBorder(x, y, z int, fore uint32) {
	s := y*c.width + x
	if CheckBounds(x, y, c.width, c.height) && c.canvas[s].Z <= z {
		if !c.canvas[s].Border || c.canvas[s].Z < z || c.canvas[s].ForeColour != fore {
			c.canvas[s].Border = true
			c.canvas[s].Z = z
			c.canvas[s].ForeColour = fore
			c.canvas[s].Mode = DRAW_GLYPH
			c.canvas[s].Dirty = true
		}
	}
}

//Chooses a border glyph for a cell at (x,y) based on the border state of it's neighbours.
func (c *Console) CalcBorderGlyph(x, y int) {
	if !CheckBounds(x, y, c.width, c.height) {
		return
	}

	s := y*c.width + x
	var g int
	var u, d, l, r bool

	if uCell := c.GetCell(x, y-1); uCell != nil && uCell.Z == c.canvas[s].Z {
		u = uCell.Border
	}
	if dCell := c.GetCell(x, y+1); dCell != nil && dCell.Z == c.canvas[s].Z {
		d = dCell.Border
	}
	if lCell := c.GetCell(x-1, y); lCell != nil && lCell.Z == c.canvas[s].Z {
		l = lCell.Border
	}
	if rCell := c.GetCell(x+1, y); rCell != nil && rCell.Z == c.canvas[s].Z {
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
		c.canvas[s].Glyph = g
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

//int32 for rect arguments. what a world.
func makeRect(x, y, w, h int) sdl.Rect {
	return sdl.Rect{X: int32(x), Y: int32(y), W: int32(w), H: int32(h)}
}
