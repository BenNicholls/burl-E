package burl

//UI Element for displaying text.
type Textbox struct {
	UIElement
	text     string
	centered bool
	lines    []string
}

func NewTextbox(w, h, x, y, z int, bord, cent bool, txt string) *Textbox {
	return &Textbox{NewUIElement(x, y, z, w, h, bord), txt, cent, WrapText(txt, w*2, h)}
}

//Returns the height required to fit a string after it has been wrapped.
func CalcWrapHeight(s string, width int) int {
	return len(WrapText(s, width*2))
}

//Replaces the string for the textbox.
func (t *Textbox) ChangeText(txt string) {
	if t.text != txt {
		t.text = txt
		t.lines = WrapText(txt, t.width*2, t.height)
	}
}

//Render function optionally takes an offset (for containering), 2 or 3 ints.
func (t *Textbox) Render(offset ...int) {
	if t.visible {
		offX, offY, offZ := processOffset(offset)

		for l, line := range t.lines {
			lineOffset := 0

			//offset if centered
			if t.centered {
				lineOffset += (t.width/2 - len(line)/4)
			}

			//draw leading spaces
			for i := 0; i < lineOffset; i++ {
				console.ChangeText(offX+t.x+i%t.width, offY+t.y+i/t.width+l, offZ+t.z, int(' '), int(' '))
			}

			//print text
			console.DrawText(offX+t.x+lineOffset, offY+t.y+l, offZ+t.z, line, 0xFFFFFFFF, 0xFF000000)

			//draw trailing spaces
			for i := lineOffset + len(line)/2; i < t.width; i++ {
				console.ChangeText(offX+t.x+i%t.width, offY+t.y+i/t.width+l, offZ+t.z, int(' '), int(' '))
			}
		}

		t.UIElement.Render(offX, offY, offZ)
	}
}

func (t *Textbox) SetCentered(c bool) {
	t.centered = c
}
