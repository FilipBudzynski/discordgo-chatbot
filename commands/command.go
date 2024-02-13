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
	UnknownCommandID
)

type Command struct {
	AuthorID  string
	GuildID   string
	ChannelID string
	CommandID CommandID
}

func ParseCommand(content string) CommandID {
	switch content {
	case "!ping", "!p":
		return PingCommandID
	case "!help", "!h":
		return HelpCommandID
	default:
		return UnknownCommandID
	}
}

// Handles commands based on the commandID
func CommandHandler(s *discordgo.Session, commandChan <-chan Command) {
	for c := range commandChan {
		switch c.CommandID {
		case PingCommandID:
			err := SendPong(s, c.ChannelID)
			if err != nil {
				fmt.Println("Error with SendPong command", err)
			}
		default:
			err := sendUnknownCommand(s, c.ChannelID)
			if err != nil {
				fmt.Println("Error with SendPong command", err)
			}
		}
	}
}
