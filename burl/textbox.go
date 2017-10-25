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

//LoremIpsum() fills the textbox with a paragraph of crazy latin. For testing!
func (t *Textbox) LoremIpsum() {
	t.ChangeText(loremipsum)
}

//Render function optionally takes an offset (for containering), 2 or 3 ints.
func (t *Textbox) Render(offset ...int) {
	if t.visible {
		offX, offY, offZ := processOffset(offset)

		for l, line := range t.lines {
			lineOffset := 0

			//blank out line, lets start fresh
			for i := 0; i < t.width; i++ {
				console.ChangeText(t.x+offX+i, t.y+offY+l, t.z+offZ, int(' '), int(' '))
			}

			//offset if centered
			if t.centered {
				lineOffset = (t.width*2 - len(line) + 1) / 2
			}

			//print text
			console.DrawText(offX+t.x+lineOffset/2, offY+t.y+l, offZ+t.z, line, COL_WHITE, COL_BLACK, lineOffset%2)

		}

		//blank out empty lines at bottom
		for y := len(t.lines); y < t.height; y++ {
			for x := 0; x < t.width; x++ {
				console.ChangeText(offX+t.x+x, offY+t.y+y, offZ+t.z, int(' '), int(' '))
			}
		}

		t.UIElement.Render(offX, offY, offZ)
	}
}

func (t *Textbox) SetCentered(c bool) {
	t.centered = c
}

const loremipsum = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed finibus velit a tempor condimentum. Nam accumsan aliquam feugiat. Pellentesque lobortis iaculis orci vel consectetur. Etiam tincidunt ipsum ac leo vehicula, et malesuada nulla dapibus. Quisque ultricies ultricies metus, in elementum enim suscipit sit amet. Ut consectetur nisl vitae metus eleifend fringilla. Vestibulum purus nunc, bibendum ullamcorper lacinia a, suscipit vel urna."