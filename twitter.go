package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dghubble/oauth1"
	"github.com/drswork/go-twitter/twitter"
)

const NORMAL_WIDTH float64 = 1024
const NORMAL_HEIGHT float64 = 768

const SUNDAY_WIDTH float64 = 1920
const SUNDAY_HEIGHT float64 = 1080

var TwitterWidth float64
var TwitterHeight float64


var client *twitter.Client
var OAuth1Config *oauth1.Config
var OAuth1Token *oauth1.Token
var Oauth1Client *http.Client

func InitTwitter(token, secret string) {
	OAuth1Config = oauth1.NewConfig(LocalConfig.TwitterConsumerKey,LocalConfig.TwitterConsumerSecret)
	OAuth1Token = oauth1.NewToken(token,secret)
	Oauth1Client = OAuth1Config.Client(oauth1.NoContext,OAuth1Token)
	client = twitter.NewClient(Oauth1Client)


	if(time.Now().Weekday() == 0) {
		TwitterWidth = SUNDAY_WIDTH
		TwitterHeight = SUNDAY_HEIGHT
	} else {
		TwitterWidth = NORMAL_WIDTH
		TwitterHeight = NORMAL_HEIGHT
	}
} 

func TwitterThread() {
	duration, err := time.ParseDuration(LocalConfig.TwitterInterval)
	if(err != nil) {
		fmt.Println(err)
		return
	}

	InitTwitter(LocalConfig.TwitterOAuthToken,LocalConfig.TwitterOAuthSecret)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	for {
		select {
			case <-sigs:
				os.Exit(0)

			case <-WaitFor(duration):
				image, _ := DefaultImage(true, TwitterWidth, TwitterHeight) // ignore errors since this is something that posts daily without user interaction.
				result, resp, err := client.Media.Upload(image,"image/png")
				if(err != nil) {
					fmt.Println(err)
					return
				}
				if(resp.StatusCode != 201) {
					fmt.Printf("%v.\n No other response given; the full struct is %v\n",resp.Status,resp)
					return
				}
				client.Statuses.Update("",&twitter.StatusUpdateParams{
					MediaIds: []int64{result.MediaID},
				})
		}
	}
}