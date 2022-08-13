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
	FunctionPool.Add("wave", func() (image.Image, error) {
		return NewNoiseLegacy(0.125, 0.125, 8, 2, false)
	})
	FunctionPool.Add("perlin-fbm", func() (image.Image, error) {
		return NewNoise(fastnoiselite.NoiseTypePerlin, fastnoiselite.FractalTypeFBm, 0.002, 5, false)
	})
	FunctionPool.Add("perlin-ridged", func() (image.Image, error) {
		return NewNoise(fastnoiselite.NoiseTypePerlin, fastnoiselite.FractalTypeRidged, 0.002, 5, false)
	})
	FunctionPool.Add("perlin-pingpong", func() (image.Image, error) {
		return NewNoise(fastnoiselite.NoiseTypePerlin, fastnoiselite.FractalTypePingPong, 0.002, 5, false)
	})
	FunctionPool.Add("perlin-domain-warp-progressive", func() (image.Image, error) {
		return NewNoise(fastnoiselite.NoiseTypePerlin, fastnoiselite.FractalTypeDomainWarpProgressive, 0.002, 5, false)
	})
	FunctionPool.Add("perlin-domain-warp-independent", func() (image.Image, error) {
		return NewNoise(fastnoiselite.NoiseTypePerlin, fastnoiselite.FractalTypeDomainWarpIndependent, 0.002, 5, false)
	})

	FunctionPool.Add("opensimplex2-fbm", func() (image.Image, error) {
		return NewNoise(fastnoiselite.NoiseTypeOpenSimplex2, fastnoiselite.FractalTypeFBm, 0.002, 5, false)
	})
	FunctionPool.Add("opensimplex2-ridged", func() (image.Image, error) {
		return NewNoise(fastnoiselite.NoiseTypeOpenSimplex2, fastnoiselite.FractalTypeRidged, 0.002, 5, false)
	})
	FunctionPool.Add("opensimplex2-pingpong", func() (image.Image, error) {
		return NewNoise(fastnoiselite.NoiseTypeOpenSimplex2, fastnoiselite.FractalTypePingPong, 0.002, 5, false)
	})
	FunctionPool.Add("opensimplex2-domain-warp-progressive", func() (image.Image, error) {
		return NewNoise(fastnoiselite.NoiseTypeOpenSimplex2, fastnoiselite.FractalTypeDomainWarpProgressive, 0.002, 5, false)
	})
	FunctionPool.Add("opensimplex2-domain-warp-independent", func() (image.Image, error) {
		return NewNoise(fastnoiselite.NoiseTypeOpenSimplex2, fastnoiselite.FractalTypeDomainWarpIndependent, 0.002, 5, false)
	})

	FunctionPool.Add("cellular-fbm", func() (image.Image, error) {
		return NewNoise(fastnoiselite.NoiseTypeCellular, fastnoiselite.FractalTypeFBm, 0.002, 5, false)
	})
	FunctionPool.Add("cellular-ridged", func() (image.Image, error) {
		return NewNoise(fastnoiselite.NoiseTypeCellular, fastnoiselite.FractalTypeRidged, 0.002, 5, false)
	})
	FunctionPool.Add("cellular-pingpong", func() (image.Image, error) {
		return NewNoise(fastnoiselite.NoiseTypeCellular, fastnoiselite.FractalTypePingPong, 0.002, 5, false)
	})
	FunctionPool.Add("cellular-domain-warp-progressive", func() (image.Image, error) {
		return NewNoise(fastnoiselite.NoiseTypeCellular, fastnoiselite.FractalTypeDomainWarpProgressive, 0.002, 5, false)
	})
	FunctionPool.Add("cellular-domain-warp-independent", func() (image.Image, error) {
		return NewNoise(fastnoiselite.NoiseTypeCellular, fastnoiselite.FractalTypeDomainWarpIndependent, 0.002, 5, false)
	})

	FunctionPool.Add("valuecubic-fbm", func() (image.Image, error) {
		return NewNoise(fastnoiselite.NoiseTypeValueCubic, fastnoiselite.FractalTypeFBm, 0.002, 5, false)
	})
	FunctionPool.Add("valuecubic-ridged", func() (image.Image, error) {
		return NewNoise(fastnoiselite.NoiseTypeValueCubic, fastnoiselite.FractalTypeRidged, 0.002, 5, false)
	})
	FunctionPool.Add("valuecubic-pingpong", func() (image.Image, error) {
		return NewNoise(fastnoiselite.NoiseTypeValueCubic, fastnoiselite.FractalTypePingPong, 0.002, 5, false)
	})
	FunctionPool.Add("valuecubic-domain-warp-progressive", func() (image.Image, error) {
		return NewNoise(fastnoiselite.NoiseTypeValueCubic, fastnoiselite.FractalTypeDomainWarpProgressive, 0.002, 5, false)
	})
	FunctionPool.Add("valuecubic-domain-warp-independent", func() (image.Image, error) {
		return NewNoise(fastnoiselite.NoiseTypeValueCubic, fastnoiselite.FractalTypeDomainWarpIndependent, 0.002, 5, false)
	})

	FunctionPool.Add("value-fbm", func() (image.Image, error) {
		return NewNoise(fastnoiselite.NoiseTypeValue, fastnoiselite.FractalTypeFBm, 0.002, 5, false)
	})
	FunctionPool.Add("value-ridged", func() (image.Image, error) {
		return NewNoise(fastnoiselite.NoiseTypeValue, fastnoiselite.FractalTypeRidged, 0.002, 5, false)
	})
	FunctionPool.Add("value-pingpong", func() (image.Image, error) {
		return NewNoise(fastnoiselite.NoiseTypeValue, fastnoiselite.FractalTypePingPong, 0.002, 5, false)
	})
	FunctionPool.Add("value-domain-warp-progressive", func() (image.Image, error) {
		return NewNoise(fastnoiselite.NoiseTypeValue, fastnoiselite.FractalTypeDomainWarpProgressive, 0.002, 5, false)
	})
	FunctionPool.Add("value-domain-warp-independent", func() (image.Image, error) {
		return NewNoise(fastnoiselite.NoiseTypeValue, fastnoiselite.FractalTypeDomainWarpIndependent, 0.002, 5, false)
	})
}

