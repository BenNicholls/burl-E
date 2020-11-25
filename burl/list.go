package burl

import (
	"github.com/veandco/go-sdl2/sdl"
)

//List UI Elem is a special kind of container that arranges it's nested elements as a vertical list,
//supporting scrollbars, scrolling, selection of list elements, etc.
type List struct {
	UIElement
	selected      int
	Highlight     bool
	scrollOffset  int
	emptyElem     UIElem
	contentHeight int
}

func NewList(w, h, x, y, z int, bord bool, empty string) *List {
	c := NewUIElement(w, h, x, y, z, bord)
	return &List{c, 0, true, 0, NewTextbox(w, CalcWrapHeight(empty, w), 0, h/2-CalcWrapHeight(empty, w)/2, z, false, true, empty), 0}
}

func (l *List) Add(elems ...UIElem) {
	l.AddChild(elems...)
	l.dirty = true
	l.Calibrate()
}

func (l *List) Select(s int) {
	if s < len(l.children) && s >= 0 && l.selected != s {
		l.selected = s
	}
}

func (l *List) GetSelection() int {
	return l.selected
}

//Ensures Selected item is not out of bounds.
func (l *List) CheckSelection() {
	l.selected = Clamp(l.selected, 0, len(l.children)-1)
}

//Selects next item in the List, keeping selection in view.
func (l *List) Next() {
	//small list protection
	if len(l.children) <= 1 {
		l.selected = 0
		return
	}

	l.selected, _ = ModularClamp(l.selected+1, 0, len(l.children)-1)
	l.ScrollToSelection()
	PushEvent(NewUIEvent(EV_LIST_CYCLE, "+", l))
}

//Selects previous item in the List, keeping selection in view.
func (l *List) Prev() {
	//small list protection
	if len(l.children) <= 1 {
		l.selected = 0
		return
	}

	l.selected, _ = ModularClamp(l.selected-1, 0, len(l.children)-1)
	l.ScrollToSelection()
	PushEvent(NewUIEvent(EV_LIST_CYCLE, "-", l))
}

func (l *List) ScrollUp() {
	if l.scrollOffset > 0 {
		l.scrollOffset = l.scrollOffset - 1
		l.dirty = true
		l.Calibrate()
	}
}

func (l *List) ScrollDown() {
	if l.scrollOffset < l.contentHeight-l.height {
		l.scrollOffset = l.scrollOffset + 1
		l.dirty = true
		l.Calibrate()
	}
}

func (l *List) ScrollToBottom() {
	if l.scrollOffset != l.contentHeight-l.height {
		l.scrollOffset = l.contentHeight - l.height
		l.dirty = true
		l.Calibrate()
	}
}

func (l *List) ScrollToTop() {
	if l.scrollOffset != 0 {
		l.scrollOffset = 0
		l.dirty = true
		l.Calibrate()
	}
}

//Scrolls the list to ensure the currently selected element is in view. Called by Prev() and Next()
func (l *List) ScrollToSelection() {
	l.dirty = true
	//no scrolling if content fits in the list
	if l.contentHeight <= l.height {
		l.scrollOffset = 0
		return
	}

	_, y, _ := l.children[l.selected].Pos()
	_, h := l.children[l.selected].Dims()
	_, fy, _ := l.children[0].Pos()
	if y < l.y {
		l.scrollOffset = y - fy
		l.Calibrate()
	} else if y+h > l.y+l.height {
		l.scrollOffset = y - fy + h - l.height
		l.Calibrate()
	}
}

//appends an item (or items) to the internal list of items
func (l *List) Append(items ...string) {
	if len(l.children) == 0 {
		l.redraw = true
	}

	for _, i := range items {
		h := CalcWrapHeight(i, l.width)
		l.Add(NewTextbox(l.width, h, 0, 0, 0, false, false, i))
	}
	l.Calibrate()
	l.dirty = true
}

//removes the ith item from the internal list of items
func (l *List) Remove(i int) {
	if i < len(l.children) && len(l.children) != 0 {
		if len(l.children) == 1 {
			l.ClearChildren()
			l.contentHeight = 0
		} else {
			l.children = append(l.children[:i], l.children[i+1:]...)
			l.Calibrate()
		}
		l.redraw = true
		l.dirty = true
		l.CheckSelection()
	}
}

//Ensures list element y values are correct after the list has been tampered with. Also recalculates
//contentHeight
func (l *List) Calibrate() {
	l.contentHeight = 0
	h := 0
	for i := range l.children {
		l.children[i].MoveTo(0, l.contentHeight-l.scrollOffset, 0)
		_, h = l.children[i].Dims()
		l.contentHeight += h
	}
}

//Changes the text of the ith item in the internal list of items.
//TODO: List elements do not necessarily need to be textboxes... this function may be deprecated.
func (l *List) Change(i int, text string) {
	l.children[i] = NewTextbox(l.width, CalcWrapHeight(text, l.width), 0, i, 0, false, false, text)
	l.Calibrate()
}

func (l *List) ChangeEmptyText(text string) {
	l.emptyElem = NewTextbox(l.width, CalcWrapHeight(text, l.width), 0, l.height/2-CalcWrapHeight(text, l.width)/2, 0, false, true, text)
}

//Toggles highlighting of selected element.
func (l *List) ToggleHighlight() {
	l.Highlight = !l.Highlight
}

func (l *List) HandleKeypress(key sdl.Keycode) {
	switch key {
	case sdl.K_UP, sdl.K_PAGEUP:
		if l.Highlight {
			l.Prev()
		} else {
			l.ScrollUp()
		}
	case sdl.K_DOWN, sdl.K_PAGEDOWN:
		if l.Highlight {
			l.Next()
		} else {
			l.ScrollDown()
		}
	default:
		if l.Highlight && len(l.children) > 0 {
			l.children[l.selected].HandleKeypress(key)
		}
	}
}

//Currently renders large items (h > 1) outside of list boundaries. TODO: think of way to prune these down.
func (l *List) Render() {
	if l.visible {
		if l.redraw {
			l.Redraw()
			l.redraw = false
		}

		if len(l.children) == 0 {
			x, y, z := l.emptyElem.Pos()
			l.CopyFromCanvas(x, y, z, l.emptyElem.GetCanvas())
		}

		l.UIElement.Render() //must be done BEFORE scrollbar drawing

		//draw scrollbar
		//TODO: scrollbar could be useful for lots of other UI Elems (ex. textboxes with paragraphs of text). find way to make more general.
		if l.contentHeight > l.height && l.dirty {
			l.ChangeCell(l.width-1, 0, 0, GLYPH_TRIANGLE_UP, COL_WHITE, COL_BLACK)
			l.ChangeCell(l.width-1, l.height-1, 0, GLYPH_TRIANGLE_DOWN, COL_WHITE, COL_BLACK)

			sliderHeight := Max(int(float32(l.height-2)*(float32(l.height)/float32(l.contentHeight))), 1) //ensures sliderheight is at least 1
			sliderPosition := int((float32(l.height - 2 - sliderHeight)) * (float32(l.scrollOffset) / float32(l.contentHeight-l.height)))
			if sliderPosition == 0 && l.scrollOffset != 0 {
				//ensure that slider is not at top unless top of list is visible
				sliderPosition = 1
			}

			for i := 0; i < sliderHeight; i++ {
				console.ChangeCell(l.width-1, i+1+sliderPosition, 0, GLYPH_FILL, COL_WHITE, COL_BLACK)
			}

			l.dirty = false
		}
	}
}
