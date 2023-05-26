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
	TwitterConsumerKey    	string
	TwitterConsumerSecret 	string
	TwitterAccessToken    	string
	TwitterAccessSecret   	string
	TwitterOAuthToken     	string
	TwitterOAuthSecret    	string
	TwitterInterval       	string

	MastodonInstanceURL 	string
	MastodonInterval    	string
	MastodonEmail       	string
	MastodonPassword    	string

	MastodonClientKey 		string
	MastodonClientSecret 	string
	MastodonAccessSecret 	string
	
	DiscordAuthToken 		string
	DiscordID        		string

	DiscordChannels 		[]string
	DiscordInterval 		string

	InProduction 			bool
}

var imageTypes = flag.String("types", "", "The types of images you want to generate and xor, seperated by commas.")
var cpuProfile = flag.String("cpuprofile", "", "Debug flag for what file to write the cpu profile to. CPU Profiling is disabled if this is blank.")
var memProfile = flag.String("memprofile", "", "Debug flag for what file to write the memory profile to. CPU Profiling is disabled if this is blank.")
var outputFile = flag.String("output", "test.png", "File to save output to.")
var widthOpt = flag.Float64("width", 1024, "Debug flag for what file to write the memory profile to. CPU Profiling is disabled if this is blank.")
var heightOpt = flag.Float64("height", 768, "File to save output to.")

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

	// if image types were given, run the command through them.
	if *imageTypes != "" {
		image, err := ImageViaTypes(*imageTypes, *widthOpt, *heightOpt)
		if err != nil {
			fmt.Println(err)
			return
		}
		f, _ := os.Create(*outputFile)
		_, err = f.Write(image)
		if err != nil {
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
		if LocalConfig.TwitterConsumerKey != "" {
			fmt.Println("starting twitter thread.")
			go TwitterThread()
		}
		if LocalConfig.DiscordID != "" {
			fmt.Println("starting discord thread.")
			go DiscordThread()
			go ServerThread()
		}
		if LocalConfig.MastodonClientKey != "" {
			fmt.Println("starting mastodon thread.")
			go MastodonThread()
		}
		select {}
	} else {
		image, err := DefaultImage(true, *widthOpt, *heightOpt)
		if err != nil {
			fmt.Println(err)
			return
		}
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

func ImageViaTypes(types_ string, width, height float64) ([]byte, error) {
	types := strings.Split(types_, ",")
	var lastImage image.Image
	var finalImage image.Image
	for _, v := range types {
		image, err := NewImage(v, width, height)
		if err != nil {
			return nil, err
		}
		if lastImage != nil {
			finalImage = xor(lastImage, image)
		} else {
			finalImage = image
		}
		lastImage = image
	}

	buf := new(bytes.Buffer)
	if err := png.Encode(buf, finalImage); err != nil {
		fmt.Println(err)
		os.Exit(1)
		return nil, nil // shut up compiler
	} else {
		bytes := buf.Bytes()
		return bytes, nil
	}
}

func DefaultImage(forceLowContrast bool, width, height float64) ([]byte, error) {
	var approved bool
	var finalgrad image.Image
	for !approved {
		var grad1, grad2 image.Image
		var wg sync.WaitGroup
		var err1, err2 error
		wg.Add(2)
		go func() {
			start := time.Now().UnixMilli()
			grad1, err1 = modules.FunctionPool.Random()(width, height)
			end := time.Now().UnixMilli()
			fmt.Printf("%vms\n", end-start)
			wg.Done()
		}()

		go func() {
			start := time.Now().UnixMilli()
			grad2, err2 = modules.FunctionPool.Random()(width, height)
			end := time.Now().UnixMilli()
			fmt.Printf("%vms\n", end-start)
			wg.Done()
		}()

		wg.Wait()

		if err1 != nil {
			return nil, err1
		}
		if err2 != nil {
			return nil, err2
		}

		finalgrad = xor(grad1, grad2)
		if forceLowContrast {
			contrast := ContrastOf(finalgrad)
			if contrast < 700 {
				approved = true
			} else {
				fmt.Println("Skipping, average contrast is too high.")
			}
		} else {
			approved = true
		}
	}

	buf := new(bytes.Buffer)
	if err := png.Encode(buf, finalgrad); err != nil {
		fmt.Println(err)
		os.Exit(1)
		return nil, nil // shut up compiler
	} else {
		bytes := buf.Bytes()
		return bytes, nil
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

func NewImage(imageType string, width, height float64) (image.Image, error) {
	imageFunc := modules.FunctionPool.Get(imageType)
	if imageFunc == nil {
		// it might just be that the module isn't loaded yet. wait 150ms
		time.Sleep(time.Millisecond * 150)
		imageFunc = modules.FunctionPool.Get(imageType)
		if imageFunc == nil {
			return nil, fmt.Errorf("Invalid type %v.\n", imageType)
		}
	}
	start := time.Now().UnixMilli()
	image, err := imageFunc(width, height)
	end := time.Now().UnixMilli()
	fmt.Printf("%vms for %v\n", end-start, imageType)
	return image, err
}

func ContrastOf(img image.Image) float64 {
	var contrastValues float64
	var contrastNum float64
	var lastPixel float64
	for y := float64(0); y < float64(img.Bounds().Max.Y); y++ {
		var contrastValues_ []float64
		// and each row
		for x := float64(0); x < float64(img.Bounds().Max.X); x++ {
			r, g, b, _ := img.At(int(x), int(y)).RGBA()
			contrastValues_ = append(contrastValues_, math.Abs(lastPixel-float64(r+g+b)))
			lastPixel = float64(r + g + b)
		}
		var sum float64
		for _, v := range contrastValues_ {
			sum += v
		}
		sum = sum / float64(len(contrastValues_))
		contrastValues += sum
		contrastNum++
	}
	return contrastValues / contrastNum
}
