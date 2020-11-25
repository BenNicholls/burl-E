package burl

//UI Element for displaying text.
type Textbox struct {
	UIElement
	text     string
	centered bool
	lines    []string
}

func NewTextbox(w, h, x, y, z int, bord, cent bool, txt string) *Textbox {
	return &Textbox{NewUIElement(w, h, x, y, z, bord), txt, cent, WrapText(txt, w*2, h)}
}

//Returns the height required to fit a string after it has been wrapped.
func CalcWrapHeight(s string, width int) int {
	return len(WrapText(s, width*2))
}

//Replaces the string for the textbox.
func (t *Textbox) ChangeString(txt string) {
	if t.text != txt {
		t.text = txt
		t.lines = WrapText(txt, t.width*2, t.height)
		t.dirty = true
	}
}

//Adds text to the contents of textbox.
func (t *Textbox) AppendText(txt string) {
	t.ChangeString(t.text + txt)
}

//LoremIpsum() fills the textbox with a paragraph of crazy latin. For testing!
func (t *Textbox) LoremIpsum() {
	t.ChangeString(loremipsum)
}

func (t *Textbox) Render() {
	if t.visible {
		for l, line := range t.lines {
			lineOffset := 0

			//offset if centered
			if t.centered {
				lineOffset = (t.width*2 - len(line)) / 2
				//blank out area before text
				for i := 0; i < lineOffset; i++ {
					t.ChangeChar(i/2, l, 0, int(' '), i%2)
					t.ChangeColours(i/2, l, 0, t.foreColour, t.backColour)
				}
			}

			if line != "" {
				t.DrawText(lineOffset/2, l, 0, line, t.foreColour, t.backColour, lineOffset%2)
			}

			//blank out area after text
			for i := lineOffset + len(line); i < t.width*2; i++ {
				t.ChangeChar(i/2, l, 0, int(' '), i%2)
				t.ChangeColours(i/2, l, 0, t.foreColour, t.backColour)
			}
		}

		//blank out empty lines at bottom
		for y := len(t.lines); y < t.height; y++ {
			for x := 0; x < t.width; x++ {
				t.ChangeText(x, y, 0, int(' '), int(' '))
				t.ChangeColours(x, y, 0, t.foreColour, t.backColour)
			}
		}

		t.UIElement.Render()
	}
}

func (t *Textbox) SetCentered(c bool) {
	t.centered = c
}

const loremipsum = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed finibus velit a tempor condimentum. Nam accumsan aliquam feugiat. Pellentesque lobortis iaculis orci vel consectetur. Etiam tincidunt ipsum ac leo vehicula, et malesuada nulla dapibus. Quisque ultricies ultricies metus, in elementum enim suscipit sit amet. Ut consectetur nisl vitae metus eleifend fringilla. Vestibulum purus nunc, bibendum ullamcorper lacinia a, suscipit vel urna."
