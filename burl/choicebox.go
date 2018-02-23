package burl

const (
	CHOICE_VERTICAL int = iota
	CHOICE_HORIZONTAL
)

//Choicebox is a textbox wherein once can cycle through some predefined choices.
//Only one choice is shown at a time. When focused, arrows appear to indicate
//which direction on the keyboard cycles the choices (vertical/horizontal).
//NOTE: This is very similar to a text-only list. Someday maybe these should be
//combined?
type ChoiceBox struct {
	Textbox

	choices   []string
	curChoice int
	direction int //see direction consts above
}

//NewChoiceBox. Textbox parameters as normal. dir parameter is either CHOICE_HORIZONTAL or
//CHOICE_VERTICAL. Defaults to horizontal. Then accepts any number of strings as choices.
func NewChoiceBox(w, h, x, y, z int, bord bool, dir int, choices ...string) (cb *ChoiceBox) {
	cb = new(ChoiceBox)
	cb.Textbox = *NewTextbox(w, h, x, y, z, bord, true, "")
	cb.choices = choices

	if dir == CHOICE_VERTICAL {
		cb.direction = dir
	} else {
		cb.direction = CHOICE_HORIZONTAL
	}

	if len(cb.choices) == 0 {
		cb.ChangeText("---")
	} else {
		cb.ChangeText(cb.choices[cb.curChoice])
	}

	return
}

func (cb *ChoiceBox) AddChoice(c string) {
	cb.choices = append(cb.choices, c)
}

func (cb *ChoiceBox) Next() {
	cb.curChoice, _ = ModularClamp(cb.curChoice+1, 0, len(cb.choices)-1)
	cb.ChangeText(cb.choices[cb.curChoice])
}

func (cb *ChoiceBox) Prev() {
	cb.curChoice, _ = ModularClamp(cb.curChoice-1, 0, len(cb.choices)-1)
	cb.ChangeText(cb.choices[cb.curChoice])
}

func (cb ChoiceBox) GetChoice() int {
	return cb.curChoice
}

func (cb ChoiceBox) Render(offset ...int) {
	if cb.visible {
		offX, offY, offZ := processOffset(offset)

		cb.Textbox.Render(offX, offY, offZ)

		//draw choice cycling triangles
		if cb.direction == CHOICE_HORIZONTAL {
			console.ChangeCell(cb.x+offX, cb.y+offY+cb.height/2, cb.z+offZ, GLYPH_TRIANGLE_LEFT, COL_WHITE, COL_BLACK)
			console.ChangeCell(cb.x+offX+cb.width-1, cb.y+offY+cb.height/2, cb.z+offZ, GLYPH_TRIANGLE_RIGHT, COL_WHITE, COL_BLACK)
		} else {
			console.ChangeCell(cb.x+offX+cb.width/2, cb.y+offY-1, cb.z+offZ, GLYPH_TRIANGLE_UP, COL_WHITE, COL_BLACK)
			console.ChangeCell(cb.x+offX+cb.width/2, cb.y+offY+cb.height, cb.z+offZ, GLYPH_TRIANGLE_DOWN, COL_WHITE, COL_BLACK)
		}

	}
}
