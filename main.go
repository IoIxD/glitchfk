package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/dghubble/oauth1"
	"github.com/drswork/go-twitter/twitter"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/mazznoer/colorgrad"
	"github.com/pelletier/go-toml/v2"
	"github.com/unixpickle/polish/polish"
)

var client *twitter.Client
var OAuth1Config *oauth1.Config
var OAuth1Token *oauth1.Token
var Oauth1Client *http.Client

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

	InitTwitter(LocalConfig.OAuthToken,LocalConfig.OAuthSecret)
	image := NewImage()

	// if we're in production, tweet out the result 
	if(LocalConfig.InProduction) {
		image, resp, err := client.Media.Upload(image,"image/png")
		if(err != nil) {
			fmt.Println(err)
			return
		}
		if(resp.StatusCode != 201) {
			fmt.Printf("%v.\n No other response given; the full struct is %v\n",resp.Status,resp)
			return
		}
		client.Statuses.Update("",&twitter.StatusUpdateParams{
			MediaIds: []int64{image.MediaID},
		})
	} else {
		f, _ := os.Create("test.png")
		f.Write(image)
		f.Close()
	}

}

func InitTwitter(token, secret string) {
	OAuth1Config = oauth1.NewConfig(LocalConfig.ConsumerKey,LocalConfig.ConsumerSecret)
	OAuth1Token = oauth1.NewToken(token,secret)
	Oauth1Client = OAuth1Config.Client(oauth1.NoContext,OAuth1Token)
	client = twitter.NewClient(Oauth1Client)
} 

