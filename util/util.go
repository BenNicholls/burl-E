package util

import "math/rand"
import "strings"

//Interface for objects that can report a bounding box of some kind.
type Bounded interface {
	Rect() (int, int, int, int)
}

//checks if key is a letter or number (ASCII-encoded)
func ValidText(key rune) bool {
	return (key >= 93 && key < 123) || (key >= 48 && key < 58)
}

//generates a tuple of cartesian directions (cannot be 0,0)
func GenerateDirection() (int, int) {
	for {
		dx, dy := rand.Intn(3)-1, rand.Intn(3)-1
		if dx != 0 || dy != 0 {
			return dx, dy
		}
	}
}

//generates a random (x,y) pair within a box defined by (x, y, w, h)
func GenerateCoord(x, y, w, h int) (int, int) {
	return rand.Intn(w) + x, rand.Intn(h) + y
}

//reports distance squared (sqrt unnecessary usually)
func Distance(x1, y1, x2, y2 int) int {
	return (x1-x2)*(x1-x2) + (y1-y2)*(y1-y2)
}

//Ensure (x,y) are inside (0, 0, w, h)
func CheckBounds(x, y, w, h int) bool {
	return x >= 0 && x < w && y >= 0 && y < h
}

//Returns the max of two integers. Duh. If tied, returns the first argument.
func Max(i, j int) int {
	if i < j {
		return j
	} else {
		return i
	}
}

//Opposite of max. If tied, returns first argument.
func Min(i, j int) int {
	if i > j {
		return j
	} else {
		return i
	}
}

//returns the intersection of two rectangularly-bound objects as a rect
//if no intersection, returns 0,0,0,0
func FindIntersectionRect(r1, r2 Bounded) (x, y, w, h int) {
	x1, y1, w1, h1 := r1.Rect()
	x2, y2, w2, h2 := r2.Rect()

	x, y, w, h = 0, 0, 0, 0

	//check for intersection
	if x1 >= x2+w2 || x2 >= x1+w1 || y1 >= y2+h2 || y2 >= y1+h1 {
		return
	}

	x = Max(x1, x2)
	y = Max(y1, y2)
	w = Min(x1+w1, x2+w2) - x
	h = Min(y1+h1, y2+h2) - y

	return
}

//wraps the provided string at WIDTH characters. optionally takes another int, used to determine the maximum number of lines.
//returns a slice of strings, each element a wrapped line.
//for words longer than width, just brutally cuts them off. no mercy.
func WrapText(str string, width int, maxlines ...int) (lines []string) {
	capped := false
	if len(maxlines) == 1 {
		lines = make([]string, 0, maxlines[0])
		capped = true
	} else {
		lines = make([]string, 0)
	}

	currentLine := ""

	for _, s := range strings.Split(str, " ") {
		//super long word make-it-not-break hack.
		if len(s) > width {
			s = s[:width]
		}

		//add a line if current word won't fit
		if len(currentLine)+len(s) > width {
			lines = append(lines, currentLine)
			currentLine = ""

			//break if number of lines == height
			if capped && len(lines) == cap(lines) {
				break
			}

		}
		currentLine += s
		if len(currentLine) != width {
			currentLine += " "
		}
	}
	//append last line if needed after we're done looping through text
	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return
}
