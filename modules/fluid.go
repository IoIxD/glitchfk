package modules

import (
	"fmt"
	"image"
)

// "Fluid" image; https://www.shadertoy.com/view/XdcGW2

func init() {
	//FunctionPool.Add("fluid", NewFluidImage)
}

func NewFluidImage() (image.Image, error) {
	fmt.Println("screenshotting fluid image (https://www.shadertoy.com/view/XdcGW2)...")
	return ShadertoyScreenshot("https://www.shadertoy.com/view/XdcGW2")
}
