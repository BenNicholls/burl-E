package ui

import "strings"
import "github.com/bennicholls/delvetown/console"

//UI Element for displaying text.
type Textbox struct {
	width, height int
	x, y, z       int
	bordered      bool
	centered      bool
	title         string
	text          string
	visible       bool
	anims         []Animator
	focused       bool
}

func NewTextbox(w, h, x, y, z int, bord, cent bool, txt string) *Textbox {
	return &Textbox{w, h, x, y, z, bord, cent, "", txt, true, make([]Animator, 0, 20), false}
}

//Returns the height required to fit a string after it has been wrapped. Reimplements the word wrapper but cruder.
func CalcWrapHeight(s string, width int) int {
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

func (t *Textbox) SetTitle(s string) {
	t.title = s
}

//TODO: validate that 't' only includes ascii characters (rune < 255 i think)
func (t *Textbox) ChangeText(txt string) {
	if t.text != txt {
		t.text = txt
	}
}

//Render function optionally takes an offset (for containering), 2 or 3 ints.
func (t *Textbox) Render(offset ...int) {
	if t.visible {
		offX, offY, offZ := processOffset(offset)

		if t.bordered {
			console.DrawBorder(offX+t.x, offY+t.y, t.z+offZ, t.width, t.height, t.title, t.focused)
		}

		//word wrap calculatrix. a mighty sinful thing.
		//TODO: support for breaking super long words. right now it just skips the word.
		lines := make([]string, t.height)
		n := 0
		for _, s := range strings.Split(t.text, " ") {
			//super long word make-it-not-break hack.
			if len(s) > t.width {
				continue
			}

			if len(lines[n])+len(s) > t.width {
				//make sure we don't overflow the textbox
				if n < len(lines)-1 {
					n++
				} else {
					break
				}
			}
			lines[n] += s
			if len(lines[n]) != t.width {
				lines[n] += " "
			}
		}

		for l := 0; l < len(lines); l++ {
			offX := offX //so we can modify the offset separately for each line

			//fill textbox with background colour
			for i := 0; i < t.width*t.height; i++ {
				console.ChangeGridPoint(offX+t.x+i%t.width, offY+t.y+l, t.z+offZ, 0, 0xFFFFFFFF, 0xFF000000)
			}

			//offset if centerred
			if t.centered {
				offX += (t.width - len(lines[l])) / 2
			}

			//print text
			for i, r := range lines[l] {
				if i >= t.width {
					break
				}
				console.ChangeGridPoint(offX+t.x+i%t.width, offY+t.y+l, t.z+offZ, int(r), 0xFFFFFFFF, 0xFF000000)
			}
		}

		for i, _ := range t.anims {
			t.anims[i].Tick()
			t.anims[i].Render(t.x+offX, t.y+offY, t.z+offZ)
		}

	}
}

func (t Textbox) Dims() (int, int) {
	return t.width, t.height
}

func (t Textbox) Pos() (int, int, int) {
	return t.x, t.y, t.z
}

func (t *Textbox) ToggleVisible() {
	t.visible = !t.visible
	console.Clear()
}

func (t *Textbox) SetVisibility(v bool) {
	t.visible = v
	console.Clear()
}

func (t *Textbox) ToggleFocus() {
	t.focused = !t.focused
}

func (t *Textbox) SetCentered(c bool) {
	t.centered = c
}

func (t *Textbox) MoveTo(x, y, z int) {
	t.x = x
	t.y = y
	t.z = z
}
