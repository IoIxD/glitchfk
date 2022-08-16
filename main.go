package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
	"sync"
	"time"

	"github.com/IoIxD/glitchfuckTwitter/modules"
	"github.com/pelletier/go-toml/v2"
)

var LocalConfig struct {
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string
	OAuthToken     string
	OAuthSecret    string
	Interval       string
	InProduction   bool
}

var imageTypes = flag.String("types", "", "The types of images you want to generate and xor, seperated by commas.")
var cpuProfile = flag.String("cpuprofile", "", "Debug flag for what file to write the cpu profile to. CPU Profiling is disabled if this is blank.")
var memProfile = flag.String("memprofile", "", "Debug flag for what file to write the memory profile to. CPU Profiling is disabled if this is blank.")
var outputFile = flag.String("output", "test.png", "File to save output to.")

func main() {
	flag.Parse()
	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			fmt.Println("could not create CPU profile: ", err)
			return
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			fmt.Println("could not start CPU profile: ", err)
			return
		}
		defer pprof.StopCPUProfile()
	}
	if *memProfile != "" {
		f, err := os.Create(*memProfile)
		if err != nil {
			fmt.Println("could not create memory profile: ", err)
			return
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			fmt.Println("could not write memory profile: ", err)
			return
		}
	}

	if *imageTypes != "" {
		types := strings.Split(*imageTypes, ",")
		var lastImage image.Image
		var finalImage image.Image
		for _, v := range types {
			image, err := NewImage(v)
			if err != nil {
				fmt.Println(err)
				return
			}
			if lastImage != nil {
				finalImage = xor(lastImage, image)
			} else {
				finalImage = image
			}
			lastImage = image
		}
		f, _ := os.Create(*outputFile)
		if err := png.Encode(f, finalImage); err != nil {
			fmt.Println(err)
		}
		return
	}

	// set up the config
	f, err := os.Open("config.toml")
	if err != nil {
		fmt.Println(err)
		// If the file simply doesn't exist, continue but
		// set the "in production" value to false.
		if !os.IsNotExist(err) {
			return
		}
		LocalConfig.InProduction = false
	} else {
		err = toml.NewDecoder(f).Decode(&LocalConfig)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if LocalConfig.InProduction {
		fmt.Println("starting twitter thread.")
		TwitterThread()
	} else {
		image := DefaultImage()
		f, _ := os.Create(*outputFile)
		f.Write(image)
		f.Close()
	}

}

func WaitFor(int time.Duration) <-chan time.Time {
	now := time.Now()
	dur := now.Truncate(int).Add(int).Sub(now)
	return time.After(dur)
}

func DefaultImage() []byte {
	var approved bool
	var finalgrad image.Image
	for(!approved) {
		var grad1, grad2 image.Image
		var wg sync.WaitGroup
		var err error

		wg.Add(2)
		go func() {
			grad1, err = modules.FunctionPool.Random()()
			if err != nil {
				fmt.Println(err)
			}
			wg.Done()
		}()

		go func() {
			grad2, err = modules.FunctionPool.Random()()
			if err != nil {
				fmt.Println(err)
			}
			wg.Done()
		}()

		wg.Wait()

		finalgrad = xor(grad1, grad2)
		contrast := ContrastOf(finalgrad)
		if(contrast < 700) {
			approved = true
		} else {
			fmt.Println("Skipping, average contrast is too high.")
		}
		fmt.Println(contrast)
	}

	buf := new(bytes.Buffer)
	if err := png.Encode(buf, finalgrad); err != nil {
		fmt.Println(err)
		os.Exit(1)
		return nil // shut up compiler
	} else {
		bytes := buf.Bytes()
		return bytes
	}
}

func xor(img1, img2 image.Image) image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, img1.Bounds().Max.X, img1.Bounds().Max.Y))
	// For each column in the image
	for y := img1.Bounds().Min.Y; y < img1.Bounds().Max.Y; y++ {
		// and each row
		for x := img1.Bounds().Min.X; x < img1.Bounds().Max.X; x++ {
			// Get the colors of each pixel in both images
			c1_r, c1_g, c1_b, _ := img1.At(x, y).RGBA()
			c2_r, c2_g, c2_b, _ := img2.At(x, y).RGBA()
			// Then set the corresponding pixel in the final image to a xor'd version of the pixels.
			img.Set(x, y, color.NRGBA{
				R: uint8(c1_r ^ c2_r),
				G: uint8(c1_g ^ c2_g),
				B: uint8(c1_b ^ c2_b),
				A: 255,
			})
		}
	}
	return img
}

var lastTypeGiven int

func NewImage(imageType string) (image.Image, error) {
	image, err := modules.FunctionPool.Get(imageType)()
	return image, err
}

func ContrastOf(img image.Image) (float64) {
	var contrastValues float64
	var contrastNum float64
	var lastPixel float64
	for y := float64(0); y < float64(img.Bounds().Max.Y); y++ {
		var contrastValues_ []float64 
		// and each row
		for x := float64(0); x < float64(img.Bounds().Max.X); x++ {
			r, g, b, _ := img.At(int(x),int(y)).RGBA()
			contrastValues_ = append(contrastValues_, math.Abs(lastPixel-float64(r+g+b)))
			lastPixel = float64(r+g+b)
		}
		var sum float64
		for _, v := range contrastValues_{
			sum += v
		}
		sum = sum/float64(len(contrastValues_))
		contrastValues += sum
		contrastNum++
	} 
	return contrastValues/contrastNum
}