package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/pelletier/go-toml/v2"
)

var LocalConfig struct {
	ConsumerKey 	string
	ConsumerSecret 	string
	AccessToken 	string
	AccessSecret	string
	OAuthToken 		string
	OAuthSecret 	string
	Interval		string
	InProduction	bool
}

const WIDTH = 1024;
const HEIGHT = 768;

const DETAIL = 32768;

func main() {	
	// flag for setting up manual generations
	var imageTypes string
	flag.StringVar(&imageTypes,"types", "", "The types of images you want to generate and xor, seperated by commas.")
	flag.Parse()
	if(imageTypes != "") {
		types := strings.Split(imageTypes,",")
		var lastImage image.Image
		var finalImage image.Image
		for _, v := range types {
			image, err := NewImage(v)
			if(err != nil) {
				fmt.Println(err)
				return
			}
			if(lastImage != nil) {
				finalImage = xor(lastImage,image)
			} else {
				finalImage = image
			}
			lastImage = image
		}
		f, _ := os.Create("test.png")
		if err := png.Encode(f,finalImage); err != nil {
			fmt.Println(err)
		}
		return
	}

	// set up the config
	f, err := os.Open("config.toml")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = toml.NewDecoder(f).Decode(&LocalConfig)
	if err != nil {
		fmt.Println(err)
		return
	}

	if(LocalConfig.InProduction) {
		TwitterThread()
	} else {
		image := DefaultImage()
		f, _ := os.Create("test.png")
		f.Write(image)
		f.Close()
	}

}

func WaitFor(int time.Duration) <-chan time.Time {
	now := time.Now()
	dur := now.Truncate(int).Add(int).Sub(now)
	return time.After(dur)
}

func DefaultImage() ([]byte) {
	grad1, err := NewImage(randomImageType())
	if(err != nil) {
		fmt.Println(err)
	}

	grad2, err := NewImage(randomImageType())
	if(err != nil) {
		fmt.Println(err)
	}

	finalgrad := xor(grad1,grad2)

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

func xor(img1, img2 image.Image) (image.Image) {
	img := image.NewNRGBA(image.Rect(0, 0, img1.Bounds().Max.X, img1.Bounds().Max.Y));
	// For each column in the image
	for y := img1.Bounds().Min.Y; y < img1.Bounds().Max.Y; y++ {
		// and each row
		for x := img1.Bounds().Min.X; x < img1.Bounds().Max.X; x++ {
			// Get the colors of each pixel in both images
			c1_r, c1_g, c1_b, _ := img1.At(x, y).RGBA() 
			c2_r, c2_g, c2_b, _ := img2.At(x, y).RGBA() 
			// Then set the corresponding pixel in the final image to a xor'd version of the pixels.
			img.Set(x,y, color.NRGBA{
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
func randomImageType() (value string) {
	values := []string{"horizontal","vertical","diagonal","radial","inverse-radial","fucked","noise","wave"}
	var choice int
	for(choice == lastTypeGiven) {
		rand.Seed(time.Now().UnixNano())
		choice = rand.Intn(len(values))
	}
	lastTypeGiven = choice
	return values[choice]
}


func NewImage(imageType string) (image.Image, error) {
	switch(imageType) {
		case "horizontal","vertical","diagonal","radial","inverse-radial","fucked":
			return NewGradientImage(imageType)
		case "noise":
			return NewNormalNoise()
		case "wave":
			return NewWave()
	}
	return nil, nil
}