package main

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

const DETAIL = 32768;

func NewGradient() (colorgrad.Gradient,error) {
	rand.Seed(time.Now().UnixNano())
	// Create a bunch of random colors
	colorNum := rand.Intn(DETAIL)
	colors := make([]color.Color,colorNum)

	for i := range colors {
		rand.Seed(time.Now().UnixNano())
		colors[i] = color.RGBA{uint8(rand.Intn(255)),uint8(rand.Intn(255)),uint8(rand.Intn(255)),255}
	}

	grad, err := colorgrad.NewGradient().Colors(colors[:]...).Domain(0, float64(colorNum)).Build()
	if(err != nil) {
		return colorgrad.Gradient{}, err
	}
	return grad, err
}
func NewGradientImage(mode string) (image.Image, error) {
	fmt.Printf("generating %v gradient...\n",mode)

	grad, err := NewGradient()
	if(err != nil) {
		return nil, err
	}

	img := image.NewNRGBA(image.Rect(0, 0, int(WIDTH), int(HEIGHT)));
	offsetX := float64(rand.Intn(int(WIDTH)))
	offsetY := float64(rand.Intn(int(HEIGHT)))
	mul := float64(rand.Intn(5))
	// For each column in the image
	for y_ := float64(0); y_ < HEIGHT; y_++ {
		y := y_+offsetY
		// and each row
		for x_ := float64(0); x_ < WIDTH; x_++ {
			x := x_-offsetX
			var colorfulColor colorful.Color
			switch(mode) {
				case "horizontal": colorfulColor = grad.At(x/WIDTH)
				case "vertical": colorfulColor = grad.At(y/HEIGHT)
				case "diagonal": colorfulColor = grad.At((x/WIDTH)+(y/HEIGHT))
				case "radial":
					pos := math.Cos(x/WIDTH)*mul+math.Sin(y/HEIGHT)*mul
					colorfulColor = grad.At((pos-1))
				case "inverse-radial":
					pos := math.Cos(y/HEIGHT)*mul+math.Sin(x/WIDTH)*mul
					colorfulColor = grad.At((pos-1))
				case "mario64windowtexture":
					adj := ((WIDTH/2)-x) 		
					opp := ((HEIGHT/2)-y)
					ctan := opp/adj
					tan := adj/opp
					pos := math.Cbrt(math.Abs(tan-ctan))
					colorfulColor = grad.At(pos)
				}
			// Set the corresponding pixel
			img.Set(int(x_),int(y_), color.NRGBA{
				R: uint8(colorfulColor.R*256),
				G: uint8(colorfulColor.G*256),
				B: uint8(colorfulColor.B*256),
				A: 255,
			})
		}
	}
	return img, nil
}
