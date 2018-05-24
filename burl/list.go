package burl

import (
	"github.com/veandco/go-sdl2/sdl"
)

//List UI Elem is a special kind of container that arranges it's nested elements as a vertical list,
//supporting scrollbars, scrolling, selection of list elements, etc.
type List struct {
	Container
	selected      int
	Highlight     bool
	scrollOffset  int
	emptyElem     UIElem
	contentHeight int
}

func NewList(w, h, x, y, z int, bord bool, empty string) *List {
	c := NewContainer(w, h, x, y, z, bord)
	return &List{*c, 0, true, 0, NewTextbox(w, CalcWrapHeight(empty, w), x, y+h/2-CalcWrapHeight(empty, w)/2, z, false, true, empty), 0}
}

func (l *List) Move(dx, dy, dz int) {
	l.Container.Move(dx, dy, dz)
	l.emptyElem.Move(dx, dy, dz)
}

func (l *List) Add(elems ...UIElem) {
	l.Container.Add(elems...)
	l.dirty = true
	l.Calibrate()
}

func (l *List) Select(s int) {
	if s < len(l.Elements) && s >= 0 && l.selected != s {
		l.selected = s
	}
}

func (l List) GetSelection() int {
	return l.selected
}

//Ensures Selected item is not out of bounds.
func (l *List) CheckSelection() {
	l.selected = Clamp(l.selected, 0, len(l.Elements)-1)
}

//Selects next item in the List, keeping selection in view.
func (l *List) Next() {
	//small list protection
	if len(l.Elements) <= 1 {
		l.selected = 0
		return
	}

	l.selected, _ = ModularClamp(l.selected+1, 0, len(l.Elements)-1)
	l.ScrollToSelection()
	PushEvent(NewUIEvent(EV_LIST_CYCLE, "+", l))
}

//Selects previous item in the List, keeping selection in view.
func (l *List) Prev() {
	//small list protection
	if len(l.Elements) <= 1 {
		l.selected = 0
		return
	}

	l.selected, _ = ModularClamp(l.selected-1, 0, len(l.Elements)-1)
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

	_, y, _ := l.Elements[l.selected].Pos()
	_, h := l.Elements[l.selected].Dims()
	_, fy, _ := l.Elements[0].Pos()
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
	if len(l.Elements) == 0 {
		l.redraw = true
	}

	for _, i := range items {
		h := CalcWrapHeight(i, l.width)
		l.Add(NewTextbox(l.width, h, l.x, l.y, l.z, false, false, i))
	}
	l.Calibrate()
	l.dirty = true
}

//removes the ith item from the internal list of items
func (l *List) Remove(i int) {
	if i < len(l.Elements) && len(l.Elements) != 0 {
		if len(l.Elements) == 1 {
			l.ClearElements()
			l.contentHeight = 0
		} else {
			l.Elements = append(l.Elements[:i], l.Elements[i+1:]...)
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
	for i := range l.Elements {
		l.Elements[i].MoveTo(l.x, l.contentHeight+l.y-l.scrollOffset, l.z)
		_, h = l.Elements[i].Dims()
		l.contentHeight += h
	}
}

//Changes the text of the ith item in the internal list of items.
//TODO: List elements do not necessarily need to be textboxes... this function may be deprecated.
func (l *List) Change(i int, item string) {
	l.Elements[i] = NewTextbox(l.width, CalcWrapHeight(item, l.width), 0, i, l.z, false, false, item)
	l.Calibrate()
}

func (l *List) ChangeEmptyText(text string) {
	l.emptyElem = NewTextbox(l.width, CalcWrapHeight(text, l.width), 0, l.height/2-CalcWrapHeight(text, l.width)/2, l.z, false, true, text)
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
		if l.Highlight && len(l.Elements) > 0 {
			l.Elements[l.selected].HandleKeypress(key)
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

		if len(l.Elements) <= 0 {
			l.emptyElem.Render()
		} else {
			for _, e := range l.Elements {
				_, y, _ := e.Pos()
				_, h := e.Dims()
				if y < l.y+l.height && y+h > l.y {
					e.Render()
				}
			}

			//TODO: implement highlight by inverting BACK and FORE when UIElement Colour Scheming goes in
			if l.Highlight {
				w, h := l.Elements[l.selected].Dims()
				_, y, _ := l.Elements[l.selected].Pos()
				for j := 0; j < h; j++ {
					for i := 0; i < w; i++ {
						console.Invert(l.x+i, j+y, l.z)
					}
				}
			}
		}

		l.Container.UIElement.Render() //must be done BEFORE scrollbar drawing

		//draw scrollbar
		//TODO: scrollbar could be useful for lots of other UI Elems (ex. textboxes with paragraphs of text). find way to make more general.
		if l.contentHeight > l.height && l.dirty {
			console.ChangeCell(l.x+l.width-1, l.y, l.z, GLYPH_TRIANGLE_UP, COL_WHITE, COL_BLACK)
			console.ChangeCell(l.x+l.width-1, l.y+l.height-1, l.z, GLYPH_TRIANGLE_DOWN, COL_WHITE, COL_BLACK)

			sliderHeight := Max(int(float32(l.height-2)*(float32(l.height)/float32(l.contentHeight))), 1) //ensures sliderheight is at least 1
			sliderPosition := int((float32(l.height - 2 - sliderHeight)) * (float32(l.scrollOffset) / float32(l.contentHeight-l.height)))
			if sliderPosition == 0 && l.scrollOffset != 0 {
				//ensure that slider is not at top unless top of list is visible
				sliderPosition = 1
			}

			for i := 0; i < sliderHeight; i++ {
				console.ChangeCell(l.x+l.width-1, l.y+i+1+sliderPosition, l.z, GLYPH_FILL, COL_WHITE, COL_BLACK)
			}

			l.dirty = false
		}
	}
}
