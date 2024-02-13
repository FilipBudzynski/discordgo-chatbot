package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// Action given when the unknown command is invoked
func sendUnknownCommand(s *discordgo.Session, channelID string) error {
	_, err := s.ChannelMessageSend(channelID, "Unknown Command")
	if err != nil {
		fmt.Println("Error sending message to channel:", err)
		return err
	}
	return nil
}
