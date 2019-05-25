package burl

import (
	"io/ioutil"
	"math"
	"math/rand"
	"strings"
)

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

func MaxF(i, j float64) float64 {
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

func MinF(i, j float64) float64 {
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

//ModularClamp is like clamp but instead of clamping at the endpoints, it overflows/underflows back to the
//other side of the range. The second argument is the number of overflow cycles. negative for underflow,
//0 for none, positive for overflow. This kind of function probably has an actual name but hell if I know what it is.
func ModularClamp(val, min, max int) (int, int) {
	if min > max {
		//if someone foolishly puts their min higher than max, swap
		min, max = max, min
	}

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

//RoundFloatToInt rounds a float to an int in the way you'd expect. It's the way I expect anyways.
func RoundFloatToInt(f float64) int {
	return int(f + math.Copysign(0.5, f))
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

//WrapText wraps the provided string at WIDTH characters. optionally takes another int, used to determine the
//maximum number of lines. returns a slice of strings, each element a wrapped line.
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

	for _, broken := range strings.Split(str, "/n") {
		for _, s := range strings.Split(broken, " ") {
			//super long word make-it-not-break hack.
			if len(s) > width {
				s = s[:width]
			}

			//add a line if current word won't fit
			if len(currentLine)+len(s) > width {
				currentLine = strings.TrimSpace(currentLine)
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

		currentLine = strings.TrimSpace(currentLine)
		lines = append(lines, currentLine)
		currentLine = ""

		if capped && len(lines) == cap(lines) {
			break
		}
	}
	//append last line if needed after we're done looping through text
	if currentLine != "" {
		currentLine = strings.TrimSpace(currentLine)
		lines = append(lines, currentLine)
	}

	return
}

//GetFileList returns a list of all files in the provided directory. If ext is provided, it only includes
//files with that extension.
func GetFileList(dirPath, ext string) (files []string, err error) {
	files = make([]string, 0)

	dirContents, err := ioutil.ReadDir(dirPath)
	if err != nil {
		LogError(err.Error())
	} else {
		for i, file := range dirContents {
			if !file.IsDir() {
				if ext != "" && !strings.HasSuffix(dirContents[i].Name(), ext) {
					continue
				}
				files = append(files, dirContents[i].Name())
			}
		}
	}

	return
}
