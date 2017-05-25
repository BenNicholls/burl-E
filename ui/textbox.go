package ui

import "github.com/bennicholls/burl/console"
import "github.com/bennicholls/burl/util"

//UI Element for displaying text.
type Textbox struct {
	UIElement
	text     string
	centered bool
	lines    []string
}

func NewTextbox(w, h, x, y, z int, bord, cent bool, txt string) *Textbox {
	return &Textbox{NewUIElement(x, y, z, w, h, bord), txt, cent, util.WrapText(txt, w*2, h)}
}

//Returns the height required to fit a string after it has been wrapped.
func CalcWrapHeight(s string, width int) int {
	return len(util.WrapText(s, width*2))
}

//Replaces the string for the textbox.
func (t *Textbox) ChangeText(txt string) {
	if t.text != txt {
		t.text = txt
		t.lines = util.WrapText(txt, t.width*2, t.height)
	}
}

//Render function optionally takes an offset (for containering), 2 or 3 ints.
func (t *Textbox) Render(offset ...int) {
	if t.visible {
		offX, offY, offZ := processOffset(offset)

		for l, line := range t.lines {
			offX := offX //so we can modify the offset separately for each line

			//clear texbox (fill with spaces).
			for i := 0; i < t.width*t.height; i++ {
				console.ChangeText(offX+t.x+i%t.width, offY+t.y+l, offZ+t.z, int(' '), int(' '))
			}

			//offset if centered
			if t.centered {
				offX += (t.width/2 - len(line)/4)
			}

			//print text
			console.DrawText(offX+t.x, offY+t.y+l, offZ+t.z, line, 0xFFFFFFFF, 0xFF000000)
		}

		t.UIElement.Render(offX, offY, offZ)
	}
}

func (t *Textbox) SetCentered(c bool) {
	t.centered = c
}