func NewImage() ([]byte) {
	grad1, err := NewGradient(randomGradientType())
	if(err != nil) {
		fmt.Println(err)
	}

	grad2, err := NewGradient(randomGradientType())
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

func NewGradient(mode string) (image.Image, error) {
	fmt.Printf("generating %v gradient...\n",mode)
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
		fmt.Println(err)
		return nil, err
	}

	img := image.NewNRGBA(image.Rect(0, 0, WIDTH, HEIGHT));
	// For each column in the image
	for y := float64(0); y < HEIGHT; y++ {
		// and each row
		for x := float64(0); x < WIDTH; x++ {
			var colorfulColor colorful.Color
			switch(mode) {
				case "horizontal": colorfulColor = grad.At(x/WIDTH)
				case "vertical": colorfulColor = grad.At(y/HEIGHT)
				case "diagonal": colorfulColor = grad.At((x/WIDTH)+(y/HEIGHT))
				case "radial":
					pos := math.Cos(x/WIDTH)+math.Sin(y/HEIGHT)
					colorfulColor = grad.At((pos-1))
				case "inverse-radial":
					pos := math.Cos(y/HEIGHT)+math.Sin(x/WIDTH)
					colorfulColor = grad.At((pos-1))
				case "fucked":
					posx := (x/WIDTH)
					posy := (y/HEIGHT)
					pos := fucked(posx,posy,(math.Cos(posx)+math.Sin(posy)))
					colorfulColor = grad.At(pos)
			}
			// Set the corresponding pixel
			img.Set(int(x),int(y), color.NRGBA{
				R: uint8(colorfulColor.R*256),
				G: uint8(colorfulColor.G*256),
				B: uint8(colorfulColor.B*256),
				A: 255,
			})
		}
	}
	if(mode == "fucked") {
		fmt.Printf("denoising gradient...\n")
		rand.Seed(time.Now().UnixNano())
		choice := rand.Intn(2)
		switch(choice) {
			case 0: return polish.PolishImage(polish.ModelTypeShallow, img), nil
			default: return polish.PolishImage(polish.ModelTypeDeep, img), nil
		}
	}
	return img, nil
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

var lastGradientTypeGiven int
func randomGradientType() (value string) {
	values := []string{"horizontal","vertical","diagonal","radial","inverse-radial","fucked"}
	var choice int
	for(choice == lastGradientTypeGiven) {
		rand.Seed(time.Now().UnixNano())
		choice = rand.Intn(len(values))
	}
	lastGradientTypeGiven = choice
	return values[choice]
}


func fucked(x, y, r float64) (float64) {
	rand.Seed(time.Now().UnixNano())
	choice := rand.Intn(64) 
	switch(choice) {
		case 0: return math.Pow(x,2)+math.Pow(y,2)+math.Pow(r,2);
		case 1: return math.Pow(x,2)+math.Pow(y,2)-math.Pow(r,2);
		case 2: return math.Pow(x,2)+math.Pow(y,2)*math.Pow(r,2);
		case 3: return math.Pow(x,2)+math.Pow(y,2)/math.Pow(r,2);
		case 4: return math.Pow(x,2)-math.Pow(y,2)+math.Pow(r,2);
		case 5: return math.Pow(x,2)-math.Pow(y,2)-math.Pow(r,2);
		case 6: return math.Pow(x,2)-math.Pow(y,2)*math.Pow(r,2);
		case 7: return math.Pow(x,2)-math.Pow(y,2)/math.Pow(r,2);
		case 8: return math.Pow(x,2)*math.Pow(y,2)+math.Pow(r,2);
		case 9: return math.Pow(x,2)*math.Pow(y,2)-math.Pow(r,2);
		case 10: return math.Pow(x,2)*math.Pow(y,2)*math.Pow(r,2);
		case 11: return math.Pow(x,2)*math.Pow(y,2)/math.Pow(r,2);
		case 12: return math.Pow(x,2)/math.Pow(y,2)+math.Pow(r,2);
		case 13: return math.Pow(x,2)/math.Pow(y,2)-math.Pow(r,2);
		case 14: return math.Pow(x,2)/math.Pow(y,2)*math.Pow(r,2);
		case 15: return math.Pow(x,2)/math.Pow(y,2)/math.Pow(r,2);
		case 16: return math.Pow(x,2)+math.Pow(y,2)+math.Pow(r,2);
		case 17: return math.Pow(x,2)-math.Pow(y,2)+math.Pow(r,2);
		case 18: return math.Pow(x,2)*math.Pow(y,2)+math.Pow(r,2);
		case 19: return math.Pow(x,2)/math.Pow(y,2)+math.Pow(r,2);
		case 20: return math.Pow(x,2)+math.Pow(y,2)-math.Pow(r,2);
		case 21: return math.Pow(x,2)-math.Pow(y,2)-math.Pow(r,2);
		case 22: return math.Pow(x,2)*math.Pow(y,2)-math.Pow(r,2);
		case 23: return math.Pow(x,2)/math.Pow(y,2)-math.Pow(r,2);
		case 24: return math.Pow(x,2)+math.Pow(y,2)*math.Pow(r,2);
		case 25: return math.Pow(x,2)-math.Pow(y,2)*math.Pow(r,2);
		case 26: return math.Pow(x,2)*math.Pow(y,2)*math.Pow(r,2);
		case 27: return math.Pow(x,2)/math.Pow(y,2)*math.Pow(r,2);
		case 28: return math.Pow(x,2)+math.Pow(y,2)/math.Pow(r,2);
		case 29: return math.Pow(x,2)-math.Pow(y,2)/math.Pow(r,2);
		case 30: return math.Pow(x,2)*math.Pow(y,2)/math.Pow(r,2);
		case 31: return math.Pow(x,2)/math.Pow(y,2)/math.Pow(r,2);

		case 32: return x+y+math.Pow(r,2);
		case 33: return x+y-math.Pow(r,2);
		case 34: return x+y*math.Pow(r,2);
		case 35: return x+y/math.Pow(r,2);
		case 36: return x-y+math.Pow(r,2);
		case 37: return x-y-math.Pow(r,2);
		case 38: return x-y*math.Pow(r,2);
		case 39: return x-y/math.Pow(r,2);
		case 40: return x*y+math.Pow(r,2);
		case 41: return x*y-math.Pow(r,2);
		case 42: return x*y*math.Pow(r,2);
		case 43: return x*y/math.Pow(r,2);
		case 44: return x/y+math.Pow(r,2);
		case 45: return x/y-math.Pow(r,2);
		case 46: return x/y*math.Pow(r,2);
		case 47: return x/y/math.Pow(r,2);
		case 48: return x+y+math.Pow(r,2);
		case 49: return x-y+math.Pow(r,2);
		case 50: return x*y+math.Pow(r,2);
		case 51: return x/y+math.Pow(r,2);
		case 52: return x+y-math.Pow(r,2);
		case 53: return x-y-math.Pow(r,2);
		case 54: return x*y-math.Pow(r,2);
		case 55: return x/y-math.Pow(r,2);
		case 56: return x+y*math.Pow(r,2);
		case 57: return x-y*math.Pow(r,2);
		case 58: return x*y*math.Pow(r,2);
		case 59: return x/y*math.Pow(r,2);
		case 60: return x+y/math.Pow(r,2);
		case 61: return x-y/math.Pow(r,2);
		case 62: return x*y/math.Pow(r,2);
		case 63: return x/y/math.Pow(r,2);
	}
	return 0;
}

/*
TODO come back to this when i know trigonometry so i can know how not to fuck this up

// https://skia.org/docs/dev/design/conical/
func conial(x, y, r float64) (float64) {
	if(r == 1) {
		lvalue := (math.Pow(x,2)+math.Pow(y,2))/((1+r)*x)
		rvalue := (math.Pow(x,2)+math.Pow(y,2))/(2*x)
		return lvalue+rvalue
	} else {
		fuckyou := (math.Sqrt((math.Pow(r,2)-1)*math.Pow(y,2)+math.Pow(r,2)*math.Pow(x,2))-x)/(math.Pow(r,2)-1)
		return fuckyou
		//if(r > 1) {
		//	return math.Abs(fuckyou)
		//} else {
		//	return fuckyou
		//}
	}
}
*/