package burl

import "github.com/veandco/go-sdl2/sdl"

//Takes r,g,b ints and creates a colour as defined by the pixelformat with alpha 255.
//TODO: rgba version of this function? variatic function that can optionally take an alpha? Hmm.
func MakeColour(r, g, b int) uint32 {
	return sdl.MapRGBA(console.format, uint8(r), uint8(g), uint8(b), 0xFF)
}

//Changes alpha of a colour.
func ChangeAlpha(colour uint32, alpha uint8) uint32 {
	r, g, b := sdl.GetRGB(colour, console.format)
	return sdl.MapRGBA(console.format, r, g, b, alpha)
}

const (
	COL_WHITE      = 0xFFFFFFFF
	COL_BLACK      = 0xFF000000
	COL_RED        = 0xFFFF0000
	COL_BLUE       = 0xFF0000FF
	COL_LIME       = 0xFF00FF00
	COL_LIGHTGREY  = 0xFF444444
	COL_GREY       = 0xFF888888
	COL_DARKGREY   = 0xFFCCCCCC
	COL_YELLOW     = 0xFFFFFF00
	COL_FUSCHIA    = 0xFFFF00FF
	COL_CYAN       = 0xFF00FFFF
	COL_MAROON     = 0xFF800000
	COL_OLIVE      = 0xFF808000
	COL_GREEN      = 0xFF008000
	COL_TEAL       = 0xFF008080
	COL_NAVY       = 0xFF000080
	COL_PURPLE     = 0xFF800080
)