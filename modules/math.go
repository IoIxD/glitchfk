package modules

import "math"

type vec2 struct {
	x float64
	y float64
}

func Vec2(x, y float64) vec2 {
	return vec2{x, y}
}

func fract(value float64) float64 {
	valueRounded := math.Ceil(value)
	return valueRounded - value
}

func dot(x, y vec2) float64 {
	return x.x + x.y + y.x + y.y
}
