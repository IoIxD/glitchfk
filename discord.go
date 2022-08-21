package main

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

var discord *discordgo.Session
var err error

var command = discordgo.ApplicationCommand{
	Name: 	"glitchfuck",
	Description: "Runs glitchfuck.",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name: 	"types",
			Description: "The types of images to generate, seperated by commas. Random image by default. No docs yet.",
			Type: discordgo.ApplicationCommandOptionString,
		},
		{
			Name: 	"forcelowcontrast",
			Description: "Don't return image if average contrast is high. Used for Twitter bot. May cause bot to not respond.",
			Type: discordgo.ApplicationCommandOptionBoolean,
		},
	},
}

var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"glitchfuck": mainCommand,
	}

func DiscordThread() {
	discord, err = discordgo.New("Bot " + LocalConfig.DiscordAuthToken)
	if(err != nil) {
		fmt.Println(err)
		os.Exit(1)
	}


	discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
	discord.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		fmt.Printf("Logged in as: %v#%v\n", s.State.User.Username, s.State.User.Discriminator)
	})

	RemoveSlashCommands()

	RefreshSlashCommands()
	go RefreshSlashCommandsThread()

	discord.Open()
	for {}
}

// Thread for refreshing the slash commands every minute.
func RefreshSlashCommandsThread() {
	ticker := time.NewTicker(3 * time.Second)
	for {
		select {
			case <-ticker.C:
				RefreshSlashCommands()
		}
	}
}

// Remove the slash commands
func RemoveSlashCommands() {
	for _, v := range discord.State.Guilds {
		registeredCommands, err := discord.ApplicationCommands(discord.State.User.ID, v.ID)
		if err != nil {
			fmt.Printf("Could not fetch registered commands: %v\n", err)
		}

		for _, n := range registeredCommands {
			err := discord.ApplicationCommandDelete(discord.State.User.ID, v.ID, n.ID)
			if(err != nil) {
				fmt.Println(err)
			}
		}
	}
}


// Refresh the slash commands
func RefreshSlashCommands() {
	for _, v := range discord.State.Guilds {
		_, err := discord.ApplicationCommandCreate(discord.State.User.ID, v.ID, &command)
		if(err != nil) {
			fmt.Println(err)
		}
	}
}


// The main command 
func mainCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var sent bool
	go func() {
		time.Sleep(time.Millisecond * 2500)
		if(!sent) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "The bot spend too long regenerating the image and Discord will not let bots take longer then three seconds. Sorry.",
				},
			})
		}
	}()
	go func() {
		options := i.ApplicationCommandData().Options

		optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
		for _, v := range options {
			optionMap[v.Name] = v
		}

		forceLowContrast := false 
		if flc, ok := optionMap["forcelowcontrast"]; ok {
			forceLowContrast = flc.Value.(bool)
		}

		var image []byte
		if types, ok := optionMap["types"]; ok {
			image, err = ImageViaTypes(types.Value.(string))
		} else {
			image, err = DefaultImage(forceLowContrast)
		}

		if(err != nil) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Error! ```\n"+err.Error()+"\n```",
				},
			})
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "­",
				Files: []*discordgo.File{
					{
						ContentType: "image/png",
						Name:        "glitchfuck.png",
						Reader:      bytes.NewReader(image),
					},
				},
			},
		})

		sent = true
	}()
}