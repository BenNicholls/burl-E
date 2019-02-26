package burl

//Takes r,g,b ints and creates a colour with alpha 255 in ARGB format.
func MakeColour(r, g, b, a int) (colour uint32) {
	colour = uint32((a % 256) << 24)
	colour |= uint32(r%256) << 16
	colour |= uint32(g%256) << 8
	colour |= uint32(b % 256)

	return
}

func MakeOpaqueColour(r, g, b int) uint32 {
	return MakeColour(r, g, b, 255)
}

//Returns the RGBA components of an ARGB8888 formatted uint32 colour.
func GetRGBA(colour uint32) (r, g, b, a uint8) {
	b = uint8(colour & 0x000000FF)
	g = uint8((colour >> 8) & 0x000000FF)
	r = uint8((colour >> 16) & 0x000000FF)
	a = uint8(colour >> 24)

	return
}

//Blends 2 colours c1 and c2. c1 is the active colour (it's on "top").
func BlendColours(c1, c2 uint32, mode BlendMode) uint32 {
	r1, g1, b1, a1 := GetRGBA(c1)
	r2, g2, b2, a2 := GetRGBA(c2)

	var r, g, b, a int

	switch mode {
	case BLEND_MULTIPLY:
		r = int(r1)*int(r2)/255
		g = int(g1)*int(g2)/255
		b = int(b1)*int(b2)/255
		a = int(a1)*int(a2)/255
	case BLEND_SCREEN:
		r = 255-int(255-r1)*int(255-r2)/255
		g = 255-int(255-g1)*int(255-g2)/255
		b = 255-int(255-b1)*int(255-b2)/255
		a = 255-int(255-a1)*int(255-a2)/255
	}

	return MakeColour(r, g, b, a)
}

type BlendMode int 

const (
	BLEND_MULTIPLY BlendMode = iota
	BLEND_SCREEN
)

const (
	COL_NONE      uint32 = 0x00000000
	COL_WHITE     uint32 = 0xFFFFFFFF
	COL_BLACK     uint32 = 0xFF000000
	COL_RED       uint32 = 0xFFFF0000
	COL_BLUE      uint32 = 0xFF0000FF
	COL_LIME      uint32 = 0xFF00FF00
	COL_LIGHTGREY uint32 = 0xFF444444
	COL_GREY      uint32 = 0xFF888888
	COL_DARKGREY  uint32 = 0xFFCCCCCC
	COL_YELLOW    uint32 = 0xFFFFFF00
	COL_FUSCHIA   uint32 = 0xFFFF00FF
	COL_CYAN      uint32 = 0xFF00FFFF
	COL_MAROON    uint32 = 0xFF800000
	COL_OLIVE     uint32 = 0xFF808000
	COL_GREEN     uint32 = 0xFF008000
	COL_TEAL      uint32 = 0xFF008080
	COL_NAVY      uint32 = 0xFF000080
	COL_PURPLE    uint32 = 0xFF800080
)

type Palette []uint32 

//Generate a palette with num items, passing from colour c1 to c2. The colours are
//lineraly interpolated evenly from one to the next. Palette is NOT circular.
//TODO: Circular palette function?
func GeneratePalette(num int, c1, c2 uint32) (p Palette) {
	p = make(Palette, num)

	r1, g1, b1, _ := GetRGBA(c1)
	r2, g2, b2, _ := GetRGBA(c2)

	for i := range p {
		p[i] = MakeOpaqueColour(Lerp(int(r1), int(r2), i, len(p)), Lerp(int(g1), int(g2), i, len(p)), Lerp(int(b1), int(b2), i, len(p)))
	}

	p[num-1] = c2 //fix end of palette rounding lerp stuff.

	return
}

//Adds the palette p2 to the end of p.
func (p *Palette) Add(p2 Palette) {
	if (*p)[len(*p) - 1] == p2[0] {
		*p = append(*p, p2[1:]...)
	} else {
		*p = append(*p, p2...)
	}
}