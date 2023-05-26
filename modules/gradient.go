package modules

import (
	"image"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/mazznoer/colorgrad"
)

const DETAIL = 32768

func init() {
	FunctionPool.Lock()
	defer FunctionPool.Unlock()
	types := []string{
		"horizontal", "vertical", "diagonal", "radial", "inverse-radial", "mario64windowtexture",
	}
	for _, v := range types {
		FunctionPool.Add(v, NewGradientFunction(v, false))
		FunctionPool.Add("segmented-"+v, NewGradientFunction(v, true))
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

			var position float64

			switch mode {
			case "horizontal":												// 27ms 	-> 15ms
				position = x / width
			case "vertical":												// 27ms	 	-> 21ms
				position = y / height
			case "diagonal": 												// 27ms 	-> 18ms
				position = (x / width) + (y / height)
			case "radial": 													// 51ms 	-> 43.4ms
				position = cos(x/width)*mul + sin(y/height)*mul
			case "inverse-radial": 											// 51ms 	-> 36ms
				position = cos(y/height)*mul + sin(x/width)*mul
			case "mario64windowtexture": 									// 56ms 	-> 35ms
				adj := ((width / 2) - x)
				opp := ((height / 2) - y)
				ctan := opp / adj
				tan := adj / opp
				position = math.Cbrt(float64(Abs(tan - ctan)))
			}

			if useFract {
				position = fract(position)
			}

			colorfulColor := grad.At(position)

			// Set the corresponding color
			col := color.NRGBA{
				R: uint8(colorfulColor.R * 256),
				G: uint8(colorfulColor.G * 256),
				B: uint8(colorfulColor.B * 256),
				A: 255,
			}

			// golang function calls are too slow for us so  we'll just copy and paste the code for img.Set
			// here.

			if !(image.Point{int(x_), int(y_)}.In(img.Rect)) {
				continue
			}
			i := img.PixOffset(int(x_), int(y_))
			c1 := color.NRGBAModel.Convert(col).(color.NRGBA)
			s := img.Pix[i : i+4 : i+4] // Small cap improves performance, see https://golang.org/issue/27857
			s[0] = c1.R
			s[1] = c1.G
			s[2] = c1.B
			s[3] = 255

		}
	}
	return img, nil
}

func fract(value float64) (float64) {
	valueRounded := float64(int64(value))
	return valueRounded - value
}

