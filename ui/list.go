package ui

import "github.com/bennicholls/delvetown/console"

//List UI Elem is a special kind of container that arranges it's nested elements as a vertical list,
//supporting scrollbars, scrolling, selection of list elements, etc.
type List struct {
	Container
	selected      int
	Highlight     bool
	scrollOffset  int
	empty         bool
	emptyElem     UIElem
	contentHeight int
}

func NewList(w, h, x, y, z int, bord bool, empty string) *List {
	c := NewContainer(w, h, x, y, z, bord)
	return &List{*c, 0, true, 0, true, NewTextbox(w, CalcWrapHeight(empty, w), 0, h/2-CalcWrapHeight(empty, w)/2, z, false, true, empty), 0}
}

func (l *List) Select(s int) {
	if s < len(l.Elements) && s >= 0 {
		l.selected = s
	}
}

func (l List) GetSelection() int {
	return l.selected
}

//Ensures Selected item is not out of bounds.
func (l *List) CheckSelection() {
	if l.selected < 0 {
		l.selected = 0
	} else if l.selected >= len(l.Elements) {
		l.selected = len(l.Elements) - 1
	}
}

//Selects next item in the List, keeping selection in view.
func (l *List) Next() {
	//small list protection
	if len(l.Elements) <= 1 {
		l.selected = 0
		return
	}

	if l.selected >= len(l.Elements)-1 {
		l.selected = 0
	} else {
		l.selected++
	}

	l.ScrollToSelection()

	PushEvent(l, CHANGE, "List Cycled +")
}

//Selects previous item in the List, keeping selection in view.
func (l *List) Prev() {
	//small list protection
	if len(l.Elements) <= 1 {
		l.selected = 0
		return
	}

	if l.selected == 0 {
		l.selected = len(l.Elements) - 1
	} else {
		l.selected--
	}

	l.ScrollToSelection()

	PushEvent(l, CHANGE, "List Cycled -")
}

func (l *List) ScrollUp() {
	if l.scrollOffset != 0 {
		l.scrollOffset = l.scrollOffset - 1
	}
}

func (l *List) ScrollDown() {
	if l.scrollOffset < l.contentHeight-l.height {
		l.scrollOffset = l.scrollOffset + 1
	}
}

func (l *List) ScrollToBottom() {
	l.scrollOffset = l.contentHeight - l.height
}

func (l *List) ScrollToTop() {
	l.scrollOffset = 0
}

func (l *List) ScrollToSelection() {
	if l.selected < l.scrollOffset {
		l.scrollOffset = l.selected
	} else if l.scrollOffset < l.selected-l.height+1 {
		l.scrollOffset = l.selected - l.height + 1
	}
}

//appends an item (or items) to the internal list of items
func (l *List) Append(items ...string) {
	if len(l.Elements) == 0 {
		l.redraw = true
	}

	for _, i := range items {
		h := CalcWrapHeight(i, l.width)
		l.Add(NewTextbox(l.width, h, 0, l.contentHeight, 0, false, false, i))
		l.contentHeight += h
	}
	l.Calibrate()
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
		l.CheckSelection()
	}
}

//Ensures list element y values are correct after the list has been tampered with. Also recalculates
//contentHeight
func (l *List) Calibrate() {
	y := 0
	h := 0
	for i := range l.Elements {
		l.Elements[i].MoveTo(0, y, 0)
		_, h = l.Elements[i].Dims()
		y += h
	}
	l.contentHeight = y
}

//Changes the text of the ith item in the internal list of items.
//TODO: List elements do not necessarily need to be textboxes... this function may be deprecated.
func (l *List) Change(i int, item string) {
	l.Elements[i] = NewTextbox(l.width, CalcWrapHeight(item, l.width), 0, i, l.z, false, false, item)
	l.Calibrate()
}

//Toggles highlighting of selected element.
func (l *List) ToggleHighlight() {
	l.Highlight = !l.Highlight
}

//Currently renders large items (h > 1) outside of list boundaries. TODO: think of way to prune these down.
func (l *List) Render(offset ...int) {
	if l.visible {
		offX, offY, offZ := processOffset(offset)

		if l.redraw {
			console.Clear(l.x+offX, l.y+offY, l.width, l.height)
			l.redraw = false
		}

		if len(l.Elements) <= 0 {
			l.emptyElem.Render(l.x+offX, l.y+offY, l.z+offZ)
		} else {

			for _, e := range l.Elements {
				_, y, _ := e.Pos()
				_, h := e.Dims()
				if y < l.scrollOffset+l.height && y+h > l.scrollOffset {
					e.Render(l.x+offX, l.y+offY-l.scrollOffset, l.z+offZ)
				}
			}

			if l.Highlight {
				w, h := l.Elements[l.selected].Dims()
				_, y, _ := l.Elements[l.selected].Pos()
				for j := 0; j < h; j++ {
					for i := 0; i < w; i++ {
						console.Invert(offX+l.x+i, offY+l.y-l.scrollOffset+j+y, offZ+l.z)
					}
				}
			}
		}

		if l.bordered {
			console.DrawBorder(l.x+offX, l.y+offY, l.z+offZ, l.width, l.height, l.title, l.focused)
		}

		//draw scrollbar
		//TODO: scrollbar could be useful for lots of other UI Elems (ex. textboxes with paragraphs of text). find way to make more general.
		if l.contentHeight > l.height {
			console.ChangeGridPoint(offX+l.x+l.width, offY+l.y, offZ+l.z, 0x1e, 0xFFFFFFFF, 0xFF000000)
			console.ChangeGridPoint(offX+l.x+l.width, offY+l.y+l.height-1, offZ+l.z, 0x1f, 0xFFFFFFFF, 0xFF000000)

			sliderHeight := int(float32(l.height-2) * (float32(l.height) / float32(l.contentHeight)))
			if sliderHeight < 1 {
				sliderHeight = 1
			}

			sliderPosition := int((float32(l.height - 2 - sliderHeight)) * (float32(l.scrollOffset) / float32(l.contentHeight-l.height)))
			if sliderPosition == 0 && l.scrollOffset != 0 {
				//ensure that slider is not at top unless top of list is visible
				sliderPosition = 1
			}

			for i := 0; i < sliderHeight; i++ {
				console.ChangeGridPoint(offX+l.x+l.width, offY+l.y+i+1+sliderPosition, offZ+l.z, 0xb1, 0xFFFFFFFF, 0xFF000000)
			}
		}
	}
}