func NewNoise(noiseType fastnoiselite.NoiseType, fractalType fastnoiselite.FractalType, frequency float64, octaves int32, hasColor bool) (image.Image, error) {
	// Create a noise image
	noise := fastnoiselite.NewNoise()
	noise.Seed = int32(time.Now().UnixNano())
	noise.SetNoiseType(noiseType)
	noise.FractalType = fractalType
	noise.Frequency = frequency
	noise.SetFractalOctaves(octaves)

	var wg sync.WaitGroup
	wg.Add(int(WIDTH * HEIGHT))

	img := image.NewNRGBA(image.Rect(0, 0, int(WIDTH), int(HEIGHT)))

	// Create a gradient to use for the colors
	colors, err := NewGradient()
	if err != nil {
		return nil, err
	}

	// For each column in the image
	for y := float64(0); y < HEIGHT; y++ {
		// and each row
		for x := float64(0); x < WIDTH; x++ {
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

func NewNoiseLegacy(density, detail float64, divide, mul float64, hasColor bool) (image.Image, error) {
	// Create a noise image
	noise := perlin.NewPerlin(density, detail, 50, time.Now().Unix())

	var wg sync.WaitGroup
	wg.Add(int(WIDTH * HEIGHT))

	// Create a gradient to use for the colors
	colors, err := NewGradient()
	if err != nil {
		return nil, err
	}

	img := image.NewNRGBA(image.Rect(0, 0, int(WIDTH), int(HEIGHT)))

	// For each column in the image
	for y := float64(0); y < HEIGHT; y++ {
		// and each row
		for x := float64(0); x < WIDTH; x++ {
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
