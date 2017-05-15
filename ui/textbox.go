package ui

import "strings"
import "github.com/bennicholls/burl/console"

//UI Element for displaying text.
type Textbox struct {
	UIElement
	text          string
	centered      bool
}

func NewTextbox(w, h, x, y, z int, bord, cent bool, txt string) *Textbox {
	return &Textbox{NewUIElement(x,y,z,w,h,bord), txt, cent}
}

//Returns the height required to fit a string after it has been wrapped. Reimplements the word wrapper but cruder.
func CalcWrapHeight(s string, width int) int {
	width = width * 2
	line := ""
	n := 0
	for _, word := range strings.Split(s, " ") {
		//super long word make-it-not-break hack.
		if len(word) > width {
			continue
		}

		if len(line)+len(word) > width {
			n++
			line = ""
		}
		line += word
		if len(line) != width {
			line += " "
		}
	}

	return n + 1
}

//Replaces the string for the textbox.
func (t *Textbox) ChangeText(txt string) {
	if t.text != txt {
		t.text = txt
	}
}

//Render function optionally takes an offset (for containering), 2 or 3 ints.
func (t *Textbox) Render(offset ...int) {
	if t.visible {
		offX, offY, offZ := processOffset(offset)

		//word wrap calculatrix. a mighty sinful thing.
		//TODO: support for breaking super long words. right now it just skips the word.
		lines := make([]string, t.height)
		n := 0
		for _, s := range strings.Split(t.text, " ") {
			//super long word make-it-not-break hack.
			if len(s) > t.width*2 {
				continue
			}

			if len(lines[n])+len(s) > t.width*2 {
				//make sure we don't overflow the textbox
				if n < len(lines)-1 {
					n++
				} else {
					break
				}
			}
			lines[n] += s
			if len(lines[n]) != t.width*2 {
				lines[n] += " "
			}
		}

		for l := 0; l < len(lines); l++ {
			offX := offX //so we can modify the offset separately for each line

			//clear texbox (fill with spaces).
			for i := 0; i < t.width*t.height; i++ {
				console.ChangeText(offX+t.x+i%t.width, offY+t.y+l, offZ+t.z, int(' '), int(' '))
			}

			//offset if centered
			if t.centered {
				offX += (t.width/2 - len(lines[l])/4)
			}

			//print text
			console.DrawText(offX+t.x, offY+t.y+l, offZ+t.z, lines[l], 0xFFFFFFFF, 0xFF000000)
		}

		t.UIElement.Render(offX, offY, offZ)
	}
}

func (t *Textbox) SetCentered(c bool) {
	t.centered = c
}

