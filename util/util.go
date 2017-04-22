package util

import "math/rand"

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

//reports distance squared (sqrt unnecessary usually)
func Distance(x1, x2, y1, y2 int) int {
	return (x1-x2)*(x1-x2) + (y1-y2)*(y1-y2)
}

//Ensure (x,y) are inside (0, 0, w, h)
func CheckBounds(x, y, w, h int) bool {
	return x >= 0 && x < w && y >= 0 && y < h
}
