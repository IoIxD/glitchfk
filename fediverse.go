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
	app := mastodon.NewClient(&mastodon.Config{
		Server:       LocalConfig.MastodonInstanceURL,
		ClientID:     LocalConfig.MastodonClientKey,
		ClientSecret: LocalConfig.MastodonClientSecret,
	})

	err := app.Authenticate(context.Background(), LocalConfig.MastodonEmail, LocalConfig.MastodonPassword)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	fmt.Println("Posting every", LocalConfig.MastodonInterval)

	duration, err := time.ParseDuration(LocalConfig.MastodonInterval)
	if err != nil {
		fmt.Println(err)
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
			image, _, _, err := DefaultImage(true, 800.0, 600.0) // ignore errors since this is something that posts daily without user interaction.
			if err != nil {
				fmt.Println(err)
				continue
			}

			attachment, err := app.UploadMediaFromReader(context.Background(), bytes.NewReader(image))
			if err != nil {
				fmt.Println(err)
				continue
			}

			app.PostStatus(context.Background(), &mastodon.Toot{
				Status:      "",
				MediaIDs:    []mastodon.ID{attachment.ID},
				Sensitive:   false,
				SpoilerText: "eyestrain",
				Visibility:  mastodon.VisibilityUnlisted,
			})
		}
	}
}
