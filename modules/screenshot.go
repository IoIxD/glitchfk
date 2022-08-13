package modules

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"os/exec"

	"golang.org/x/image/draw"
)

// go bridge for the node.js code that screenshots shadertoy webpages

func ShadertoyScreenshot(url string) (image.Image, error) {
	// run the command
	os.Setenv("WEBSITE",url)
	cmd := exec.Command("node","./CEF/index.js")

	// get the stderr and stdout pipes.
	stderrPipe, err := cmd.StderrPipe()
	if(err != nil) {
		return nil, err
	}
	stdoutPipe, err := cmd.StdoutPipe()
	if(err != nil) {
		return nil, err
	}

	// start the command.
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	
	// read from the pipes from earlier
	stderr, _ := io.ReadAll(stderrPipe)
	if(len(stderr) >= 1) {
		return nil, fmt.Errorf("%v",string(stderr))
	}

	stdout, _ := io.ReadAll(stdoutPipe)
	if(len(stdout) >= 1) {
		return nil, fmt.Errorf("%v",string(stdout))
	}

	// Wait.
	cmd.Wait();

	// Open the resulting image
	file, err := os.Open("./screenshot.png")
	if(err != nil) {
		return nil, err
	}

	// So now we want to decode it...
	finalImageUncropped, err := png.Decode(file)
	if(err != nil) {
		return nil, err
	}

	// But we want to crop it, first.
	finalImageUnresized := image.NewNRGBA(image.Rect(0,0,finalImageUncropped.Bounds().Max.X-22,finalImageUncropped.Bounds().Max.Y-106))
	
	for y := 106; y < finalImageUncropped.Bounds().Max.Y; y++ {
		for x := 11; x < finalImageUncropped.Bounds().Max.X-11; x++ {
			value := finalImageUncropped.At(x,y)
			finalImageUnresized.Set(x-11,y-106,value)
		}
	}

	// and then resize it.
	finalImage := image.NewNRGBA(image.Rect(0,0,int(WIDTH),int(HEIGHT)))
	draw.NearestNeighbor.Scale(finalImage,finalImage.Bounds(),finalImageUnresized,finalImageUnresized.Bounds(), draw.Over,nil)

	return finalImage, nil
}