package modules

import (
	"image"
	"image/color"
	"math"
	"sync"

	//"math/rand"
	"time"

	"github.com/aquilax/go-perlin"
	"github.com/furui/fastnoiselite-go"
)

var density = 1.0
var detail = 1.0 / 256.0
var depth = 256.0

func init() {
	FunctionPool.Add("wave", func(width, height float64) (image.Image, error) {
		return NewNoiseLegacy(0.125, 0.125, 8, 2, false, width, height)
	})
	types := map[string]fastnoiselite.NoiseType{
		"perlin": fastnoiselite.NoiseTypePerlin,
		"opensimplex2": fastnoiselite.NoiseTypeOpenSimplex2,
		"cellular": fastnoiselite.NoiseTypeCellular,
		"valuecubic": fastnoiselite.NoiseTypeValueCubic,
		"value": fastnoiselite.NoiseTypeValue,
	}
	fractal := map[string]fastnoiselite.FractalType{
		"fbm": fastnoiselite.FractalTypeFBm,
		"ridged": fastnoiselite.FractalTypeRidged,
		"pingpong": fastnoiselite.FractalTypePingPong,
		"domain-warp-progressive": fastnoiselite.FractalTypeDomainWarpProgressive,
		"domain-warp-independent": fastnoiselite.FractalTypeDomainWarpIndependent,
	}

	for k1, v1 := range types {
		for k2, v2 := range fractal {
			FunctionPool.Add(k1+"-"+k2, func(width, height float64) (image.Image, error) {
				return NewNoise(v1, v2, 0.002, 5, false, width, height)
			})
			FunctionPool.Add(k1+"-"+k2+"-colored", func(width, height float64) (image.Image, error) {
				return NewNoise(v1, v2, 0.002, 5, true, width, height)
			})
		}
	}
}

func NewNoise(noiseType fastnoiselite.NoiseType, fractalType fastnoiselite.FractalType, frequency float64, octaves int32, hasColor bool, width, height float64) (image.Image, error) {
	// Create a noise image
	noise := fastnoiselite.NewNoise()
	noise.Seed = int32(time.Now().UnixNano())
	noise.SetNoiseType(noiseType)
	noise.FractalType = fractalType
	noise.Frequency = frequency
	noise.SetFractalOctaves(octaves)

	var wg sync.WaitGroup
	wg.Add(int(width * height))

	img := image.NewNRGBA(image.Rect(0, 0, int(width), int(height)))

	// Create a gradient to use for the colors
	colors, err := NewGradient()
	if err != nil {
		return nil, err
	}

	// For each column in the image
	for y := float64(0); y < height; y++ {
		// and each row
		for x := float64(0); x < width; x++ {
			// branch off into another thread
			go func(x, y float64) {
				// generate the noise value.
				value := math.Abs(noise.GetNoise2D(fastnoiselite.FNLfloat(x), fastnoiselite.FNLfloat(y)))
				var theColor color.NRGBA
				if hasColor {
					theColor_ := colors.At(value)
					theColor = color.NRGBA{
						R: uint8(theColor_.R * 255),
						G: uint8(theColor_.G * 255),
						B: uint8(theColor_.B * 255),
						A: 255,
					}
				} else {
					theColor = color.NRGBA{uint8(value*255), uint8(value*255), uint8(value*255), 255}
				}
				// Set the corresponding pixel
				img.Set(int(x), int(y), theColor)
				wg.Done()
			}(x, y)
		}
	}

	wg.Wait()
	return img, nil
}

func NewNoiseLegacy(density, detail float64, divide, mul float64, hasColor bool, width, height float64) (image.Image, error) {
	// Create a noise image
	noise := perlin.NewPerlin(density, detail, 50, time.Now().Unix())

	var wg sync.WaitGroup
	wg.Add(int(width * height))

	// Create a gradient to use for the colors
	colors, err := NewGradient()
	if err != nil {
		return nil, err
	}

	img := image.NewNRGBA(image.Rect(0, 0, int(width), int(height)))

	// For each column in the image
	for y := float64(0); y < height; y++ {
		// and each row
		for x := float64(0); x < width; x++ {
			// branch off into another thread
			go func(x, y float64) {
				// generate the noise value.
				value := math.Abs(noise.Noise2D(x, y)) / divide
				var theColor color.NRGBA
				if hasColor {
					theColor_ := colors.At(value)
					theColor = color.NRGBA{
						R: uint8(theColor_.R * mul),
						G: uint8(theColor_.G * mul),
						B: uint8(theColor_.B * mul),
						A: 255,
					}
				} else {
					theColor = color.NRGBA{uint8(value), uint8(value), uint8(value), 255}
				}
				// Set the corresponding pixel
				img.Set(int(x), int(y), theColor)
				wg.Done()
			}(x, y)
		}
	}

	wg.Wait()
	return img, nil
}

/*
func NewWave() (image.Image, error) {
	fmt.Printf("generating wave...\n")

	var wg sync.WaitGroup;

	// Create a fucked noise image
	noise := perlin.NewPerlin(0.125,0.125,50,time.Now().Unix())

	img := image.NewNRGBA(image.Rect(0, 0, int(width), int(height)));
	// For each column in the image
	for y := float64(0); y < height; y++ {
		// and each row
		for x := float64(0); x < width; x++ {
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
