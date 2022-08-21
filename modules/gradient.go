package modules

import (
	"fmt"
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
		"horizontal", "vertical", "diagonal", "radial", "inverse-radial", "mario64windowtexture",
	}
	for _, v := range types {
		FunctionPool.Add(v, NewGradientFunction(v, false))
		FunctionPool.Add("segmented-"+v, NewGradientFunction(v, true))

		// add them three more times to artifically increase their chances over the hoard of noise functions
		for i := 0; i < 3; i++ {
			FunctionPool.Add(v+"_"+fmt.Sprintf("%v",i), NewGradientFunction(v, false))
			FunctionPool.Add("segmented-"+v+"_"+fmt.Sprintf("%v",i), NewGradientFunction(v, true))
		}
	}
}

func NewGradientFunction(mode string, fract bool) ImageFunction {
	return func(width, height float64) (image.Image, error) {
		return NewGradientImage(mode, fract, width, height)
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
func NewGradientImage(mode string, useFract bool, width, height float64) (image.Image, error) {
	grad, err := NewGradient()
	if err != nil {
		return nil, err
	}

	img := image.NewNRGBA(image.Rect(0, 0, int(width), int(height)))
	offsetX := float64(rand.Intn(int(width)))
	offsetY := float64(rand.Intn(int(height)))
	mul := float64(rand.Intn(5))
	// For each column in the image
	for y_ := float64(0); y_ < height; y_++ {
		y := y_ + offsetY
		// and each row
		for x_ := float64(0); x_ < width; x_++ {
			x := x_ - offsetX
			var colorfulColor colorful.Color
			var position float64
			switch mode {
			case "horizontal":
				position = x / width
			case "vertical":
				position = y / height
			case "diagonal":
				position = (x / width) + (y / height)
			case "radial":
				position = math.Cos(x/width)*mul + math.Sin(y/height)*mul
			case "inverse-radial":
				position = math.Cos(y/height)*mul + math.Sin(x/width)*mul
			case "mario64windowtexture":
				adj := ((width / 2) - x)
				opp := ((height / 2) - y)
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