package util

import "math/rand"
import "strings"
import "math"

//Bounded defines objects that can report a bounding box of some kind.
type Bounded interface {
	Rect() (int, int, int, int)
}

//ValidText checks if key is a letter or number or basic punctuation (ASCII-encoded)
//TODO: this is NOT comprehensive. Improve this later.
func ValidText(key rune) bool {
	return (key >= 93 && key < 123) || (key >= 37 && key < 58)
}

//RandomDirection generates a tuple of cartesian directions (cannot be 0,0)
func RandomDirection() (int, int) {
	for {
		dx, dy := rand.Intn(3)-1, rand.Intn(3)-1
		if dx != 0 || dy != 0 {
			return dx, dy
		}
	}
}

//GenerateCoord generates a random (x,y) pair within a box defined by (x, y, w, h)
func GenerateCoord(x, y, w, h int) (int, int) {
	return rand.Intn(w) + x, rand.Intn(h) + y
}

//Distance calculates the distance squared (sqrt unnecessary usually)
func Distance(x1, y1, x2, y2 int) int {
	return (x1-x2)*(x1-x2) + (y1-y2)*(y1-y2)
}

//ManhattanDistance calculates the manhattan (or taxicab) distance on a square grid.
func ManhattanDistance(x1, y1, x2, y2 int) int {
	return Abs(x2-x1) + Abs(y2-y1)
}

//CheckBounds ensures (x,y) is inside (0, 0, w, h)
func CheckBounds(x, y, w, h int) bool {
	return x >= 0 && x < w && y >= 0 && y < h
}

//IsInside checks if the point (px, py) is within the rect (x, y, w h).
func IsInside(px, py, x, y, w, h int) bool {
	return px >= x && px < x+w && py >= y && py < y+h
}

//Pow is an integer power function. Doesn't ~~do~~ negative exponents. Totally does 0 though.
func Pow(val, exp int) int {
	v := 1
	for i := 0; i < exp; i++ {
		v = v * val
	}
	return v
}

//Abs returns the absolute value of val
func Abs(val int) int {
	if val < 0 {
		return val * (-1)
	}
	return val
}

//Max returns the max of two integers. Duh.
func Max(i, j int) int {
	if i < j {
		return j
	} else {
		return i
	}
}

//Min is the opposite of max.
func Min(i, j int) int {
	if i > j {
		return j
	} else {
		return i
	}
}

//Clamp checks if min <= val <= max.
//If val < min, returns min. If val > max, returns max. Otherwise returns val.
func Clamp(val, min, max int) int {
	if val <= min {
		return min
	} else if val >= max {
		return max
	} else {
		return val
	}
}

//ModularClamp is like clamp but instead of clamping at the endpoints, it overflows/underflows back to the other side of the range.
//The second argument is the number of overflow cycles. negative for underflow, 0 for none, positive for overflow.
//This kind of function probably has an actual name but hell if I know what it is.
func ModularClamp(val, min, max int) (int, int) {
	if val < min {
		r := max - min + 1
		underflows := (min-val-1)/r + 1
		return val + r*underflows, -underflows
	} else if val > max {
		r := max - min + 1
		overflows := (val-max-1)/r + 1
		return val - r*overflows, overflows
	} else {
		return val, 0
	}
}

//RoundFLoatToInt rounds a float to an int in the way you'd expect. It's the way I expect anyways. 
func RoundFloatToInt(f float64) int {
	return int(f + math.Copysign(0.5, f))
}

//FindIntersectionRect calculates the intersection of two rectangularly-bound objects as a rect
//if no intersection, returns (0,0,0,0)
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

//Lerp linearly interpolates a range (min-max) over (steps) intervals, and returns the (val)th step.
//Currently does this via a conversion to float64, so there might be some rounding errors in here I don't know about.
func Lerp(min, max, val, steps int) int {
	if val >= steps {
		return max
	} else if val <= 0 {
		return min
	}

	stepVal := float64(max-min) / float64(steps)
	return int(float64(min) + stepVal*float64(val))
}

//WrapText wraps the provided string at WIDTH characters. optionally takes another int, used to determine the maximum number of lines.
//returns a slice of strings, each element a wrapped line.
//for words longer than width it just brutally cuts them off. no mercy.
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
