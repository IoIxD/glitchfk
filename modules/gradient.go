package modules

import (
	"image"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/mazznoer/colorgrad"
)

const DETAIL = 32768

func init() {
	types := []string{
		"horizontal", "vertical", "diagonal", "radial", "inverse-radial",
		"area", "mario64windowtexture",
	}
	for _, v := range types {
		FunctionPool.Add(v, NewGradientFunction(v, false))
		FunctionPool.Add("segmented-"+v, NewGradientFunction(v, true))
	}
}

func NewGradientFunction(mode string, fract bool) func() (image.Image, error) {
	return func() (image.Image, error) {
		return NewGradientImage(mode, fract)
	}
}

func NewGradient() (colorgrad.Gradient, error) {
	rand.Seed(time.Now().UnixNano())
	// Create a bunch of random colors
	colorNum := rand.Intn(DETAIL)
	colors := make([]color.Color, colorNum)

	for i := range colors {
		rand.Seed(time.Now().UnixNano())
		colors[i] = color.RGBA{uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255)), 255}
	}

	grad, err := colorgrad.NewGradient().Colors(colors[:]...).Domain(0, float64(colorNum)).Build()
	if err != nil {
		return colorgrad.Gradient{}, err
	}
	return grad, err
}
func NewGradientImage(mode string, useFract bool) (image.Image, error) {
	grad, err := NewGradient()
	if err != nil {
		return nil, err
	}

	img := image.NewNRGBA(image.Rect(0, 0, int(WIDTH), int(HEIGHT)))
	offsetX := float64(rand.Intn(int(WIDTH)))
	offsetY := float64(rand.Intn(int(HEIGHT)))
	mul := float64(rand.Intn(5))
	// For each column in the image
	for y_ := float64(0); y_ < HEIGHT; y_++ {
		y := y_ + offsetY
		// and each row
		for x_ := float64(0); x_ < WIDTH; x_++ {
			x := x_ - offsetX
			var colorfulColor colorful.Color
			var position float64
			switch mode {
			case "horizontal":
				position = x / WIDTH
			case "vertical":
				position = y / HEIGHT
			case "diagonal":
				position = (x / WIDTH) + (y / HEIGHT)
			case "radial":
				position = math.Cos(x/WIDTH)*mul + math.Sin(y/HEIGHT)*mul
			case "inverse-radial":
				position = math.Cos(y/HEIGHT)*mul + math.Sin(x/WIDTH)*mul
			case "mario64windowtexture":
				adj := ((WIDTH / 2) - x)
				opp := ((HEIGHT / 2) - y)
				ctan := opp / adj
				tan := adj / opp
				position = math.Cbrt(math.Abs(tan - ctan))
			}

			if useFract {
				position = fract(position)
			}

			colorfulColor = grad.At(position)

			// Set the corresponding pixel
			img.Set(int(x_), int(y_), color.NRGBA{
				R: uint8(colorfulColor.R * 256),
				G: uint8(colorfulColor.G * 256),
				B: uint8(colorfulColor.B * 256),
				A: 255,
			})
		}
	}
	return img, nil
}

func fract(value float64) (float64) {
	valueRounded := math.Ceil(value)
	return valueRounded - value
}