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
		{
			Name: 	"width",
			Description: "Width of the image. Default is 640.",
			Type: discordgo.ApplicationCommandOptionInteger,
			MaxValue: 1024,
		},
		{
			Name: 	"height",
			Description: "Height of the image. Default is 480.",
			Type: discordgo.ApplicationCommandOptionInteger,
			MaxValue: 768,
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

	RefreshSlashCommands(true)
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
				RefreshSlashCommands(false)
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
func RefreshSlashCommands(nameServers bool) {
	for _, v := range discord.State.Guilds {
		_, err := discord.ApplicationCommandCreate(discord.State.User.ID, v.ID, &command)
		if(err != nil) {
			fmt.Println(err)
		}
		if(nameServers) {
			fmt.Printf("initialized command in %v\n", v.Name)
		}
	}
}


// The main command 
func mainCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var sent bool
	// Thread to check if the command has finished within the time discord allows.
	go func() {
		time.Sleep(time.Millisecond * 2800)
		if(!sent) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "The bot spend too long regenerating the image and Discord will not let bots take longer then three seconds. Sorry.",
				},
			})
			return
		}
	}()
	// Thread that generates the image.
	go func() {
		options := i.ApplicationCommandData().Options
		optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
		for _, v := range options {
			optionMap[v.Name] = v
		}

		// unmarshal options into go.
		forceLowContrast := false 
		if flc, ok := optionMap["forcelowcontrast"]; ok {
			forceLowContrast_, ok := flc.Value.(bool)
			if(ok) {
				forceLowContrast = forceLowContrast_
			}
		}
		var width float64 = 640
		if widthop, ok := optionMap["width"]; ok {
			switch t := widthop.Value.(type) {
				case int64:
					width = float64(widthop.Value.(int64))
				case uint64:
					width = float64(widthop.Value.(uint64))
				case float64:
					width = float64(widthop.Value.(float64))
				default: 
					fmt.Println(t)
			}
		}

		var height float64 = 480
		if heightop, ok := optionMap["height"]; ok {
			switch t := heightop.Value.(type) {
				case int64:
					height = float64(heightop.Value.(int64))
				case uint64:
					height = float64(heightop.Value.(uint64))
				case float64:
					height = float64(heightop.Value.(float64))
				default: 
					fmt.Println(t)
			}
		}
		var image []byte
		var content string
		if types, ok := optionMap["types"]; ok {
			image, err = ImageViaTypes(types.Value.(string),width,height)
			content = "`"+types.Value.(string)+"`"
		} else {
			image, err = DefaultImage(forceLowContrast,width,height)
			content = "­"
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
				Content: content,
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