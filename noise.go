package main

import (
	"fmt"
	"image"
	"image/color"
	"math"

	//"math/rand"
	"time"

	"github.com/anthonynsimon/bild/perlin"
)

var density = 1.0
var detail = 1.0/256.0
var depth = 256.0

func NewNormalNoise() (image.Image, error) {
	fmt.Printf("generating noise...\n")
	return NewNoise(density,detail)
}

func NewWave() (image.Image, error) {
	fmt.Printf("generating wave...\n")

	// Create a fucked noise image
	noise := perlin.NewPerlin(0.125,0.125,50,time.Now().Unix())

	img := image.NewNRGBA(image.Rect(0, 0, WIDTH, HEIGHT));
	// For each column in the image
	for y := float64(0); y < HEIGHT; y++ {
		// and each row
		for x := float64(0); x < WIDTH; x++ {
			value := math.Abs(noise.Noise2D(x,y))/8
			// Set the corresponding pixel
			img.Set(int(x),int(y), color.NRGBA{
				R: uint8(value),
				G: uint8(value),
				B: uint8(value),
				A: 255,
			})
		}
	}
	return img, nil
}


func NewNoise(density,detail float64) (image.Image, error) {
	// Create a noise image
	noise := perlin.NewPerlin(density,detail,100,time.Now().Unix())

	// Create a gradient to use for the colors
	colors, err := NewGradient()
	if(err != nil) {
		return nil, err
	}

	img := image.NewNRGBA(image.Rect(0, 0, WIDTH, HEIGHT));
	// For each column in the image
	for y := float64(0); y < HEIGHT; y++ {
		// and each row
		for x := float64(0); x < WIDTH; x++ {
			value := math.Abs(noise.Noise2D(x,y))
			fmt.Println(value)
			theColor := colors.At(value*2.0)
			// Set the corresponding pixel
			img.Set(int(x),int(y), color.NRGBA{
				R: uint8(theColor.R*255),
				G: uint8(theColor.G*255),
				B: uint8(theColor.B*255),
				A: 255,
			})
		}
	}
	return img, nil
}
