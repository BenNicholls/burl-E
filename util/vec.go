package util

import "math"

type Vec2 struct {
	X, Y int
}

func (v *Vec2) Set(x, y int) {
	v.X, v.Y = x, y
}

func (v Vec2) Get() (int, int) {
	return v.X, v.Y
}

func (v *Vec2) Mod(dx, dy int) {
	v.X += dx
	v.Y += dy
}

func (v1 Vec2) Add(v2 Vec2) Vec2 {
	return Vec2{v1.X + v2.X, v1.Y + v2.Y}
}

//returns vec2 = v1 - v2
func (v1 Vec2) Sub(v2 Vec2) Vec2 {
	return Vec2{v1.X - v2.X, v1.Y - v2.Y}
}
 
func (v Vec2) MagFloat() float64 {
	return math.Sqrt(float64(v.X)*float64(v.X) + float64(v.Y)*float64(v.X))
 }

func (v Vec2) Mag() int {
	return int(v.MagFloat())
}

func (v Vec2) ToPolar() Vec2Polar {
	return Vec2Polar{v.MagFloat(), math.Atan2(float64(v.Y), float64(v.X))}
}

type Vec2Polar struct {
	R, Phi float64
}

func (v *Vec2Polar) Set(r, phi float64) {
	v.R, v.Phi = r, phi
}

func (v Vec2Polar) Get() (float64, float64) {
	return v.R, v.Phi
}

func (v Vec2Polar) ToRect() Vec2 {
	return Vec2{int(v.R*math.Cos(v.Phi)), int(v.R*math.Sin(v.Phi))}
}