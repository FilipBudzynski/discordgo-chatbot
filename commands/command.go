package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type CommandID int

// Const to add different commands
const (
	PingCommandID CommandID = iota
	HelpCommandID
	PlayCommandID
	UnknownCommandID
)

type Command struct {
	CommandID CommandID
	Args      []string
	Message   *discordgo.MessageCreate
}

func ParseCommand(content string) CommandID {
	switch content {
	case "!ping":
		return PingCommandID
	case "!help", "!h":
		return HelpCommandID
	case "!play", "!p":
		return PlayCommandID
	default:
		return UnknownCommandID
	}
}

// Handles commands based on the commandID
func CommandHandler(s *discordgo.Session, commandChan <-chan Command) {
	for c := range commandChan {
		GuildID := c.Message.GuildID
		UserID := c.Message.Author.ID

		switch c.CommandID {
		case PingCommandID:
			err := SendPong(s, c.Message.ChannelID)
			if err != nil {
				fmt.Println("Error with SendPong command", err)
			}
		case PlayCommandID:
			vs, err := s.State.VoiceState(GuildID, UserID)
			if err != nil {
				fmt.Println("Could not find the VoiceState", err)
				return
			}
			PlayLink(s, c.Args[1], GuildID, UserID, vs.ChannelID)
		default:
			err := sendUnknownCommand(s, c.Message.ChannelID)
			if err != nil {
				fmt.Println("Error with SendPong command", err)
			}
		}
	}
}
