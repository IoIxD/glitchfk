package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"sync"

	//"math/rand"
	"time"

	"github.com/aquilax/go-perlin"
)

var density = 1.0
var detail = 1.0/256.0
var depth = 256.0

func NewNormalNoise() (image.Image, error) {
	fmt.Printf("generating noise...\n")
	return NewNoise(density,detail,1,255,true)
}

func NewWave() (image.Image, error) {
	fmt.Printf("generating wave...\n")
	return NewNoise(0.125,0.125,8,2,false)
}

func NewNoise(density,detail float64, divide, mul float64, hasColor bool) (image.Image, error) {
	start := time.Now().UnixMilli()

	// Create a noise image
	noise := perlin.NewPerlin(density,detail,50,time.Now().Unix())

	var wg sync.WaitGroup;
	wg.Add(int(WIDTH*HEIGHT))

	// Create a gradient to use for the colors
	colors, err := NewGradient()
	if(err != nil) {
		return nil, err
	}

	img := image.NewNRGBA(image.Rect(0, 0, int(WIDTH), int(HEIGHT)));

	// For each column in the image
	for y := float64(0); y < HEIGHT; y++ {
		// and each row
		for x := float64(0); x < WIDTH; x++ {
			// branch off into another thread
			go func(x,y float64) {
				// generate the noise value.
				value := math.Abs(noise.Noise2D(x,y))/divide
				var theColor color.NRGBA
				if(hasColor) {
					theColor_ := colors.At(value)
					theColor = color.NRGBA{
						R: uint8(theColor_.R*mul),
						G: uint8(theColor_.G*mul),
						B: uint8(theColor_.B*mul),
						A: 255,
					}
				} else {
					theColor = color.NRGBA{uint8(value),uint8(value),uint8(value),255}
				}
				// Set the corresponding pixel
				img.Set(int(x), int(y), theColor)
				wg.Done()
			}(x,y)
		}
	}

	wg.Wait()

	end := time.Now().UnixMilli()
	fmt.Println(end-start, "ms")
	return img, nil
}

/*
func NewWave() (image.Image, error) {
	fmt.Printf("generating wave...\n")

	var wg sync.WaitGroup;

	// Create a fucked noise image
	noise := perlin.NewPerlin(0.125,0.125,50,time.Now().Unix())

	img := image.NewNRGBA(image.Rect(0, 0, int(WIDTH), int(HEIGHT)));
	// For each column in the image
	for y := float64(0); y < HEIGHT; y++ {
		// and each row
		for x := float64(0); x < WIDTH; x++ {
			wg.Add(1)
			go func(x,y float64) {
				value := math.Abs(noise.Noise2D(x,y))/8
				// Set the corresponding pixel
				img.Set(int(x),int(y), color.NRGBA{
					R: uint8(value),
					G: uint8(value),
					B: uint8(value),
					A: 255,
				})
			}(x,y)
			wg.Done()
		}
	}
	wg.Wait()
	return img, nil
}

*/