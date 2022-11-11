package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mattn/go-mastodon"
)

func MastodonThread() {
	appParent, err := mastodon.RegisterApp(context.Background(), &mastodon.AppConfig{
		Server:     LocalConfig.MastodonInstanceURL,
		ClientName: "glitchfuck",
		Scopes:     "read write follow",
		Website:    "https://github.com/IoIxD/glitchfuck",
	})

	app := mastodon.NewClient(&mastodon.Config{
		Server:       "https://wetdry.world",
		ClientID:     appParent.ClientID,
		ClientSecret: appParent.ClientSecret,
	})

	err = app.Authenticate(context.Background(), LocalConfig.MastodonEmail, LocalConfig.MastodonPassword)
	if err != nil {
		fmt.Println("Authentication error, ", err)
		os.Exit(1)
		return
	}

	duration, err := time.ParseDuration(LocalConfig.MastodonInterval)
	if err != nil {
		fmt.Println("Time parsing erro,r", err)
		os.Exit(1)
		return
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	for {
		select {
		case <-sigs:
			os.Exit(0)

		case <-WaitFor(duration):
			image, err := DefaultImage(true, 800.0, 600.0) // ignore errors since this is something that posts daily without user interaction.
			if err != nil {
				fmt.Println(err)
				return
			}

			attachment, err := app.UploadMediaFromReader(context.Background(), bytes.NewReader(image))
			if err != nil {
				fmt.Println(err)
				return
			}

			app.PostStatus(context.Background(), &mastodon.Toot{
				Status:     "",
				MediaIDs:   []mastodon.ID{attachment.ID},
				Sensitive:  false,
				Visibility: mastodon.VisibilityUnlisted,
			})
		}
	}
}
